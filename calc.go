package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"calc/calculator"
)

func main() {
	var expression string

	// если аргументы переданы - считаем выражение из них
	if len(os.Args) > 1 {
		expression = strings.Join(os.Args[1:], "")
	} else {
		// иначе читаем одну строку из stdin
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "Ошибка: пустой ввод")
			os.Exit(1)
		}

		expression = scanner.Text()
	}

	result, err := calculator.Calculate(expression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
