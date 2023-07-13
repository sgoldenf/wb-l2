package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	u, err := initGrepUtil()
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range u.cfg.files {
		u.currentFile = file
		var r *os.File
		var err error
		if file != "" {
			r, err = os.Open(file)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			r = os.Stdin
		}
		u.grep(r)
	}
}

var (
	errNoPattern = errors.New("no pattern to match in arguments")
)

type grepUtil struct {
	cfg         *config
	currentFile string
	beforeBuf   []string
	afterBuf    uint
	prevLineNum uint
	curLineNum  uint
	count       uint
}

type config struct {
	files   []string
	pattern string
	before  uint
	after   uint
	ctx     uint
	c       bool
	i       bool
	v       bool
	fixed   bool
	n       bool
}

func initGrepUtil() (*grepUtil, error) {
	var u grepUtil
	cfg, err := parseFlags()
	if err != nil {
		return nil, err
	}
	u.cfg = cfg
	if cfg.ctx > 0 {
		cfg.before = cfg.ctx
		cfg.after = cfg.ctx
	}
	if cfg.c {
		cfg.before = 0
		cfg.after = 0
	}
	u.beforeBuf = make([]string, 0, cfg.before)
	if cfg.i {
		cfg.pattern = strings.ToLower(cfg.pattern)
	}
	return &u, nil
}

func parseFlags() (*config, error) {
	var cfg config
	flag.UintVar(&cfg.before, "B", 0, "print n lines before match")
	flag.UintVar(&cfg.after, "A", 0, "print n lines after match")
	flag.UintVar(&cfg.ctx, "C", 0, "print n lines before and after match")
	flag.BoolVar(&cfg.c, "c", false, "print line count")
	flag.BoolVar(&cfg.i, "i", false, "ignore case")
	flag.BoolVar(&cfg.v, "v", false, "inverted")
	flag.BoolVar(&cfg.fixed, "F", false, "exact string match")
	flag.BoolVar(&cfg.n, "n", false, "print line number")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return nil, errNoPattern
	}
	cfg.pattern = args[0]
	if len(args) > 1 {
		cfg.files = make([]string, 0, len(args)-1)
		cfg.files = append(cfg.files, args[1:]...)
	} else {
		cfg.files = []string{""}
	}
	return &cfg, nil
}

func (u *grepUtil) grep(r io.Reader) {
	s := bufio.NewScanner(r)
	u.count = 0
	u.curLineNum = 0
	for s.Scan() {
		u.curLineNum++
		line := s.Text()
		match := u.checkMatch(line)
		if match {
			u.processMatch(line)
		}
		if !u.cfg.c && !match {
			u.processNotMatch(line)
		}
	}
	if u.cfg.c {
		if len(u.cfg.files) > 1 {
			fmt.Printf("%s:", u.currentFile)
		}
		fmt.Println(u.count)
	}
}

func (u *grepUtil) checkMatch(str string) bool {
	str = strings.TrimSpace(str)
	if u.cfg.i {
		str = strings.ToLower(str)
	}
	var match bool
	if u.cfg.fixed {
		match = str == u.cfg.pattern
	} else {
		match = strings.Contains(str, u.cfg.pattern)
	}
	if u.cfg.v {
		match = !match
	}
	return match
}

func (u *grepUtil) processMatch(match string) {
	u.count++
	u.afterBuf = u.cfg.after
	if !u.cfg.c {
		u.printMatch(match)
	}
}

func (u *grepUtil) processNotMatch(line string) {
	if u.afterBuf > 0 {
		u.afterBuf--
		u.printLine(line, true)
	} else if u.cfg.before > 0 {
		u.addToBeforeBuf(line)
	}
}

func (u *grepUtil) printMatch(match string) {
	if len(u.beforeBuf) > 0 {
		u.curLineNum -= uint(len(u.beforeBuf))
		for i := range u.beforeBuf {
			u.printLine(u.beforeBuf[i], true)
			u.curLineNum++
		}
		u.beforeBuf = u.beforeBuf[0:0:u.cfg.before]
	}
	u.printLine(match, false)
}

func (u *grepUtil) printLine(line string, fromBuf bool) {
	if (u.cfg.before > 0 || u.cfg.after > 0) && (u.curLineNum != u.prevLineNum+1) {
		fmt.Println("--")
	}
	if len(u.cfg.files) > 1 {
		fmt.Printf("%s:", u.currentFile)
	}
	if u.cfg.n {
		fmt.Print(u.curLineNum)
		if fromBuf {
			fmt.Print("-")
		} else {
			fmt.Print(":")
		}
	}
	fmt.Println(line)
	u.prevLineNum = uint(u.curLineNum)
}

func (u *grepUtil) addToBeforeBuf(line string) {
	if len(u.beforeBuf) < int(u.cfg.before) {
		u.beforeBuf = append(u.beforeBuf, line)
	} else {
		copy(u.beforeBuf, u.beforeBuf[1:])
		u.beforeBuf[u.cfg.before-1] = line
	}
}
