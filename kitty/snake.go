package kitty

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Snake struct {
	MaxLen int
	Color  tcell.Color

	body        []Point
	curLen      int
	step        int
	phase       float64
	initialized bool
	progress    float64
	speed       float64
	speedTarget float64
	zoomTicks   int
	zoomOffTicks int
	zoomOffTargetSet bool
	zoomOffX float64
	zoomOffY float64
	respawnWait int

	posX            float64
	posY            float64
	heading         float64
	turnTarget      float64
	turnSpeed       float64
	amplitude       float64
	amplitudeTarget float64
}

var snakeRandSeeded bool

func (s *Snake) Draw(screen tcell.Screen) {
	width, height := screen.Size()
	fg := s.Color
	if fg == tcell.ColorDefault || fg == 0 {
		fg = color.Green
	}
	for _, p := range s.body {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				x := p.X + dx
				y := p.Y + dy
				if x < 0 || y < 0 || x >= width || y >= height {
					continue
				}
				screen.SetContent(x, y, tcell.RuneBlock, nil, tcell.StyleDefault.Foreground(fg))
			}
		}
	}
}

func (s *Snake) Update(screen tcell.Screen) {
	if s.MaxLen <= 0 {
		s.MaxLen = 10
	}
	if s.respawnWait > 0 {
		s.respawnWait--
		return
	}
	if !s.initialized {
		s.initSnake(screen)
	}

	width, height := screen.Size()
	s.updateSpeed()
	s.progress += s.speed
	for s.progress >= 1.0 {
		s.progress -= 1.0
		s.updateSteering(width, height)
		s.phase += 0.4
		s.step++

		s.posX += math.Cos(s.heading)
		s.posY += math.Sin(s.heading)

		head := s.nextHead()
		s.body = append(s.body, head)
		if s.curLen < s.MaxLen {
			s.curLen++
		}
		if len(s.body) > s.curLen {
			s.body = s.body[len(s.body)-s.curLen:]
		}
	}

	if s.shouldReset(width, height) {
		s.respawnWait = 20 + rand.Intn(40)
		s.initialized = false
		s.body = s.body[:0]
	}
}

func (s *Snake) initSnake(screen tcell.Screen) {
	width, height := screen.Size()
	if width <= 0 || height <= 0 {
		return
	}
	if !snakeRandSeeded {
		rand.Seed(time.Now().UnixNano())
		snakeRandSeeded = true
	}

	side := rand.Intn(4)
	maxAmp := 6
	if width < maxAmp*2 {
		maxAmp = max(1, width/4)
	}
	if height < maxAmp*2 {
		maxAmp = max(1, height/4)
	}

	s.amplitude = float64(max(1, maxAmp))
	s.amplitudeTarget = s.amplitude
	s.phase = 0
	s.step = 0
	s.curLen = 1
	s.progress = 0
	s.speed = 1.0
	s.speedTarget = 1.0
	s.zoomTicks = 0
	s.zoomOffTicks = 0
	s.zoomOffTargetSet = false
	s.respawnWait = 0
	s.body = s.body[:0]
	s.initialized = true
	if s.Color == tcell.ColorDefault || s.Color == 0 {
		s.Color = randomSnakeColor()
	}

	if side == 0 { // left -> right
		s.posX = -1
		s.posY = randRange(0, float64(max(1, height-1)))
		s.heading = randRange(-0.6, 0.6)
	} else if side == 1 { // right -> left
		s.posX = float64(width)
		s.posY = randRange(0, float64(max(1, height-1)))
		s.heading = randRange(math.Pi-0.6, math.Pi+0.6)
	} else if side == 2 { // top -> bottom
		s.posX = randRange(0, float64(max(1, width-1)))
		s.posY = -1
		s.heading = randRange(math.Pi/2-0.6, math.Pi/2+0.6)
	} else { // bottom -> top
		s.posX = randRange(0, float64(max(1, width-1)))
		s.posY = float64(height)
		s.heading = randRange(-math.Pi/2-0.6, -math.Pi/2+0.6)
	}

	s.turnTarget = s.heading
	s.turnSpeed = randRange(0.03, 0.12)
}

