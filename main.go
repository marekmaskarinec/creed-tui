package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

var (
	styleDefault, styleSel tcell.Style
)

func main() {
	styleDefault = tcell.StyleDefault
	styleSel = styleDefault.Reverse(true)

	ed, err := newEd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer ed.srn.Fini()

	for i:=1; i < len(os.Args); i++ {
		ed.buffers = append(ed.buffers, newBuffer())
		ed.bufferi = len(ed.buffers) - 1
		err := ed.buffers[ed.bufferi].inst.ExecCommand(
			fmt.Sprintf("\"%s\" read", os.Args[i]))
		if err != nil {
			fmt.Fprintf(ed, "%v", err)
		}
	}

	for {
		err := ed.handle()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			break
		}

		ed.draw()
	}
}
