package unique

import (
	"strconv"
	"strings"
)

// структура с флагами
type Commands struct {
	Count      bool // -c
	Duplicates bool // -d
	Unique     bool // -u
	IgnoreCase bool // -i
	SkipFields int  // -f
	SkipChars  int  // -s
}

// функция для удаления префикса из строки и приведения к нижнему регистру (если необходимо)
func normalize(line string, ignoreCase bool, skipFields int, skipChars int) string {

	resultLine := []rune(line) // переменная для результата (нормализованная строка)

	pointer := 0      // указатель (текущий индекс строки)
	fieldCounter := 0 // счетчик полей

	// цикл удаления полей
	for pointer < len(resultLine) {
		// если удалили нужное кол-во полей, то выходим из цикла
		if fieldCounter == skipFields {
			break
		}

		// пропускаем все пробелы
		for (pointer < len(resultLine)) && (resultLine[pointer] == ' ') {
			pointer++
		}

		// пропускаем одно поле (набор непробельных символов подряд)
		for (pointer < len(resultLine)) && (resultLine[pointer] != ' ') {
			pointer++
		}

		// инкрементируем счетчик полей
		fieldCounter += 1
	}

	// убираем разделяющий пробел после последнего поля
	if skipFields > 0 && pointer < len(resultLine) && resultLine[pointer] == ' ' {
		pointer++
	}

	charCounter := 0 // счетчик символов (для удаления)
	// цикл пропуска символов
	for (pointer < len(resultLine)) && (charCounter < skipChars) {
		pointer++
		charCounter++
	}

	// к этому моменту указатель пропустил все ненужные символы
	// и если строка не закончилась, то вырезаем из нее полезную часть
	var result string
	if pointer < len(resultLine) {
		result = string(resultLine[pointer:])
	} else {
		result = ""
	}

	if ignoreCase {
		result = strings.ToLower(result)
	}

	return result
}

// основная функция unique
func Unique(commands Commands, lines []string) []string {
	var result []string // для записи результата выполнения функции

	// проверяем на соответствие определенному флагу
	if commands.Count {
		// условие с подсчетом повторяющихся строк
		lineCounter := 1   // счетчик слов
		previousLine := "" // начальное значение предыдущей строки

		// цикл перебора всех строк
		for index, line := range lines {

			// нормализация строки
			normalizedLine := normalize(line, commands.IgnoreCase, commands.SkipFields, commands.SkipChars)

			// условие для начального элемента с нулевым индексом
			if index == 0 {
				// встретили новое значение сразу
				result = append(result, line) // запоминаем только первую строку, т.к. разные строки после нормализации могут быть похожи
				previousLine = normalizedLine // запоминаем нормализованную форму, т.к. сравниваем по ним
				continue
			}

			if previousLine == normalizedLine {
				// встретили похожее значение
				lineCounter++ // инкрементируем счетчик похожих строк
			} else {
				// встретили новое значение после подсчета
				result[len(result)-1] = strconv.Itoa(lineCounter) + " " + result[len(result)-1] // составляем строку "число похожих строк + сама строка"
				result = append(result, line)                                                   // запоминаем новую строку
				lineCounter = 1                                                                 // возвращаем счетчик
				previousLine = normalizedLine                                                   // запоминаем новое нормализованное значение
			}
		}

		if len(result) > 0 {
			result[len(result)-1] = strconv.Itoa(lineCounter) + " " + result[len(result)-1]
		}

	} else if commands.Duplicates {
		// условие с выводом только продублированных строк

		lineCounter := 1   // счетчик продублированных строк
		previousLine := "" // начальное значение предыдущей строки
		curentLine := ""   // текущее начальное значение продублированной последовательности

		for index, line := range lines {
			normalizedLine := normalize(line, commands.IgnoreCase, commands.SkipFields, commands.SkipChars)

			if index == 0 {
				curentLine = line
				previousLine = normalizedLine
				continue
			}

			if normalizedLine == previousLine {
				lineCounter++
				continue
			}

			if lineCounter > 1 {
				result = append(result, curentLine)
			}

			previousLine = normalizedLine
			lineCounter = 1
			curentLine = line
		}

		if len(lines) > 0 && lineCounter > 1 {
			result = append(result, curentLine)
		}
	} else if commands.Unique {
		// условие с выводом только неповторяющихся строк (флаг -u)

		lineCounter := 1   // счетчик продублированных строк
		previousLine := "" // начальное значение предыдущей строки
		curentLine := ""   // текущее начальное значение продублированной последовательности

		for index, line := range lines {
			normalizedLine := normalize(line, commands.IgnoreCase, commands.SkipFields, commands.SkipChars)

			if index == 0 {
				curentLine = line
				previousLine = normalizedLine
				continue
			}

			if normalizedLine == previousLine {
				lineCounter++
				continue
			}

			if lineCounter == 1 {
				result = append(result, curentLine)
			}

			previousLine = normalizedLine
			lineCounter = 1
			curentLine = line
		}

		if len(lines) > 0 && lineCounter == 1 {
			result = append(result, curentLine)
		}
	} else {
		// поведение по умолчанию - выводим по одному представителю каждой группы подряд идущих строк

		previousLine := "" // начальное значение предыдущей строки
		curentLine := ""   // текущее начальное значение последовательности

		for index, line := range lines {
			normalizedLine := normalize(line, commands.IgnoreCase, commands.SkipFields, commands.SkipChars)

			if index == 0 {
				curentLine = line
				previousLine = normalizedLine
				continue
			}

			if normalizedLine == previousLine {
				// строка совпала с предыдущей - пропускаем
				continue
			}

			// встретили новую группу - сохраняем представителя предыдущей
			result = append(result, curentLine)
			previousLine = normalizedLine
			curentLine = line
		}

		// не забываем сохранить последнюю группу
		if len(lines) > 0 {
			result = append(result, curentLine)
		}
	}

	return result
}
