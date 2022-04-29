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

	for {
		err := ed.handle()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			break
		}

		ed.draw()
	}
}