func (s *Snake) updateSpeed() {
	if s.zoomOffTicks > 0 {
		s.zoomOffTicks--
		s.speedTarget = randRange(6.0, 9.0)
	} else if s.zoomTicks > 0 {
		s.zoomTicks--
	} else {
		// small random drift while in normal mode
		s.speedTarget += (rand.Float64() - 0.5) * 0.05
		s.speedTarget = clampFloat(s.speedTarget, 0.4, 1.6)
		// occasional zoom burst
		if rand.Float64() < 0.02 {
			s.speedTarget = randRange(2.5, 5.0)
			s.zoomTicks = 10 + rand.Intn(20)
		}
		// rare zoom-off to exit
		if rand.Float64() < 0.006 {
			s.zoomOffTicks = 20 + rand.Intn(30)
			s.zoomOffTargetSet = false
		}
	}
	// smooth change toward target
	s.speed += (s.speedTarget - s.speed) * 0.1
}

func (s *Snake) updateSteering(width, height int) {
	if s.zoomOffTicks > 0 {
		if !s.zoomOffTargetSet {
			s.zoomOffX, s.zoomOffY = randomEdgePoint(width, height)
			s.zoomOffTargetSet = true
		}
		angle := math.Atan2(s.zoomOffY-s.posY, s.zoomOffX-s.posX)
		s.turnTarget = angle
		s.turnSpeed = 0.25
		s.amplitudeTarget = 1.0
		// rotate heading toward target
		delta := normalizeAngle(s.turnTarget - s.heading)
		if delta > s.turnSpeed {
			delta = s.turnSpeed
		} else if delta < -s.turnSpeed {
			delta = -s.turnSpeed
		}
		s.heading = normalizeAngle(s.heading + delta)
		return
	}
	// drift target heading a bit for chaos
	s.turnTarget += (rand.Float64()-0.5)*0.08 + math.Sin(s.phase)*0.01
	// occasional bigger turn
	if rand.Float64() < 0.03 {
		s.turnTarget = s.heading + randRange(-1.2, 1.2)
	}
	// sometimes aim toward a random edge to allow any exit
	if rand.Float64() < 0.015 {
		x, y := randomEdgePoint(width, height)
		angle := math.Atan2(y-s.posY, x-s.posX)
		s.turnTarget = angle
	}

	// smooth amplitude changes
	if rand.Float64() < 0.02 {
		s.amplitudeTarget = randRange(2.0, 10.0)
	}
	s.amplitude += (s.amplitudeTarget - s.amplitude) * 0.05

	// rotate heading toward target
	delta := normalizeAngle(s.turnTarget - s.heading)
	if delta > s.turnSpeed {
		delta = s.turnSpeed
	} else if delta < -s.turnSpeed {
		delta = -s.turnSpeed
	}
	s.heading = normalizeAngle(s.heading + delta)
}

func (s *Snake) nextHead() Point {
	perpX := -math.Sin(s.heading)
	perpY := math.Cos(s.heading)
	offset := math.Sin(s.phase) * s.amplitude
	x := s.posX + perpX*offset
	y := s.posY + perpY*offset
	return Point{X: int(math.Round(x)), Y: int(math.Round(y))}
}

func (s *Snake) shouldReset(width, height int) bool {
	if s.step <= s.MaxLen {
		return false
	}
	for _, p := range s.body {
		if p.X >= 0 && p.Y >= 0 && p.X < width && p.Y < height {
			return false
		}
	}
	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampFloat(v, minV, maxV float64) float64 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func randRange(minV, maxV float64) float64 {
	return minV + rand.Float64()*(maxV-minV)
}

func randomSnakeColor() tcell.Color {
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
	}
	return colors[rand.Intn(len(colors))]
}

func normalizeAngle(a float64) float64 {
	for a > math.Pi {
		a -= math.Pi * 2
	}
	for a < -math.Pi {
		a += math.Pi * 2
	}
	return a
}

func randomEdgePoint(width, height int) (float64, float64) {
	if width <= 0 || height <= 0 {
		return 0, 0
	}
	edge := rand.Intn(4)
	if edge == 0 { // left
		return -1, randRange(0, float64(height-1))
	}
	if edge == 1 { // right
		return float64(width), randRange(0, float64(height-1))
	}
	if edge == 2 { // top
		return randRange(0, float64(width-1)), -1
	}
	// bottom
	return randRange(0, float64(width-1)), float64(height)
}
