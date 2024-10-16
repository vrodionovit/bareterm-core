package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/creack/pty"
	"golang.org/x/text/encoding"
)

const (
	stateNormal = iota
	stateEscape
	stateCSI
	stateOSC
)

type Terminal struct {
	parseState       int
	escapeBuffer     []byte
	cursorX, cursorY int
	width, height    int
	ColorState       ColorState
	currentEncoding  EncodingMode
	decoder          *encoding.Decoder
}

func NewTerminal(width, height int) *Terminal {
	return &Terminal{
		width:  width,
		height: height,
		ColorState: ColorState{
			Foreground: ansiColors[7], // White
			Background: ansiColors[0], // Black
		},
	}
}

func (t *Terminal) handleOutput(buf []byte) {
	for i := 0; i < len(buf); {
		switch t.parseState {
		case stateNormal:
			if buf[i] == 0x1B { // ESC
				t.parseState = stateEscape
				t.escapeBuffer = []byte{buf[i]}
				i++
			} else {
				// Декодируем только обычные символы
				r, size, err := t.decodeSingleChar(buf[i:])
				if err != nil {
					fmt.Printf("Error decoding character at position %d: %v\n", i, err)
					i++
					continue
				}
				t.handleOutputChar(r)
				i += size
			}
		case stateEscape:
			t.escapeBuffer = append(t.escapeBuffer, buf[i])
			switch buf[i] {
			case '[':
				t.parseState = stateCSI
			case ']':
				t.parseState = stateOSC
			default:
				t.handleEscape(string(t.escapeBuffer))
				t.parseState = stateNormal
				t.escapeBuffer = nil
			}
			i++
		case stateCSI:
			t.escapeBuffer = append(t.escapeBuffer, buf[i])
			if (buf[i] >= 0x40 && buf[i] <= 0x7E) && buf[i] != '[' {
				t.handleCSI(string(t.escapeBuffer[2:])) // Skip ESC[
				t.parseState = stateNormal
				t.escapeBuffer = nil
			}
			i++
		case stateOSC:
			t.escapeBuffer = append(t.escapeBuffer, buf[i])
			if buf[i] == 0x07 || (buf[i] == '\\' && t.escapeBuffer[len(t.escapeBuffer)-2] == 0x1B) {
				t.handleOSC(string(t.escapeBuffer))
				t.parseState = stateNormal
				t.escapeBuffer = nil
			}
			i++
		}
	}
}

func (t *Terminal) handleOutputChar(r rune) {
	fmt.Printf("%c", r)
	t.cursorX++
	if t.cursorX >= t.width {
		t.cursorX = 0
		t.cursorY++
	}
	if t.cursorY >= t.height {
		// Implement scrolling here
	}
}

func (t *Terminal) handleOSC(sequence string) {
	// Удаляем начальный ESC] и конечный BEL или ST
	sequence = strings.TrimPrefix(sequence, "\x1B]")
	sequence = strings.TrimSuffix(sequence, "\x07")
	sequence = strings.TrimSuffix(sequence, "\x1B\\")

	// Разделяем команду и параметры
	parts := strings.SplitN(sequence, ";", 2)
	if len(parts) < 2 {
		fmt.Printf("Неверная OSC последовательность: %q\n", sequence)
		return
	}

	command := parts[0]
	param := parts[1]

	switch command {
	case "0", "1", "2":
		t.setWindowTitle(param)
	case "4":
		t.setColorPalette(param)
	case "10", "11", "12", "13", "14", "15", "16", "17":
		t.setDynamicColor(command, param)
	case "52":
		t.manipulateSelectionData(param)
	default:
		fmt.Printf("Неизвестная OSC команда: %s\n", command)
	}
}

func (t *Terminal) handleCSI(sequence string) {
	fmt.Printf("последовательность: %s\n", sequence)

	switch {
	case strings.HasSuffix(sequence, "H") || strings.HasSuffix(sequence, "f"):
		t.moveCursor(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "A"):
		t.moveCursorUp(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "B"):
		t.moveCursorDown(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "C"):
		t.moveCursorForward(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "D"):
		t.moveCursorBackward(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "m"):
		t.setGraphicsMode(sequence[:len(sequence)-1])
	case sequence == "2J":
		t.clearScreen()
	case strings.HasSuffix(sequence, "K"):
		t.eraseLine(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "S"):
		t.scrollUp(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "T"):
		t.scrollDown(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "n"):
		t.deviceStatusReport(sequence[:len(sequence)-1])
	case strings.HasSuffix(sequence, "h") || strings.HasSuffix(sequence, "l"):
		t.setMode(sequence, strings.HasSuffix(sequence, "h"))
	case strings.HasSuffix(sequence, "r"):
		t.setScrollingRegion(sequence[:len(sequence)-1])
	default:
		fmt.Printf("Неизвестная CSI последовательность: %s\n", sequence)
	}
}

func getAvailableShells() []string {
	commonShells := []string{"bash", "sh", "zsh", "fish", "powershell", "pwsh"}
	availableShells := []string{}

	for _, shell := range commonShells {
		path, err := exec.LookPath(shell)
		if err == nil {
			availableShells = append(availableShells, path)
		}
	}

	if runtime.GOOS != "windows" {
		// Проверяем /etc/shells на Unix-подобных системах
		if shellsFile, err := os.ReadFile("/etc/shells"); err == nil {
			for _, line := range strings.Split(string(shellsFile), "\n") {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					if _, err := os.Stat(line); err == nil {
						availableShells = append(availableShells, line)
					}
				}
			}
		}
	} else {
		// Дополнительные проверки для Windows
		programFiles := os.Getenv("ProgramFiles")
		possiblePaths := []string{
			filepath.Join(programFiles, "PowerShell", "7", "pwsh.exe"),
			filepath.Join(programFiles, "PowerShell", "6", "pwsh.exe"),
			filepath.Join(os.Getenv("WINDIR"), "System32", "WindowsPowerShell", "v1.0", "powershell.exe"),
		}
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				availableShells = append(availableShells, path)
			}
		}
	}

	// Удаляем дубликаты
	uniqueShells := []string{}
	seen := make(map[string]bool)
	for _, shell := range availableShells {
		if _, ok := seen[shell]; !ok {
			seen[shell] = true
			uniqueShells = append(uniqueShells, shell)
		}
	}

	return uniqueShells
}

func selectShell() string {
	shells := getAvailableShells()
	fmt.Println("Доступные оболочки:")
	for i, shell := range shells {
		fmt.Printf("%d. %s\n", i+1, shell)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Выберите оболочку (введите номер): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if num, err := strconv.Atoi(input); err == nil && num > 0 && num <= len(shells) {
			return shells[num-1]
		}

		fmt.Println("Неверный ввод. Пожалуйста, введите номер из списка.")
	}
}

// func handleKeyPress(key rune, ptmx *os.File) {
// 	switch key {
// 	case '\t':
// 		// Отправляем Tab в pty
// 		ptmx.Write([]byte{'\t'})
// 		// Обработка других клавиш...
// 	}
// }

func main() {
	selectedShell := selectShell()
	cmd := exec.Command(selectedShell)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer ptmx.Close()

	term := NewTerminal(80, 24)
	term.SetEncoding(EncodingUTF8)

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := ptmx.Read(buf)
			if err != nil {
				if err == io.EOF {
					return
				}
				panic(err)
			}
			fmt.Printf("Raw output: %q\n", buf[:n])
			term.handleOutput(buf[:n])
		}
	}()

	io.Copy(ptmx, os.Stdin)
}
