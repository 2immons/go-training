package main

import (
	"flag"
	"fmt"
)

func main() {
	// Определение флагов
	src := flag.String("src", "", "Source file path")
	dst := flag.String("dst", "", "Destination file path")

	// Парсинг командной строки
	flag.Parse()

	// Вывод значений флагов
	fmt.Println("Source:", *src)
	fmt.Println("Destination:", *dst)
}
