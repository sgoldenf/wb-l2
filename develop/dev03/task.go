package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	cfg := parseFlags()
	lines, err := parseFiles(cfg)
	if err != nil {
		log.Fatal(err)
	}
	sortLines(cfg, lines)
	if cfg.u {
		processFlagU(&lines)
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}

var (
	errInvalidArgument = errors.New("invalid argument")
)

type kFlags []uint

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
func (k *kFlags) Set(value string) error {
	column, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return err
	}
	if column == 0 {
		return errInvalidArgument
	}
	*k = append(*k, uint(column))
	return nil
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (k *kFlags) String() string {
	return fmt.Sprint(*k)
}

type config struct {
	files []string
	k     kFlags
	regex *regexp.Regexp
	n     bool
	r     bool
	u     bool
}

func parseFlags() *config {
	var cfg config
	flag.Var(&cfg.k, "k", "sort by column (0 - default sort)")
	flag.BoolVar(&cfg.n, "n", false, "numeric sort")
	flag.BoolVar(&cfg.r, "r", false, "sort in reverse order")
	flag.BoolVar(&cfg.u, "u", false, "unique sort")
	flag.Parse()
	if cfg.n {
		cfg.regex = regexp.MustCompile(`^(\+|-)?\d+(.\d+)?(e(\+|-)?\d+)?`)
	}
	files := flag.Args()
	if len(files) > 0 {
		cfg.files = make([]string, 0, len(files))
		cfg.files = append(cfg.files, files...)
	}
	return &cfg
}

func parseFiles(cfg *config) ([]string, error) {
	lines := []string{}
	if cfg.files != nil {
		for _, file := range cfg.files {
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}
			parseFile(&lines, f)
		}
	} else {
		parseFile(&lines, os.Stdin)
	}
	return lines, nil
}

func parseFile(lines *[]string, r io.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		*lines = append(*lines, s.Text())
	}
}

func sortLines(cfg *config, lines []string) {
	sort.SliceStable(lines, func(i, j int) bool {
		str1 := lines[i]
		str2 := lines[j]
		less := compareLines(cfg, str1, str2)
		if cfg.r {
			less = !less
		}
		return less
	})
}

func compareLines(cfg *config, str1, str2 string) bool {
	var less bool
	if cfg.k != nil {
		for _, k := range cfg.k {
			col1, hasColumn1 := column(str1, k)
			col2, hasColumn2 := column(str2, k)
			if hasColumn1 && !hasColumn2 {
				less = false
			} else if !hasColumn1 && hasColumn2 {
				less = true
			} else if hasColumn1 && hasColumn2 && col1 != col2 {
				less = compare(cfg, col1, col2)
			} else {
				less = compare(cfg, str1, str2)
			}
			if hasColumn1 && hasColumn2 && col1 == col2 {
				continue
			}
			break
		}
	} else {
		less = compareLines(cfg, str1, str2)
	}
	return less
}

func compare(cfg *config, str1, str2 string) bool {
	if cfg.n {
		return compareNumeric(str1, str2, cfg.regex)
	}
	return str1 < str2
}

func column(str string, column uint) (string, bool) {
	columns := strings.Fields(str)
	if int(column) > len(columns) {
		return str, false
	}
	return columns[column-1], true
}

func compareNumeric(str1, str2 string, regex *regexp.Regexp) bool {
	b1 := regex.Find([]byte(str1))
	b2 := regex.Find([]byte(str2))
	if b1 == nil && b2 == nil {
		return str1 < str2
	}
	if b1 == nil {
		return true
	}
	if b2 == nil {
		return false
	}
	num1, err1 := strconv.ParseFloat(string(b1), 64)
	num2, err2 := strconv.ParseFloat(string(b2), 64)
	if err1 != nil && err2 != nil {
		return str1 < str2
	}
	if err1 != nil {
		return true
	}
	if err2 != nil {
		return false
	}
	if num1 == num2 {
		return str1 < str2
	}
	return num1 < num2
}

func processFlagU(lines *[]string) {
	for i := 0; i < len(*lines)-1; i++ {
		if (*lines)[i] == (*lines)[i+1] {
			*lines = append((*lines)[:i], (*lines)[i+1:]...)
		}
	}
}
