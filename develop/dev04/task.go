package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	m := anagramMap([]string{"тяпка", "пятак", "пятка", "слиток", "листок", "столик", "пяТак", "лиСТОк", "слиТОК", "один"})
	fmt.Println(m)
}

func anagramMap(words []string) map[string][]string {
	anagrams := make(map[string][]string)
	for _, word := range words {
		low := strings.ToLower(word)
		anagram := findAnagram(anagrams, low)
		if _, ok := anagrams[anagram]; !ok {
			anagrams[anagram] = []string{}
		}
		if !isMapped(anagrams[anagram], low) {
			anagrams[anagram] = append(anagrams[anagram], low)
		}
	}
	for k, v := range anagrams {
		if len(v) < 2 {
			delete(anagrams, k)
			continue
		}
		sort.Strings(anagrams[k])
	}
	return anagrams
}

func findAnagram(anagrams map[string][]string, word string) string {
	for k := range anagrams {
		if isAnagram(k, word) {
			return k
		}
	}
	return word
}

func isAnagram(word1, word2 string) bool {
	if len(word1) != len(word2) {
		return false
	}
	map1 := make(map[rune]int, len(word1))
	map2 := make(map[rune]int, len(word2))
	for _, r := range word1 {
		map1[r]++
	}
	for _, r := range word2 {
		map2[r]++
	}
	for k, v1 := range map1 {
		if v2, ok := map2[k]; !ok || v2 != v1 {
			return false
		}
		delete(map2, k)
	}
	return len(map2) == 0
}

func isMapped(anagrams []string, word string) bool {
	for _, anagram := range anagrams {
		if anagram == word {
			return true
		}
	}
	return false
}
