package kitty

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Butterfly struct {
	Color tcell.Color

	active       bool
	respawnWait  int
	initDelaySet bool

	x         float64
	baseY     float64
	vx        float64
	waveAmp   float64
	wavePhase float64
	flapPhase float64
	dir       int
	flutterTicks int
	burstTicks   int
	turnBias     float64
}

var butterflyRandSeeded bool

func (b *Butterfly) Update(screen tcell.Screen) {
	width, height := screen.Size()
	if width <= 0 || height <= 0 {
		return
	}
	if !butterflyRandSeeded {
		rand.Seed(time.Now().UnixNano())
		butterflyRandSeeded = true
	}
	if !b.active && !b.initDelaySet {
		b.respawnWait = rand.Intn(80)
		b.initDelaySet = true
	}
	if b.respawnWait > 0 {
		b.respawnWait--
		return
	}
	if !b.active {
		b.initButterfly(width, height)
		return
	}

	// flutter and dart behavior for prey-like motion
	b.wavePhase += 0.18 + rand.Float64()*0.08
	b.flapPhase += 0.7 + rand.Float64()*0.25

	if b.flutterTicks > 0 {
		b.flutterTicks--
		b.vx = randRange(0.3, 0.8)
		b.waveAmp = clampFloat(b.waveAmp+randRange(-0.15, 0.15), 0.5, 4.0)
	} else if b.burstTicks > 0 {
		b.burstTicks--
		b.vx = randRange(1.6, 2.6)
		if rand.Float64() < 0.15 {
			b.turnBias *= -1
		}
	} else {
		if rand.Float64() < 0.02 {
			b.flutterTicks = 10 + rand.Intn(18)
		}
		if rand.Float64() < 0.02 {
			b.burstTicks = 6 + rand.Intn(12)
		}
		b.vx = clampFloat(b.vx+randRange(-0.08, 0.08), 0.5, 1.6)
	}

	if rand.Float64() < 0.01 {
		b.turnBias = randRange(-1.0, 1.0)
	}

	b.x += b.vx * float64(b.dir)

	if b.dir > 0 && b.x > float64(width+2) {
		b.active = false
		b.respawnWait = 40 + rand.Intn(80)
		return
	}
	if b.dir < 0 && b.x < -2 {
		b.active = false
		b.respawnWait = 40 + rand.Intn(80)
		return
	}
}

func (b *Butterfly) Draw(screen tcell.Screen) {
	if !b.active {
		return
	}
	width, height := screen.Size()
	cx := int(math.Round(b.x))
	wobble := math.Sin(b.wavePhase*1.7) * 0.8
	cy := int(math.Round(b.baseY + math.Sin(b.wavePhase)*b.waveAmp + wobble))

	if cx < 0 || cy < 0 || cx >= width || cy >= height {
		return
	}

	fg := b.Color
	if fg == tcell.ColorDefault || fg == 0 {
		fg = randomButterflyColor()
		b.Color = fg
	}

	up := math.Sin(b.flapPhase+b.turnBias) > 0
	leftWing := '/'
	rightWing := '\\'
	if !up {
		leftWing = '\\'
		rightWing = '/'
	}

	b.drawCell(screen, cx-1, cy-1, leftWing, fg, width, height)
	b.drawCell(screen, cx+1, cy-1, rightWing, fg, width, height)
	b.drawCell(screen, cx, cy, tcell.RuneVLine, fg, width, height)
	b.drawCell(screen, cx-1, cy+1, leftWing, fg, width, height)
	b.drawCell(screen, cx+1, cy+1, rightWing, fg, width, height)
}

func (b *Butterfly) drawCell(screen tcell.Screen, x, y int, r rune, fg tcell.Color, width, height int) {
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	screen.SetContent(x, y, r, nil, tcell.StyleDefault.Foreground(fg))
}

func (b *Butterfly) initButterfly(width, height int) {
	b.active = true
	b.wavePhase = randRange(0, math.Pi*2)
	b.flapPhase = randRange(0, math.Pi*2)
	b.waveAmp = randRange(0.5, 2.5)
	b.vx = randRange(0.6, 1.4)
	b.dir = 1
	if rand.Intn(2) == 0 {
		b.dir = -1
	}
	minY := 1
	maxY := height - 2
	if maxY < minY {
		maxY = minY
	}
	b.baseY = float64(minY + rand.Intn(maxY-minY+1))
	if b.dir > 0 {
		b.x = -2
	} else {
		b.x = float64(width + 2)
	}
	b.flutterTicks = 0
	b.burstTicks = 0
	b.turnBias = randRange(-1.0, 1.0)
	b.Color = randomButterflyColor()
}

func randomButterflyColor() tcell.Color {
	colors := []tcell.Color{
		color.Fuchsia,
		color.Purple,
		color.Orange,
		color.Yellow,
		color.Aqua,
		color.Lime,
		color.White,
	}
	return colors[rand.Intn(len(colors))]
}

