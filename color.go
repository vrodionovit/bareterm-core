package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// ANSI цвета
var ansiColors = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{170, 0, 0, 255},     // Red
	{0, 170, 0, 255},     // Green
	{170, 85, 0, 255},    // Yellow
	{0, 0, 170, 255},     // Blue
	{170, 0, 170, 255},   // Magenta
	{0, 170, 170, 255},   // Cyan
	{170, 170, 170, 255}, // White
}

// Яркие ANSI цвета
var ansiBrightColors = []color.RGBA{
	{85, 85, 85, 255},    // Bright Black
	{255, 85, 85, 255},   // Bright Red
	{85, 255, 85, 255},   // Bright Green
	{255, 255, 85, 255},  // Bright Yellow
	{85, 85, 255, 255},   // Bright Blue
	{255, 85, 255, 255},  // Bright Magenta
	{85, 255, 255, 255},  // Bright Cyan
	{255, 255, 255, 255}, // Bright White
}

type ColorState struct {
	Foreground color.RGBA
	Background color.RGBA
}

func (t *Terminal) setForegroundColor(colorIndex int) {
	if colorIndex < 0 || colorIndex > 15 {
		fmt.Printf("Неверный индекс цвета переднего плана: %d\n", colorIndex)
		return
	}

	var newColor color.RGBA
	if colorIndex < 8 {
		newColor = ansiColors[colorIndex]
	} else {
		newColor = ansiBrightColors[colorIndex-8]
	}

	t.ColorState.Foreground = newColor
	fmt.Printf("Установлен цвет переднего плана: %v\n", newColor)
	t.updateTerminalColors()
}

func (t *Terminal) setBackgroundColor(colorIndex int) {
	if colorIndex < 0 || colorIndex > 15 {
		fmt.Printf("Неверный индекс цвета фона: %d\n", colorIndex)
		return
	}

	var newColor color.RGBA
	if colorIndex < 8 {
		newColor = ansiColors[colorIndex]
	} else {
		newColor = ansiBrightColors[colorIndex-8]
	}

	t.ColorState.Background = newColor
	fmt.Printf("Установлен цвет фона: %v\n", newColor)
	t.updateTerminalColors()
}

func (t *Terminal) resetForegroundColor() {
	t.ColorState.Foreground = ansiColors[7] // White
	fmt.Println("Сброс цвета переднего плана к значению по умолчанию (White)")
	t.updateTerminalColors()
}

func (t *Terminal) resetBackgroundColor() {
	t.ColorState.Background = ansiColors[0] // Black
	fmt.Println("Сброс цвета фона к значению по умолчанию (Black)")
	t.updateTerminalColors()
}

func (t *Terminal) updateTerminalColors() {
	// Здесь должен быть код для фактического обновления цветов в вашем терминале
	// Например, если вы используете библиотеку для рендеринга терминала, вы можете вызвать соответствующий метод здесь
	fmt.Printf("Обновление цветов терминала: Передний план %v, Фон %v\n", t.ColorState.Foreground, t.ColorState.Background)
}

func (t *Terminal) handleColor(params string) {
	if params == "" {
		params = "0" // Если параметры не указаны, сбрасываем все атрибуты
	}

	attributes := strings.Split(params, ";")
	for _, attr := range attributes {
		value, err := strconv.Atoi(attr)
		if err != nil {
			fmt.Printf("Неверный параметр SGR: %s\n", attr)
			continue
		}

		switch {
		case value == 0:
			t.resetGraphicsMode()
		case value == 1:
			t.setBold(true)
		case value == 2:
			t.setDim(true)
		case value == 3:
			t.setItalic(true)
		case value == 4:
			t.setUnderline(true)
		case value == 5:
			t.setBlink(true)
		case value == 7:
			t.setReverse(true)
		case value == 8:
			t.setHidden(true)
		case value == 9:
			t.setStrikethrough(true)
		case value >= 30 && value <= 37:
			t.setForegroundColor(value - 30)
		case value == 39:
			t.resetForegroundColor()
		case value >= 40 && value <= 47:
			t.setBackgroundColor(value - 40)
		case value == 49:
			t.resetBackgroundColor()
		case value == 38:
			t.setExtendedForegroundColor(attributes)
		case value == 48:
			t.setExtendedBackgroundColor(attributes)
		case value >= 90 && value <= 97:
			t.setForegroundColor(value - 90 + 8) // Яркие цвета
		case value >= 100 && value <= 107:
			t.setBackgroundColor(value - 100 + 8) // Яркие цвета
		default:
			fmt.Printf("Неподдерживаемый параметр SGR: %d\n", value)
		}
	}
}

func (t *Terminal) resetGraphicsMode() {
	fmt.Println("Сброс всех атрибутов графического режима")
	// Реализация сброса всех атрибутов
}

func (t *Terminal) setBold(on bool) {
	fmt.Printf("Установка жирного шрифта: %v\n", on)
	// Реализация установки жирного шрифта
}

func (t *Terminal) setDim(on bool) {
	fmt.Printf("Установка тусклого шрифта: %v\n", on)
	// Реализация установки тусклого шрифта
}

func (t *Terminal) setItalic(on bool) {
	fmt.Printf("Установка курсива: %v\n", on)
	// Реализация установки курсива
}

func (t *Terminal) setUnderline(on bool) {
	fmt.Printf("Установка подчеркивания: %v\n", on)
	// Реализация установки подчеркивания
}

func (t *Terminal) setBlink(on bool) {
	fmt.Printf("Установка мигания: %v\n", on)
	// Реализация установки мигания
}

func (t *Terminal) setReverse(on bool) {
	fmt.Printf("Установка обратного видео: %v\n", on)
	// Реализация установки обратного видео
}

func (t *Terminal) setHidden(on bool) {
	fmt.Printf("Установка скрытого текста: %v\n", on)
	// Реализация установки скрытого текста
}

func (t *Terminal) setStrikethrough(on bool) {
	fmt.Printf("Установка зачеркивания: %v\n", on)
	// Реализация установки зачеркивания
}

func (t *Terminal) setExtendedForegroundColor(params []string) {
	fmt.Printf("Установка расширенного цвета переднего плана: %v\n", params)
	// Реализация установки расширенного цвета переднего плана
}

func (t *Terminal) setExtendedBackgroundColor(params []string) {
	fmt.Printf("Установка расширенного цвета фона: %v\n", params)
	// Реализация установки расширенного цвета фона
}
