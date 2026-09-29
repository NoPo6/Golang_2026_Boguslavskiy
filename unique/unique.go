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
func normalize(line string, commands Commands) string {
	resultLine := []rune(line) // переменная для результата (нормализованная строка)

	pointer := 0      // указатель (текущий индекс строки)
	fieldCounter := 0 // счетчик полей

	// цикл удаления полей
	for pointer < len(resultLine) {
		// если удалили нужное кол-во полей, то выходим из цикла
		if fieldCounter == commands.SkipFields {
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
	if commands.SkipFields > 0 && pointer < len(resultLine) && resultLine[pointer] == ' ' {
		pointer++
	}

	charCounter := 0 // счетчик символов (для удаления)
	// цикл пропуска символов
	for (pointer < len(resultLine)) && (charCounter < commands.SkipChars) {
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

	if commands.IgnoreCase {
		result = strings.ToLower(result)
	}

	return result
}

// основная функция unique
func Unique(commands Commands, lines []string) []string {
	// группа подряд идущих одинаковых строк
	type group struct {
		line       string // исходная строка (первая в группе)
		normalized string // нормализованный вид для сравнения
		count      int    // сколько раз повторилась
	}

	groups := make([]group, 0, len(lines))

	// группируем подряд идущие строки по нормализованному виду
	for _, line := range lines {
		normalized := normalize(line, commands)

		// если группа есть и совпадает с предыдущей - увеличиваем счетчик
		if len(groups) > 0 && groups[len(groups)-1].normalized == normalized {
			groups[len(groups)-1].count++
			continue
		}

		// иначе начинаем новую группу
		groups = append(groups, group{line: line, normalized: normalized, count: 1})
	}

	result := make([]string, 0, len(groups)) // для записи результата выполнения функции

	// применяем флаги к группам
	for _, g := range groups {
		if commands.Count {
			// -c: число повторений + пробел + исходная строка
			result = append(result, strconv.Itoa(g.count)+" "+g.line)
			continue
		}

		if commands.Duplicates {
			// -d: только группы, где было больше одной строки
			if g.count > 1 {
				result = append(result, g.line)
			}
			continue
		}

		if commands.Unique {
			// -u: только группы, где строка встретилась один раз
			if g.count == 1 {
				result = append(result, g.line)
			}
			continue
		}

		// поведение по умолчанию - по одному представителю каждой группы
		result = append(result, g.line)
	}

	return result
}
