package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gdamore/tcell/v2"
	creed "github.com/marekmaskarinec/creed/lib"
)

type Buffer struct {
	// inst is a pointer, since there can be multiple buffers with one instance
	// open at the same time
	inst *creed.Instance
}

func newBuffer() Buffer {
	inst := new(creed.Instance)
	*inst = creed.NewInstance("", os.Stderr)
	return Buffer{inst: inst}
}

func (buf *Buffer) drawRune(srn tcell.Screen, x, y, idx int) {
	r := ' '
	if idx < len(buf.inst.Buf) {
		r = buf.inst.Buf[idx]
	}

	s := styleDefault
	if idx >= buf.inst.Sel.Index &&
		idx < buf.inst.Sel.Index+buf.inst.Sel.Length {
		s = styleSel
	}
	srn.SetContent(x, y, r, nil, s)
}

func (buf *Buffer) drawBorder(srn tcell.Screen, box Box, selected bool) {
	boxStyle := styleDefault
	if selected {
		boxStyle = styleSel
	}
	box.draw(srn, boxStyle)
}

func (buf *Buffer) drawContent(srn tcell.Screen, box Box) {
	min, _ := getLFs(buf.inst.Buf, buf.inst.Sel.Index-1)
	sx := box.x
	sy := box.y + box.h/2

	idx := min - 1
	for y := sy - 1; y > box.y && idx >= 0; y-- {
		mi, mx := getLFs(buf.inst.Buf, idx)
		idx = mi - 1

		x := sx
		for i := mi + 1; i < mx+1; i++ {
			buf.drawRune(srn, x, y, i)
			if buf.inst.Buf[i] == '\t' {
				x++
			}
			x++
		}
	}

	idx = min + 1
	for y := sy; y <= box.y+box.h && idx < len(buf.inst.Buf); y++ {
		mi, mx := getLFs(buf.inst.Buf, idx)
		if mi == idx {
			buf.drawRune(srn, sx, y, idx)
		}
		idx = mx + 1

		x := sx
		for i := mi + 1; i < mx+1 && i <= len(buf.inst.Buf); i++ {
			buf.drawRune(srn, x, y, i)
			if i < len(buf.inst.Buf) && buf.inst.Buf[i] == '\t' {
				x++
			}
			x++
			if x >= box.x + box.w - 1 {
				x = sx + 4
				y++
			}
		}
	}
}

func (buf *Buffer) drawStatusBar(srn tcell.Screen, box Box) {
	str := filepath.Base(buf.inst.Filename)
	if !buf.inst.Saved {
		str += "[*]"
	}

	if buf.inst.Filename == "" {
		str = "<no file>"
	}

	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		srn.SetContent(box.x+3+i, box.y, runes[i], nil, styleSel)
	}
}

func (buf *Buffer) drawLineNumbers(srn tcell.Screen, box Box) int {
	x := box.x + 1
	lno := 0
	lineCount := 0

	for i:=0; i < len(buf.inst.Buf); i++ {
		if buf.inst.Buf[i] == '\n' {
			lineCount++

			if i < buf.inst.Sel.Index {
				lno++
			}
		}
	}

	l := lno - box.h/2
	maxLen := len(strconv.Itoa(lineCount))
	for y:=box.y; y <= box.y + box.h; y++ {
		if l < 0 || l >= lineCount {
			srn.SetContent(x + maxLen - 1, y, '~', nil, styleDefault)
		} else {
			s := strconv.Itoa(l)
			off := maxLen - len(s)

			for i:=0; i < len(s); i++ {
				srn.SetContent(x + i + off, y, rune(s[i]), nil, styleDefault)
			}
		}

		l++
	}

	return maxLen + 1
}

func (buf *Buffer) draw(srn tcell.Screen, box Box, selected bool) {
	if buf.inst.Sel.Length == 0 {
		buf.inst.Sel.Length = 1
	}

	buf.drawBorder(srn, box, selected)
	buf.drawStatusBar(srn, box)

	box.x++; box.w -= 2
	box.y++; box.h -= 2
	
	numberWidth := buf.drawLineNumbers(srn, box) + 1
	box.x += numberWidth
	box.w -= numberWidth
	buf.drawContent(srn, box)
}

func (buf *Buffer) insertRune(r rune) {
	s := string(r)
	switch r {
	case '"', '\\':
		s = "\\" + s
	}
	buf.inst.ExecCommand(fmt.Sprintf("\"%s\" p", s))
}

func (buf *Buffer) deleteRune() {
	// special case
	if buf.inst.Sel.Index == len(buf.inst.Buf) && len(buf.inst.Buf) > 0 {
		buf.inst.Buf = buf.inst.Buf[:len(buf.inst.Buf)-1]
	}

	clip := buf.inst.Clipboard
	buf.inst.Sel.Length = 1
	buf.inst.ExecCommand("1 -- d")

	buf.inst.Clipboard = clip
}

func (buf *Buffer) moveUp() {
	min, max := getLFs(buf.inst.Buf, buf.inst.Sel.Index)
	// handle empty rows
	if max-min <= 1 {
		buf.inst.Sel.Index--
	}

	offset := buf.inst.Sel.Index - min

	min, max = getLFs(buf.inst.Buf, min-1)
	buf.inst.Sel.Index = min + offset
	if buf.inst.Sel.Index > max {
		buf.inst.Sel.Index = max
	}

	// handle empty rows
	if max-min <= 1 {
		buf.inst.Sel.Index++
	}

	if buf.inst.Sel.Index < 0 {
		buf.inst.Sel.Index = 0
	}
}

func (buf *Buffer) moveDown() {
	min, max := getLFs(buf.inst.Buf, buf.inst.Sel.Index)

	// handle empty rows
	if max-min <= 1 {
		buf.inst.Sel.Index++
	}
	offset := buf.inst.Sel.Index - min

	min, max = getLFs(buf.inst.Buf, max+1)
	buf.inst.Sel.Index = min + offset
	if buf.inst.Sel.Index > max {
		buf.inst.Sel.Index = max - 1
	}

	// handle empty rows
	if max-min <= 1 {
		buf.inst.Sel.Index++
	}
}
