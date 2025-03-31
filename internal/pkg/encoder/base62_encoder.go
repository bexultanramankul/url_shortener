package encoder

import (
	"bytes"
)

const base62Characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base62Length = int64(len(base62Characters))

// Encode преобразует число в Base62 строку длиной 6 символов
func Encode(number int64) string {
	if number < 0 {
		return "000000" // Можно заменить на ошибку
	}

	// Используем bytes.Buffer для быстрого конкатенирования строк
	var base62 bytes.Buffer
	base62.Grow(6) // Предварительное выделение памяти

	for number > 0 {
		base62.WriteByte(base62Characters[number%base62Length])
		number /= base62Length
	}

	// Дополним строку нулями слева
	result := padLeft(base62.String(), '0', 6)

	// Переворачиваем только финальный результат, а не всю строку в процессе
	return reverseString(result)
}

// reverseString реверсирует строку (быстрее без дополнительных копирований)
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// padLeft дополняет строку слева указанным символом до нужной длины
func padLeft(s string, pad rune, length int) string {
	for len(s) < length {
		s = string(pad) + s
	}
	return s
}
