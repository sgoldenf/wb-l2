Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


	func test() (x int) {
		defer func() {
			x++
		}()
		x = 1
		return
	}


	func anotherTest() int {
		var x int
		defer func() {
			x++
		}()
		x = 1
		return x
	}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
2
1

defer добавляет вызов функции в стеке -> выполняются от последнего к первому
defer вычисляет аргументы в том месте, где его вызвали
в функции test увеличивается именованное возвращаемое значение
в функции AnotherTest увеличивается локальная переменная на стеке
```
