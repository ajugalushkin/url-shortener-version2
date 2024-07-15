package main

import "os"

// main функция для теста вызова os.Exit в функции main
func main() {
	println("here is it")
	os.Exit(1) // want `os.Exit called in main func in main package`
}
