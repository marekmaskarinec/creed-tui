package main

import (
	"errors"
	"fmt"

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
	writer string

	srn tcell.Screen
}

func newEd() (*Ed, error) {
	ed := Ed{}
	var err error
	ed.srn, err = tcell.NewScreen()
	if err != nil {
		return &ed, err
	}

	if err = ed.srn.Init(); err != nil {
		return &ed, err
	}

	inst := new(creed.Instance)
	*inst = creed.NewInstance("", &ed)
	ed.cli = Buffer{inst: inst}
	ed.cliHeight = 2

	return &ed, err
}

// Implement io.Writer
func (ed *Ed) Write(p []byte) (n int, err error) {
	ed.writer = ""
	for i:=0; i < len(p); i++ {
		ed.writer += string(rune(p[i]))
	}
	return len(p), nil	
}

func (ed *Ed) handleKey(ev *tcell.EventKey) error {
	buf := &ed.cli
	if len(ed.buffers) != 0 && !ed.commandMode {
		buf = &ed.buffers[ed.bufferi]
	}

	switch ev.Key() {
	case tcell.KeyRune:
		buf.insertRune(ev.Rune())
	case tcell.KeyTab:
		buf.insertRune('\t')
	case tcell.KeyBackspace2:
		buf.deleteRune()
	case tcell.KeyEscape:
		ed.commandMode = !ed.commandMode
	case tcell.KeyEnter:
		if ed.commandMode {
			if len(ed.buffers) != 0 {
				execBuf := &ed.buffers[ed.bufferi]
				if err := execBuf.inst.ExecCommand(string(ed.cli.inst.Buf)); err != nil {
					fmt.Fprintf(ed, "%v", err)
				}

				ed.cli.inst.Buf = []rune{}
				ed.cli.inst.Sel = creed.Sel{
					Index:  0,
					Length: 1,
				}
			} else {
				fmt.Fprintf(ed, "no buffer")
			}
		} else {
			buf.insertRune('\n')
		}

	case tcell.KeyLeft:
		buf.inst.ExecCommand("1 --")
	case tcell.KeyRight:
		buf.inst.ExecCommand("1 ++")

	case tcell.KeyUp:
		if len(ed.buffers) != 0 {
			ed.buffers[ed.bufferi].moveUp()
		}

	case tcell.KeyDown:
		if len(ed.buffers) != 0 {
			ed.buffers[ed.bufferi].moveDown()
		}

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

func (ed *Ed) drawWriter(box Box) {
	if len(ed.writer) == 0 { return }

	x := box.x + 1
	y := box.y + 1

	for i:=0; i < len(ed.writer); i++ {
		if ed.writer[i] == '\n' || x >= box.x + box.w {
			y++
			x = box.x + 1
		}

		ed.srn.SetContent(x, y, rune(ed.writer[i]), nil, styleDefault)
		x++
	}
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
  cliBox := Box{
		x: 0, y: h - ed.cliHeight - 1,
		w: w/2 - 1, h: ed.cliHeight,
	}
	ed.cli.draw(ed.srn, cliBox, ed.commandMode)
	if ed.commandMode {
 	 cliBox.draw(ed.srn, styleSel)
	}

	ed.drawWriter(Box{
		x: w/2, y: h - ed.cliHeight - 1,
		w: w/2, h: ed.cliHeight + 1})

	if len(ed.buffers) != 0 {
		box.x = mainBufferX
		ed.buffers[ed.bufferi].draw(ed.srn, box, true)
	}

	ed.srn.Show()
}
