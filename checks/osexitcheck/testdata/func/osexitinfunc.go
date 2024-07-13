package main

import (
	"os"
)

// main функция для теста
func main() {
	println("here is it")
	Exit(1)
}

// Exit проверяем вызов os.Exit из функции Exit
func Exit(code int) {
	os.Exit(code)
}
