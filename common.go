package main

import "github.com/gdamore/tcell/v2"

type Box struct {
	x, y, w, h int
}

func (box *Box) draw(srn tcell.Screen, style tcell.Style) {
	for x := box.x; x < box.x+box.w; x++ {
		srn.SetContent(x, box.y, '-', nil, style)
		srn.SetContent(x, box.y+box.h, '-', nil, style)
	}

	for y := box.y; y < box.y+box.h; y++ {
		srn.SetContent(box.x, y, '|', nil, style)
		srn.SetContent(box.x+box.w, y, '|', nil, style)
	}

	srn.SetContent(box.x, box.y, '+', nil, style)
	srn.SetContent(box.x+box.w, box.y, '+', nil, style)
	srn.SetContent(box.x, box.y+box.h, '+', nil, style)
	srn.SetContent(box.x+box.w, box.y+box.h, '+', nil, style)
}

// Returns the indices of the LF of the current line
func getLFs(buf []rune, idx int) (int, int) {
	min := -1
	max := len(buf)

	for i := idx; i >= 0; i-- {
		if i >= len(buf) {
			continue
		}

		if buf[i] == '\n' {
			min = i
			break
		}
	}

	for i := idx; i < len(buf); i++ {
		if i < 0 {
			continue
		}

		if buf[i] == '\n' {
			max = i
			break
		}
	}

	return min, max
}
