package main

import (
	"errors"
	"os"

	"github.com/gdamore/tcell/v2"
	creed "github.com/marekmaskarinec/creed/lib"
)

type Ed struct {
	bufferi     int
	commandMode bool
	buffers     []Buffer
	// the cli is a buffer too to simplify things
	cli       Buffer
	cliHeight int

	srn tcell.Screen
}

func newEd() (Ed, error) {
	ed := Ed{}
	var err error
	ed.srn, err = tcell.NewScreen()
	if err != nil {
		return ed, err
	}

	if err = ed.srn.Init(); err != nil {
		return ed, err
	}

	inst := new(creed.Instance)
	*inst = creed.NewInstance("hello world", os.Stderr)
	ed.buffers = []Buffer{
		Buffer{
			inst: inst,
		},
	}
	inst.Sel.Length = 1

	inst = new(creed.Instance)
	*inst = creed.NewInstance("", os.Stderr)
	ed.cli = Buffer{inst: inst}
	ed.cliHeight = 2

	return ed, err
}

func (ed *Ed) handleKey(ev *tcell.EventKey) error {
	buf := &ed.cli
	if len(ed.buffers) != 0 && !ed.commandMode {
		buf = &ed.buffers[ed.bufferi]
	}

	switch ev.Key() {
	case tcell.KeyRune:
		buf.insertRune(ev.Rune())
	case tcell.KeyBackspace2:
		buf.deleteRune()
	case tcell.KeyEscape:
		ed.commandMode = !ed.commandMode
	case tcell.KeyEnter:
		if ed.commandMode {
			var execBuf *Buffer
			if len(ed.buffers) != 0 {
				execBuf = &ed.buffers[ed.bufferi]
			}
			if err := execBuf.inst.ExecCommand(string(ed.cli.inst.Buf)); err != nil {
				//return err
			}

			ed.cli.inst.Buf = []rune{}
			ed.cli.inst.Sel = creed.Sel{
				Index:  0,
				Length: 1,
			}
		} else {
			buf.insertRune('\n')
		}

	case tcell.KeyLeft:
		buf.inst.ExecCommand("1 --")
	case tcell.KeyRight:
		buf.inst.ExecCommand("1 ++")

	case tcell.KeyUp:
		buf.moveUp()

	case tcell.KeyDown:
		buf.moveDown()

	case tcell.KeyCtrlT:
		ed.buffers = append(ed.buffers, newBuffer())
		ed.bufferi = len(ed.buffers) - 1

	case tcell.KeyCtrlW:
		if len(ed.buffers) != 0 {
			ed.buffers = append(ed.buffers[:ed.bufferi], ed.buffers[ed.bufferi+1:]...)
		}
		ed.bufferi = 0

	case tcell.KeyPgUp:
		if len(ed.buffers) != 0 {
			ed.bufferi--
			if ed.bufferi < 0 {
				ed.bufferi = len(ed.buffers) - 1
			}
		}

	case tcell.KeyPgDn:
		ed.bufferi++
		if ed.bufferi >= len(ed.buffers) {
			ed.bufferi = 0
		}

	// TODO: buf operations

	case tcell.KeyCtrlD:
		return errors.New("bye")
	}

	return nil
}

func (ed *Ed) handle() error {
	ev := ed.srn.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventKey:
		return ed.handleKey(ev)
	}

	return nil
}

func (ed *Ed) draw() {
	ed.srn.Clear()
	w, h := ed.srn.Size()
	if len(ed.buffers) != 0 {
		w /= len(ed.buffers)
	}

	box := Box{
		x: 0, y: 0,
		w: w - 1, h: h - 1 - ed.cliHeight,
	}

	mainBufferX := 0

	for i := 0; i < len(ed.buffers); i++ {
		// make sure the selected buffer is drawn on the top
		if i == ed.bufferi {
			mainBufferX = box.x
			box.x += w
			continue
		}
		ed.buffers[i].draw(ed.srn, box, false)
		box.x += w
	}

	w, h = ed.srn.Size()
	ed.cli.draw(ed.srn, Box{
		x: 0, y: h - ed.cliHeight - 1,
		w: w - 1, h: ed.cliHeight,
	}, ed.commandMode)
	if len(ed.buffers) != 0 {
		box.x = mainBufferX
		ed.buffers[ed.bufferi].draw(ed.srn, box, true)
	}

	ed.srn.Show()
}
