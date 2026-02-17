package kitty

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type SwayString struct {
	MinLen int
	MaxLen int
	Color  tcell.Color

	respawnWait int
	spawnDelay  int
	initDelaySet bool
	initialDelayMax int
	step        int
	lifeSteps   int
	length      int
	phase       float64
	breezePhase float64
	breezeTicks int
	breezeDir   float64

	anchorX int
	anchorY int
	dirX    float64
	dirY    float64
	perpX   float64
	perpY   float64
	swingAmp float64
}

var swayRandSeeded bool

func (s *SwayString) Update(screen tcell.Screen) {
	if s.respawnWait > 0 {
		s.respawnWait--
		return
	}
	if s.lifeSteps == 0 {
		if !s.initDelaySet {
			ensureSwaySeeded()
			maxDelay := s.initialDelayMax
			if maxDelay <= 0 {
				maxDelay = 40
			}
			s.spawnDelay = rand.Intn(maxDelay + 1)
			s.initDelaySet = true
		}
		if s.spawnDelay > 0 {
			s.spawnDelay--
			return
		}
		s.initString(screen)
		return
	}

	s.step++
	s.phase += 0.25
	s.breezePhase += 0.12
	if s.breezeTicks > 0 {
		s.breezeTicks--
	} else if rand.Float64() < 0.015 {
		s.breezeTicks = 40 + rand.Intn(80)
		s.breezeDir = randRange(-1.0, 1.0)
	}
	if s.step >= s.lifeSteps {
		s.lifeSteps = 0
		s.respawnWait = 20 + rand.Intn(60)
	}
}

func (s *SwayString) Draw(screen tcell.Screen) {
	if s.lifeSteps == 0 {
		return
	}
	width, height := screen.Size()
	fg := s.Color
	if fg == tcell.ColorDefault || fg == 0 {
		fg = randomStringColor()
		s.Color = fg
	}

	u := float64(s.step) / float64(max(1, s.lifeSteps-1))
	lengthFactor := 1 - math.Abs(1-2*u)
	curLen := int(math.Round(float64(s.length) * lengthFactor))
	if curLen < 1 {
		return
	}

	var breezeOffset float64
	if s.breezeTicks > 0 {
		breezeOffset = math.Sin(s.breezePhase) * (0.8 + 0.4*math.Sin(s.breezePhase*0.5)) * s.swingAmp * s.breezeDir
	}
	for i := 0; i < curLen; i++ {
		flex := float64(i) / float64(max(1, curLen-1))
		localSwing := math.Sin(s.phase+u*math.Pi*2+float64(i)*0.45) * s.swingAmp * (0.2 + 0.8*flex)
		bend := math.Sin(s.phase*0.7+float64(i)*0.25) * (0.15 + 0.85*flex)
		wind := breezeOffset * (0.2 + 0.8*flex)
		fx := float64(i)*s.dirX + (localSwing+bend+wind)*s.perpX
		fy := float64(i)*s.dirY + (localSwing+bend+wind)*s.perpY
		x := s.anchorX + int(math.Round(fx))
		y := s.anchorY + int(math.Round(fy))
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				px := x + dx
				py := y + dy
				if px < 0 || py < 0 || px >= width || py >= height {
					continue
				}
				screen.SetContent(px, py, tcell.RuneBlock, nil, tcell.StyleDefault.Foreground(fg))
			}
		}
	}
}

func (s *SwayString) initString(screen tcell.Screen) {
	width, height := screen.Size()
	if width <= 0 || height <= 0 {
		return
	}
	ensureSwaySeeded()

	minLen := s.MinLen
	maxLen := s.MaxLen
	if minLen <= 0 {
		minLen = 4
	}
	if maxLen < minLen {
		maxLen = minLen + 6
	}

	s.length = minLen + rand.Intn(maxLen-minLen+1)
	s.lifeSteps = 40 + rand.Intn(80)
	s.step = 0
	s.phase = randRange(0, math.Pi*2)
	s.breezePhase = randRange(0, math.Pi*2)
	s.breezeTicks = 0
	s.breezeDir = randRange(-1.0, 1.0)
	s.swingAmp = randRange(1.5, 6.5)

	if s.Color == tcell.ColorDefault || s.Color == 0 {
		s.Color = randomStringColor()
	}

	edge := rand.Intn(4)
	if edge == 0 { // left
		s.anchorX = 0
		s.anchorY = rand.Intn(height)
		s.dirX = 1
		s.dirY = 0
	} else if edge == 1 { // right
		s.anchorX = width - 1
		s.anchorY = rand.Intn(height)
		s.dirX = -1
		s.dirY = 0
	} else if edge == 2 { // top
		s.anchorX = rand.Intn(width)
		s.anchorY = 0
		s.dirX = 0
		s.dirY = 1
	} else { // bottom
		s.anchorX = rand.Intn(width)
		s.anchorY = height - 1
		s.dirX = 0
		s.dirY = -1
	}

	s.perpX = -s.dirY
	s.perpY = s.dirX
}

func NewSwayString(cfg SwayStringConfig) *SwayString {
	if cfg.MinLen <= 0 {
		cfg.MinLen = 18
	}
	if cfg.MaxLen < cfg.MinLen {
		cfg.MaxLen = cfg.MinLen + 6
	}
	if cfg.InitialDelayMax <= 0 {
		cfg.InitialDelayMax = 40
	}
	return &SwayString{
		MinLen:          cfg.MinLen,
		MaxLen:          cfg.MaxLen,
		Color:           cfg.Color,
		initialDelayMax: cfg.InitialDelayMax,
	}
}

func randomStringColor() tcell.Color {
	colors := []tcell.Color{
		color.Red,
		color.Orange,
		color.Yellow,
		color.Green,
		color.Teal,
		color.Aqua,
		color.Blue,
		color.Navy,
		color.Purple,
		color.Fuchsia,
		color.Maroon,
		color.Lime,
		color.White,
	}
	return colors[rand.Intn(len(colors))]
}

func ensureSwaySeeded() {
	if !swayRandSeeded {
		rand.Seed(time.Now().UnixNano())
		swayRandSeeded = true
	}
}