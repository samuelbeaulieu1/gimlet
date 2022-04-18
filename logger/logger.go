package logger

import (
	"fmt"
	"os"
	"strings"
)

const (
	Black = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func SetTextColor(tag string, color int) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, tag)
}

func PrintDebug(format string, values ...any) {
	format = addNewLine(format)
	fmt.Fprintf(os.Stdout, SetTextColor("[Gimlet debug] ", Blue)+format, values...)
}

func PrintInfo(format string, values ...any) {
	format = addNewLine(format)
	fmt.Fprintf(os.Stdout, SetTextColor("[Gimlet info] ", Yellow)+format, values...)
}

func PrintError(format string, values ...any) {
	format = addNewLine(format)
	fmt.Fprintf(os.Stdout, SetTextColor("[Gimlet error] ", Red)+format, values...)
}

func addNewLine(format string) string {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	return format
}
