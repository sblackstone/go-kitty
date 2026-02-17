package kitty

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type BouncyBall struct {
	active      bool
	respawnWait int

	x      float64
	y      float64
	vx     float64
	vy     float64
	gravity float64
	radius  int
	dir     int
}

var bouncyBallSeeded bool

func (s *BouncyBall) Update(screen tcell.Screen) {
	if s.radius == 0 {
		s.radius = 3
	}
	if s.gravity == 0 {
		s.gravity = 0.35
	}
	if s.respawnWait > 0 {
		s.respawnWait--
		return
	}
	if !s.active {
		s.initBall(screen)
		return
	}

	width, height := screen.Size()
	ground := float64(height - 1)

	s.vy += s.gravity
	s.x += s.vx
	s.y += s.vy

	if s.y+float64(s.radius) >= ground {
		s.y = ground - float64(s.radius)
		s.vy = -s.vy * 0.7
		if math.Abs(s.vy) < 0.6 {
			s.vy = -randRange(2.5, 4.0)
		}
	}

	if s.dir > 0 && s.x-float64(s.radius) > float64(width) {
		s.active = false
		s.respawnWait = 60 + rand.Intn(140)
		return
	}
	if s.dir < 0 && s.x+float64(s.radius) < 0 {
		s.active = false
		s.respawnWait = 60 + rand.Intn(140)
		return
	}
}

func (s *BouncyBall) Draw(screen tcell.Screen) {
	if !s.active {
		return
	}
	width, height := screen.Size()
	centerX := int(math.Round(s.x))
	centerY := int(math.Round(s.y))

	fg := color.White
	r := float64(s.radius)
	for dy := -s.radius; dy <= s.radius; dy++ {
		for dx := -s.radius; dx <= s.radius; dx++ {
			fx := float64(dx)
			fy := float64(dy)
			if fx*fx+fy*fy > r*r {
				continue
			}
			x := centerX + dx
			y := centerY + dy
			if x < 0 || y < 0 || x >= width || y >= height {
				continue
			}
			screen.SetContent(x, y, tcell.RuneBlock, nil, tcell.StyleDefault.Foreground(fg))
		}
	}
}

func (s *BouncyBall) initBall(screen tcell.Screen) {
	width, height := screen.Size()
	if width <= 0 || height <= 0 {
		return
	}
	if !bouncyBallSeeded {
		rand.Seed(time.Now().UnixNano())
		bouncyBallSeeded = true
	}

	s.active = true
	s.dir = 1
	if rand.Intn(2) == 0 {
		s.dir = -1
	}
	// cross the screen in a few bounces
	s.vx = randRange(1.8, 3.2) * float64(s.dir)
	s.vy = -randRange(3.0, 5.0)
	ground := float64(height - 1)
	s.y = ground - float64(s.radius)
	if s.dir > 0 {
		s.x = -float64(s.radius)
	} else {
		s.x = float64(width + s.radius)
	}
}
