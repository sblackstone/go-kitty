package kitty

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type BouncySquare struct {
	X1  int
	Y1  int
	Len int
	Vx  int
	Vy  int
}

func (s *BouncySquare) Update(screen tcell.Screen) {
	width, height := screen.Size()

	s.X1 += s.Vx
	s.Y1 += s.Vy
	if s.X1 < 0 {
		s.X1 = 0
		s.Vx = -s.Vx
	}
	if s.Y1 < 0 {
		s.Y1 = 0
		s.Vy = -s.Vy
	}
	if s.X1+s.Len >= width {
		s.X1 = width - s.Len
		s.Vx = -s.Vx
	}
	if s.Y1+s.Len >= height {
		s.Y1 = height - s.Len
		s.Vy = -s.Vy
	}
}

func (s *BouncySquare) Draw(screen tcell.Screen) {
	for x := s.X1; x < s.X1+s.Len; x++ {
		for y := s.Y1; y < s.Y1+s.Len; y++ {
			screen.SetContent(x, y, tcell.RuneHLine, nil, tcell.StyleDefault.Foreground(color.Red).Background(color.White))
		}
	}
}
