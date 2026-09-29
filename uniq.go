package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"uniq/unique"
)

// структура, которая говорит, переданы ли файловые аргументы
type FileFlags struct {
	HasInput  bool
	HasOutput bool
}

// функция для чтения всех строк из произвольного io.Reader
func readLinesFromFile(fileName io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(fileName)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// функция для записи строк в произвольный io.Writer
func writeLinesToFile(writer io.Writer, lines []string) error {
	bufferedWriter := bufio.NewWriter(writer)
	defer bufferedWriter.Flush()

	// пишем построчно, а не через fmt.Println(output), чтобы не получить "[a b c]"
	for _, line := range lines {
		if _, err := fmt.Fprintln(bufferedWriter, line); err != nil {
			return err
		}
	}

	return nil
}

// функция для проверки количества файловых аргументов
func validateFiles(numberOfFiles int) (FileFlags, error) {
	switch numberOfFiles {
	case 2:
		return FileFlags{HasInput: true, HasOutput: true}, nil
	case 1:
		return FileFlags{HasInput: true, HasOutput: false}, nil
	case 0:
		return FileFlags{HasInput: false, HasOutput: false}, nil
	default:
		err := errors.New("Неверное количество файловых аргументов!")
		return FileFlags{}, err
	}
}

// функция для вывода правильного использования утилиты
func printUsage() {
	fmt.Fprintln(os.Stderr, "Использование: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
}

func main() {
	log.SetFlags(0)
	// объявляем флаги
	count := flag.Bool("c", false, "подсчитать количество встречаний строки во входных данных")
	duplicates := flag.Bool("d", false, "вывести только те строки, которые повторились во входных данных")
	uniqueFlag := flag.Bool("u", false, "вывести только те строки, которые не повторились во входных данных")
	ignoreCase := flag.Bool("i", false, "не учитывать регистр букв")
	skipFields := flag.Int("f", 0, "не учитывать первые num_fields полей в строке")
	skipChars := flag.Int("s", 0, "не учитывать первые num_chars символов в строке")

	// подменяем стандартный usage на собственный
	flag.Usage = printUsage

	// парсим флаги
	flag.Parse()

	// проверяем взаимоисключение флагов -c, -d, -u
	// делаем это именно ПОСЛЕ flag.Parse(), иначе значения флагов ещё дефолтные
	if (*count && *duplicates) || (*count && *uniqueFlag) || (*duplicates && *uniqueFlag) {
		printUsage()
		log.Fatalf("Ошибка: флаги -c, -d, -u взаимоисключающие")
	}

	// проверяем количество файловых аргументов
	fileFlags, err := validateFiles(flag.NArg())
	if err != nil {
		printUsage()
		log.Fatalf("Ошибка: %v", err)
	}

	// готовим входной поток: либо файл, либо stdin
	var inputFile *os.File
	if fileFlags.HasInput {
		inputFile, err = os.Open(flag.Args()[0])
		if err != nil {
			log.Fatalf("Ошибка: %v", err)
		}
		defer inputFile.Close()
	} else {
		inputFile = os.Stdin
	}

	// читаем все строки из входного потока
	lines, err := readLinesFromFile(inputFile)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	// готовим выходной поток: либо файл, либо stdout
	var outputFile *os.File
	if fileFlags.HasOutput {
		// используем именно flag.Args()[1], а не "output.txt"
		outputFile, err = os.Create(flag.Args()[1])
		if err != nil {
			log.Fatalf("Ошибка: %v", err)
		}
		defer outputFile.Close()
	} else {
		outputFile = os.Stdout
	}

	// формируем структуру команд для пакета unique
	commands := unique.Commands{
		Count:      *count,
		Duplicates: *duplicates,
		Unique:     *uniqueFlag,
		IgnoreCase: *ignoreCase,
		SkipFields: *skipFields,
		SkipChars:  *skipChars,
	}

	// вызываем основную логику уникализации
	output := unique.Unique(commands, lines)

	// пишем результат в выходной поток
	if err := writeLinesToFile(outputFile, output); err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
}
