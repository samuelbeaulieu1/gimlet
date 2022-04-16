package logger

import (
	"fmt"
	"os"
	"strings"
)

func PrintDebug(format string, values ...any) {
	format = addNewLine(format)
	fmt.Fprintf(os.Stdout, "[debug] "+format, values...)
}

func PrintInfo(format string, values ...any) {
	format = addNewLine(format)
	fmt.Fprintf(os.Stdout, "[info] "+format, values...)
}

func PrintError(format string, values ...any) {
	format = addNewLine(format)
	fmt.Fprintf(os.Stdout, "[error] "+format, values...)
}

func addNewLine(format string) string {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	return format
}
