package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	cfg, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		cut(cfg, s.Text())
	}
}

var (
	errInvalidArgument = errors.New("invalid argument")
	errNofields        = errors.New("no fields provided")
)

type fFlags []uint

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
func (f *fFlags) Set(value string) error {
	field, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return err
	}
	if field == 0 {
		return errInvalidArgument
	}
	*f = append(*f, uint(field))
	return nil
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (f *fFlags) String() string {
	return fmt.Sprint(*f)
}

type config struct {
	f fFlags
	d string
	s bool
}

func parseFlags() (*config, error) {
	var cfg config
	flag.Var(&cfg.f, "f", "cut field from list")
	flag.StringVar(&cfg.d, "d", "\t", "delimiter")
	flag.BoolVar(&cfg.s, "s", false, "suppress lines with no field delimiter characters")
	flag.Parse()
	if cfg.f == nil {
		flag.Usage()
		return nil, errNofields
	}
	sort.Slice(cfg.f, func(i, j int) bool {
		return cfg.f[i] < cfg.f[j]
	})
	for i := 0; i < len(cfg.f)-1; i++ {
		if cfg.f[i] == cfg.f[i+1] {
			cfg.f = append(cfg.f[:i], cfg.f[i+1:]...)
		}
	}
	if len(cfg.d) > 1 {
		return nil, errInvalidArgument
	}
	return &cfg, nil
}

func cut(cfg *config, str string) {
	if !strings.Contains(str, cfg.d) {
		if !cfg.s {
			fmt.Println(str)
		}
		return
	}
	fields := strings.Split(str, cfg.d)
	toPrint := []string{}
	if len(fields) >= int(cfg.f[0]) {
		for i := 0; i < len(cfg.f) && int(cfg.f[i]) <= len(fields); i++ {
			toPrint = append(toPrint, fields[cfg.f[i]-1])
		}
	}
	fmt.Println(strings.Join(toPrint, cfg.d))
}

// package main

// import (
// 	"flag"
// 	"fmt"
// )

// type arrayFlags []string

// func (i *arrayFlags) String() string {
// 	return "my string representation"
// }

// func (i *arrayFlags) Set(value string) error {
// 	*i = append(*i, value)
// 	return nil
// }

// var myFlags arrayFlags

// func main() {
// 	flag.Var(&myFlags, "list1", "Some description for this param.")
// 	flag.Parse()
// 	fmt.Println(myFlags)
// }
