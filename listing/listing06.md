Что выведет программа? Объяснить вывод программы. Рассказать про внутреннее устройство слайсов и что происходит при передачи их в качестве аргументов функции.

```go
package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
```

Ответ:
```
[3 2 3]

В main создается слайс с размером и емкостью 3
В modifySlice сначала перезаписывается 0-й элемент
Затем происходит append, т.к. size==cap, то выделяется новая область в памяти
После этого i и s указывают на разные участки памяти и изменения 1-го элемента не несут никакого эффекта для слайса s
```
