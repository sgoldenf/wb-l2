package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	fmt.Print("Input string to unpack: ")
	var input string
	fmt.Scanln(&input)
	str, err := stringUnpack(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println(str)
	}
}

var (
	errorInvalidString = errors.New("invalid string")
)

func stringUnpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}
	if str[0] >= '0' && str[0] <= '9' {
		return "", errorInvalidString
	}
	var b strings.Builder
	var r rune
	count := 0
	for _, cur := range str {
		if cur >= '0' && cur <= '9' {
			if count == -1 {
				count = 0
			}
			count = count*10 + (int(cur) - 48)
		} else {
			addRune(&b, r, count)
			count = -1
			r = cur
		}
	}
	addRune(&b, r, count)
	return b.String(), nil
}

func addRune(b *strings.Builder, r rune, count int) {
	if count == -1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		b.WriteRune(r)
	}
}
