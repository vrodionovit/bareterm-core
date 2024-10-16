package main

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// EncodingMode - перечисление для различных режимов кодировки
type EncodingMode int

const (
	EncodingASCII       EncodingMode = iota // ASCII кодировка
	EncodingUTF8                            // UTF-8 кодировка
	EncodingISO8859_1                       // ISO-8859-1 кодировка
	EncodingWindows1251                     // Windows-1251 кодировка
	// Добавьте другие кодировки по необходимости
)

// SetEncoding устанавливает режим кодировки для терминала
func (t *Terminal) SetEncoding(mode EncodingMode) error {
	t.currentEncoding = mode
	var err error
	switch mode {
	case EncodingASCII, EncodingUTF8:
		t.decoder = nil // Для ASCII и UTF-8 используем встроенную поддержку Go
	case EncodingISO8859_1:
		t.decoder = charmap.ISO8859_1.NewDecoder() // Создаем декодер для ISO-8859-1
	case EncodingWindows1251:
		t.decoder = charmap.Windows1251.NewDecoder() // Создаем декодер для Windows-1251
	default:
		return fmt.Errorf("unsupported encoding mode: %d", mode) // Возвращаем ошибку для неподдерживаемых режимов
	}
	return err
}

// DecodeInput декодирует входные данные в соответствии с текущей кодировкой
func (t *Terminal) DecodeInput(input []byte) (string, error) {
	switch t.currentEncoding {
	case EncodingASCII:
		return t.decodeASCII(input), nil // Декодируем ASCII
	case EncodingUTF8:
		return string(input), nil // UTF-8 не требует дополнительного декодирования
	default:
		if t.decoder == nil {
			return "", fmt.Errorf("decoder not initialized for encoding mode: %d", t.currentEncoding)
		}
		// Используем декодер для преобразования входных данных
		reader := transform.NewReader(bytes.NewReader(input), t.decoder)
		decoded, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}
}

// decodeASCII декодирует ASCII-данные, заменяя не-ASCII символы на '?'
func (t *Terminal) decodeASCII(input []byte) string {
	result := make([]rune, 0, len(input))
	for _, b := range input {
		if b < 128 {
			result = append(result, rune(b)) // ASCII-символы добавляем как есть
		} else {
			result = append(result, '?') // Заменяем не-ASCII символы на '?'
		}
	}
	return string(result)
}

// decodeSingleChar декодирует один символ в соответствии с текущей кодировкой
func (t *Terminal) decodeSingleChar(buf []byte) (rune, int, error) {
	switch t.currentEncoding {
	case EncodingASCII:
		if buf[0] < 128 {
			return rune(buf[0]), 1, nil // Возвращаем ASCII-символ
		}
		return '?', 1, nil // Заменяем не-ASCII символ на '?'
	case EncodingUTF8:
		r, size := utf8.DecodeRune(buf) // Декодируем UTF-8 руну
		return r, size, nil
	default:
		if t.decoder == nil {
			return 0, 0, fmt.Errorf("decoder not initialized for encoding mode: %d", t.currentEncoding)
		}
		// Используем декодер для преобразования одного байта
		decoded, n, err := transform.Bytes(t.decoder, buf[:1])
		if err != nil {
			return 0, 0, err
		}
		r, _ := utf8.DecodeRune(decoded) // Декодируем результат в руну
		return r, n, nil
	}
}
