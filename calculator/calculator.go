package calculator

import (
	"errors"
	"strconv"
	"strings"
)

// ошибки парсера
var (
	ErrUnexpectedChar  = errors.New("неожиданный символ")
	ErrUnexpectedEnd   = errors.New("неожиданный конец выражения")
	ErrDivisionByZero  = errors.New("деление на ноль")
	ErrEmptyExpression = errors.New("пустое выражение")
	ErrMismatchedParen = errors.New("несогласованные скобки")
	ErrTooDeep         = errors.New("слишком большая вложенность")
)

const maxDepth = 1000

// структура для разбора выражения
type parser struct {
	input string
	pos   int
	depth int
}

// функция для вычисления выражения
func Calculate(expression string) (float64, error) {
	// убираем пробелы, чтобы не отвлекаться на них при разборе
	cleaned := strings.ReplaceAll(expression, " ", "")

	// пустое выражение считать нечего
	if cleaned == "" {
		return 0, ErrEmptyExpression
	}

	p := &parser{input: cleaned, pos: 0}

	// считаем самое верхнеуровневое выражение
	result, err := p.parseSum()
	if err != nil {
		return 0, err
	}

	// если после разбора остались символы - выражение некорректно
	if p.pos != len(p.input) {
		return 0, ErrUnexpectedChar
	}

	return result, nil
}

// функция для разбора суммы и разности
func (p *parser) parseSum() (float64, error) {
	// разбираем первый множитель
	left, err := p.parseProduct()
	if err != nil {
		return 0, err
	}

	// цикл разбора цепочки + и -
	for p.pos < len(p.input) {
		op := p.input[p.pos]

		// если это не + и не -, то заканчиваем разбор выражения
		if op != '+' && op != '-' {
			break
		}

		p.pos++

		// после бинарного + или - не может идти ещё один + или -
		if p.pos < len(p.input) && (p.input[p.pos] == '+' || p.input[p.pos] == '-') {
			return 0, ErrUnexpectedChar
		}

		// разбираем следующую операцию
		right, err := p.parseProduct()
		if err != nil {
			return 0, err
		}

		if op == '+' {
			left += right
		} else {
			left -= right
		}
	}

	return left, nil
}

// функция для разбора произведения и деления
func (p *parser) parseProduct() (float64, error) {
	// разбираем первый унарный операнд
	left, err := p.parseUnary()
	if err != nil {
		return 0, err
	}

	// цикл разбора цепочки *, / и неявного умножения
	for p.pos < len(p.input) {
		op := p.input[p.pos]

		// случай явного умножения или деления
		if op == '*' || op == '/' {
			p.pos++

			// после бинарного * или / не может идти + или -
			if p.pos < len(p.input) && (p.input[p.pos] == '+' || p.input[p.pos] == '-') {
				return 0, ErrUnexpectedChar
			}

			right, err := p.parseUnary()
			if err != nil {
				return 0, err
			}

			if op == '*' {
				left *= right
			} else {
				// деление на ноль недопустимо
				if right == 0 {
					return 0, ErrDivisionByZero
				}
				left /= right
			}
			continue
		}

		// случай неявного умножения: (9+3)(7+2), 2(3), (3)2
		// если следующий символ открывает новую группу - это умножение
		if op == '(' || (op >= '0' && op <= '9') {
			right, err := p.parseUnary()
			if err != nil {
				return 0, err
			}
			left *= right
			continue
		}

		break
	}

	return left, nil
}

// функция для разбора унарного минуса
func (p *parser) parseUnary() (float64, error) {
	// берём не более одного минуса
	negative := false
	if p.pos < len(p.input) && p.input[p.pos] == '-' {
		negative = true
		p.pos++
	}

	// разбираем операнд после минуса
	value, err := p.parsePrimary()
	if err != nil {
		return 0, err
	}

	// меняем знак
	if negative {
		value = -value
	}

	return value, nil
}

// функция для разбора первичного операнда (число или скобка)
func (p *parser) parsePrimary() (float64, error) {
	// проверяем конец выражения
	if p.pos >= len(p.input) {
		return 0, ErrUnexpectedEnd
	}

	ch := p.input[p.pos]

	// случай скобки
	if ch == '(' {
		p.depth++
		// проверяем глубину рекурсии
		if p.depth > maxDepth {
			return 0, ErrTooDeep
		}

		p.pos++

		// разбираем выражение внутри скобок
		value, err := p.parseSum()

		p.depth--

		if err != nil {
			return 0, err
		}

		// ожидаем закрывающую скобку
		if p.pos >= len(p.input) || p.input[p.pos] != ')' {
			return 0, ErrMismatchedParen
		}

		p.pos++
		return value, nil
	}

	// встретили число
	if ch >= '0' && ch <= '9' {
		start := p.pos
		dotSeen := false

		// считываем все цифры и точки подряд
		for p.pos < len(p.input) {
			c := p.input[p.pos]

			if c >= '0' && c <= '9' {
				p.pos++
				continue
			}

			// встретили точку
			if c == '.' {
				// вторая точка в одном числе - сразу ошибка
				if dotSeen {
					return 0, ErrUnexpectedChar
				}
				dotSeen = true
				p.pos++
				continue
			}

			break
		}

		// превращаем в число с плавающей точкой
		value, err := strconv.ParseFloat(p.input[start:p.pos], 64)
		if err != nil {
			return 0, err
		}

		return value, nil
	}

	return 0, ErrUnexpectedChar
}
