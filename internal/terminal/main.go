package terminal

import (
	tm "github.com/buger/goterm"
)

func Clear() {
	tm.Clear()
}

func Append(line string) {
	tm.Println(line)
	tm.Flush()
}

func Replace(line string) {
	tm.MoveCursorUp(1)
	tm.Print("\033[K")
	Append(line)
}
