module github.com/marekmaskarinec/creed-tui

go 1.18

replace github.com/marekmaskarinec/creed/lib => ../creed/lib

require (
	github.com/gdamore/tcell/v2 v2.5.1
	github.com/marekmaskarinec/creed/lib v0.0.0-00010101000000-000000000000
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.0.0-20220318055525-2edf467146b5 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/text v0.3.7 // indirect
)
