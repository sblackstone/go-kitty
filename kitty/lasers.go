package kitty

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type LaserPointer struct {
	Color tcell.Color

	active        bool
	respawnWait   int
	initDelaySet  bool
	initialDelayMax int

	x       float64
	y       float64
	baseSpeed float64
	speed     float64
	targetX   float64
	targetY   float64
	pauseTicks int
	dashTicks  int
}

var laserRandSeeded bool

func (l *LaserPointer) Update(screen tcell.Screen) {
	width, height := screen.Size()
	if width <= 0 || height <= 0 {
		return
	}
	ensureLaserSeeded()

	if !l.active && !l.initDelaySet {
		maxDelay := l.initialDelayMax
		if maxDelay <= 0 {
			maxDelay = 80
		}
		l.respawnWait = rand.Intn(maxDelay + 1)
		l.initDelaySet = true
	}
	if l.respawnWait > 0 {
		l.respawnWait--
		return
	}
	if !l.active {
		l.initLaser(width, height)
		return
	}

	if l.pauseTicks > 0 {
		l.pauseTicks--
		return
	}
	if l.dashTicks > 0 {
		l.dashTicks--
		l.speed = randRange(2.5, 4.0)
	} else {
		l.speed += (l.baseSpeed - l.speed) * 0.12
		if rand.Float64() < 0.01 {
			l.pauseTicks = 4 + rand.Intn(8)
		}
		if rand.Float64() < 0.05 {
			l.dashTicks = 6 + rand.Intn(12)
		}
	}

	dx := l.targetX - l.x
	dy := l.targetY - l.y
	dist := math.Hypot(dx, dy)
	if dist < 1.2 || rand.Float64() < 0.04 {
		l.targetX = randRange(1, float64(width-2))
		l.targetY = randRange(1, float64(height-2))
		return
	}
	step := l.speed / math.Max(dist, 0.001)
	l.x += dx * step
	l.y += dy * step
}

func (l *LaserPointer) Draw(screen tcell.Screen) {
	if !l.active {
		return
	}
	width, height := screen.Size()
	cx := int(math.Round(l.x))
	cy := int(math.Round(l.y))
	if cx < 0 || cy < 0 || cx >= width || cy >= height {
		return
	}

	fg := l.Color
	if fg == tcell.ColorDefault || fg == 0 {
		fg = color.Red
		l.Color = fg
	}
	glow := color.DarkRed

	screen.SetContent(cx, cy, tcell.RuneBlock, nil, tcell.StyleDefault.Foreground(fg))
	for _, p := range []Point{{X: cx - 1, Y: cy}, {X: cx + 1, Y: cy}, {X: cx, Y: cy - 1}, {X: cx, Y: cy + 1}} {
		if p.X < 0 || p.Y < 0 || p.X >= width || p.Y >= height {
			continue
		}
		screen.SetContent(p.X, p.Y, tcell.RuneBlock, nil, tcell.StyleDefault.Foreground(glow))
	}
	for _, p := range []Point{{X: cx - 1, Y: cy - 1}, {X: cx + 1, Y: cy - 1}, {X: cx - 1, Y: cy + 1}, {X: cx + 1, Y: cy + 1}} {
		if p.X < 0 || p.Y < 0 || p.X >= width || p.Y >= height {
			continue
		}
		screen.SetContent(p.X, p.Y, tcell.RuneBullet, nil, tcell.StyleDefault.Foreground(glow))
	}
}

func (l *LaserPointer) Position(width, height int) (Point, bool) {
	if !l.active {
		return Point{}, false
	}
	cx := int(math.Round(l.x))
	cy := int(math.Round(l.y))
	if cx < 0 || cy < 0 || cx >= width || cy >= height {
		return Point{}, false
	}
	return Point{X: cx, Y: cy}, true
}

func (l *LaserPointer) initLaser(width, height int) {
	l.active = true
	l.baseSpeed = randRange(1.0, 2.2)
	l.speed = l.baseSpeed
	l.x = randRange(1, float64(width-2))
	l.y = randRange(1, float64(height-2))
	l.targetX = randRange(1, float64(width-2))
	l.targetY = randRange(1, float64(height-2))
	l.pauseTicks = 0
	l.dashTicks = 0
	if l.Color == tcell.ColorDefault || l.Color == 0 {
		l.Color = color.Red
	}
}

func NewLaserPointer(cfg LaserConfig) *LaserPointer {
	if cfg.InitialDelayMax <= 0 {
		cfg.InitialDelayMax = 80
	}
	return &LaserPointer{
		Color:           cfg.Color,
		initialDelayMax: cfg.InitialDelayMax,
	}
}

func ensureLaserSeeded() {
	if !laserRandSeeded {
		rand.Seed(time.Now().UnixNano())
		laserRandSeeded = true
	}
}
