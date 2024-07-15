package main

import "os"

// main функция для теста вызова os.Exit в методе Exit
func main() {
	println("here is it")
	t := Exiter{}
	t.Exit(1)
}

// Exiter структура для формирования теста
type Exiter struct {
}

// Exit метод в который оборачиваем вызов os.Exit
func (e Exiter) Exit(code int) {
	os.Exit(code)
}
