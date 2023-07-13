Что выведет программа? Объяснить вывод программы.

```go
package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}
```

Ответ:
```
Числа от 0 до 9, затем deadlock, т.к. канал ch не закрыт и запись чтения из пустого канала заблокирует main горутину

```
