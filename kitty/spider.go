package kitty

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Spider struct {
	Color tcell.Color

	active          bool
	respawnWait     int
	initDelaySet    bool
	initialDelayMax int

	x            float64
	y            float64
	targetX      float64
	targetY      float64
	speed        float64
	pauseTicks   int
	legPhase     float64
	explosionTicks int
	explosionX     int
	explosionY     int
	
	// Web building
	webSegments    []Point
	webSpokes      []Point
	dropSilk       []Point
	webState       string // "dropping", "building", "done", "hunting", "eating", "climbing"
	webBuildStep   int
	dropStartY     float64
	dropTargetY    float64
	centerX        float64
	centerY        float64
	
	// Hunting
	preyX          float64
	preyY          float64
	eatingTicks    int
	webIncomplete  bool
}

var spiderRandSeeded bool

func (s *Spider) Update(screen tcell.Screen) {
	width, height := screen.Size()
	if width <= 0 || height <= 0 {
		return
	}
	if s.explosionTicks > 0 {
		s.explosionTicks--
		return
	}
	ensureSpiderSeeded()

	if !s.active && !s.initDelaySet {
		maxDelay := s.initialDelayMax
		if maxDelay <= 0 {
			maxDelay = 60
		}
		s.respawnWait = rand.Intn(maxDelay + 1)
		s.initDelaySet = true
	}
	if s.respawnWait > 0 {
		s.respawnWait--
		return
	}
	if !s.active {
		s.initSpider(width, height)
		return
	}

	s.legPhase += 0.5

	// State machine for spider web building
	switch s.webState {
	case "dropping":
		// Drop down from top on silk thread
		s.y += 0.8
		// Add silk thread
		webPoint := Point{X: int(math.Round(s.x)), Y: int(math.Round(s.y))}
		if len(s.dropSilk) == 0 || s.dropSilk[len(s.dropSilk)-1] != webPoint {
			s.dropSilk = append(s.dropSilk, webPoint)
		}
		if s.y >= s.dropTargetY {
			s.webState = "building"
			s.centerX = s.x
			s.centerY = s.y
			s.webBuildStep = 0
		}
		return

	case "building":
		// Build classic spider web
		s.buildClassicWeb(width, height)
		return

	case "done":
		// Web is complete, spider rests at center
		if s.pauseTicks > 0 {
			s.pauseTicks--
		} else {
			// After resting, climb back up
			s.webState = "climbing"
			s.x = s.centerX
			s.y = s.centerY
		}
		return
	
	case "hunting":
		// Move towards prey
		dx := s.preyX - s.x
		dy := s.preyY - s.y
		dist := math.Hypot(dx, dy)
		if dist < 1.0 {
			// Reached prey, start eating
			s.webState = "eating"
			s.eatingTicks = 20
		} else {
			// Move quickly towards prey
			speed := 1.5
			s.x += (dx / dist) * speed
			s.y += (dy / dist) * speed
		}
		return
	
	case "eating":
		// Eating animation
		if s.eatingTicks > 0 {
			s.eatingTicks--
		} else {
			// Return to center or resume building
			if s.webIncomplete {
				// Resume building web
				s.webState = "building"
				s.webIncomplete = false
			} else {
				// Return to center
				dx := s.centerX - s.x
				dy := s.centerY - s.y
				dist := math.Hypot(dx, dy)
				if dist < 1.0 {
					s.x = s.centerX
					s.y = s.centerY
					s.webState = "done"
					s.pauseTicks = 100 + rand.Intn(100)
				} else {
					speed := 1.0
					s.x += (dx / dist) * speed
					s.y += (dy / dist) * speed
				}
			}
		}
		return
	
	case "climbing":
		// Climb back up to top
		s.y -= 0.8
		if s.y <= 0 {
			// Reached top, despawn and clear web
			s.active = false
			s.respawnWait = 200 + rand.Intn(300)
			s.webSegments = []Point{}
			s.webSpokes = []Point{}
			s.dropSilk = []Point{}
		}
		return
	}
}

func (s *Spider) buildClassicWeb(width, height int) {
	// Build a classic spider web with radial spokes and concentric circles
	const numSpokes = 8
	const numRings = 4
	const maxRadius = 8.0
	
	step := s.webBuildStep
	
	// Phase 1: Build spokes progressively (first 100 steps)
	if step < numSpokes*25 {
		spokeIndex := step / 25
		spokeProgress := float64(step % 25) / 25.0
		
		if spokeIndex < numSpokes {
			angle := (float64(spokeIndex) / float64(numSpokes)) * 2 * math.Pi
			radius := maxRadius * spokeProgress
			
			// Move spider outward along the spoke
			s.x = s.centerX + math.Cos(angle)*radius
			s.y = s.centerY + math.Sin(angle)*radius
			
			// Add spoke segment
			px := int(math.Round(s.x))
			py := int(math.Round(s.y))
			if px >= 0 && py >= 0 && px < width && py < height {
				s.webSpokes = append(s.webSpokes, Point{X: px, Y: py})
			}
		}
		s.webBuildStep++
		return
	}
	
	// Phase 2: Build concentric circles
	circleStep := step - numSpokes*25
	if circleStep < numRings*50 {
		ringIndex := circleStep / 50
		if ringIndex < numRings {
			radius := maxRadius * float64(ringIndex+1) / float64(numRings)
			angle := (float64(circleStep%50) / 50.0) * 2 * math.Pi
			
			// Move spider along the ring being built
			s.x = s.centerX + math.Cos(angle)*radius
			s.y = s.centerY + math.Sin(angle)*radius
			
			// Add web segment
			px := int(math.Round(s.x))
			py := int(math.Round(s.y))
			if px >= 0 && py >= 0 && px < width && py < height {
				s.webSegments = append(s.webSegments, Point{X: px, Y: py})
			}
		}
		s.webBuildStep++
	} else {
		// Web complete, move spider to center
		s.x = s.centerX
		s.y = s.centerY
		s.webState = "done"
		s.pauseTicks = 300 + rand.Intn(200)
	}
}

func (s *Spider) Draw(screen tcell.Screen) {
	if s.explosionTicks > 0 {
		s.drawExplosion(screen)
		return
	}
	width, height := screen.Size()
	
	// Draw drop silk (not collidable) - same color as spokes
	spokeColor := tcell.ColorDarkGray
	for _, p := range s.dropSilk {
		if p.X >= 0 && p.Y >= 0 && p.X < width && p.Y < height {
			screen.SetContent(p.X, p.Y, '|', nil, tcell.StyleDefault.Foreground(spokeColor))
		}
	}
	
	// Draw web spokes first (radial lines)
	for _, p := range s.webSpokes {
		if p.X >= 0 && p.Y >= 0 && p.X < width && p.Y < height {
			screen.SetContent(p.X, p.Y, '|', nil, tcell.StyleDefault.Foreground(spokeColor))
		}
	}
	
	// Draw web segments (concentric circles)
	webColor := tcell.ColorGray
	for _, p := range s.webSegments {
		if p.X >= 0 && p.Y >= 0 && p.X < width && p.Y < height {
			screen.SetContent(p.X, p.Y, '-', nil, tcell.StyleDefault.Foreground(webColor))
		}
	}
	
	if !s.active {
		return
	}
	
	cx := int(math.Round(s.x))
	cy := int(math.Round(s.y))
	if cx < 0 || cy < 0 || cx >= width || cy >= height {
		return
	}

	fg := s.Color
	if fg == tcell.ColorDefault || fg == 0 {
		fg = randomSpiderColor()
		s.Color = fg
	}

	// body
	s.drawCell(screen, cx, cy, 'o', fg, width, height)

	// legs alternate with phase
	legOffset := int(math.Round(math.Sin(s.legPhase)))
	s.drawCell(screen, cx-1, cy+legOffset, '-', fg, width, height)
	s.drawCell(screen, cx+1, cy-legOffset, '-', fg, width, height)
	s.drawCell(screen, cx-1, cy-1+legOffset, '/', fg, width, height)
	s.drawCell(screen, cx+1, cy-1-legOffset, '\\', fg, width, height)
	s.drawCell(screen, cx-1, cy+1+legOffset, '\\', fg, width, height)
	s.drawCell(screen, cx+1, cy+1-legOffset, '/', fg, width, height)
}

func (s *Spider) Hit(x, y int) {
	s.active = false
	s.explosionTicks = 6
	s.explosionX = x
	s.explosionY = y
	s.respawnWait = 40 + rand.Intn(80)
	s.webSegments = []Point{}
	s.webSpokes = []Point{}
	s.dropSilk = []Point{}
}

func (s *Spider) HitPoint(width, height int) (int, int, bool) {
	if !s.active {
		return 0, 0, false
	}
	cx := int(math.Round(s.x))
	cy := int(math.Round(s.y))
	if cx < 0 || cy < 0 || cx >= width || cy >= height {
		return 0, 0, false
	}
	return cx, cy, true
}

func (s *Spider) GetWebPoints() []Point {
	points := make([]Point, 0, len(s.webSegments)+len(s.webSpokes))
	points = append(points, s.webSegments...)
	points = append(points, s.webSpokes...)
	return points
}

func (s *Spider) HuntPrey(x, y float64) {
	if s.webState == "done" || s.webState == "building" {
		s.preyX = x
		s.preyY = y
		// Remember if web was still being built
		s.webIncomplete = (s.webState == "building")
		s.webState = "hunting"
	}
}

func (s *Spider) IsHunting() bool {
	return s.webState == "hunting" || s.webState == "eating"
}

func (s *Spider) GetCenterPoint() (float64, float64) {
	return s.centerX, s.centerY
}

func (s *Spider) drawCell(screen tcell.Screen, x, y int, r rune, fg tcell.Color, width, height int) {
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	screen.SetContent(x, y, r, nil, tcell.StyleDefault.Foreground(fg))
}

func (s *Spider) drawExplosion(screen tcell.Screen) {
	width, height := screen.Size()
	fg := color.Yellow
	if s.explosionTicks <= 2 {
		fg = color.Red
	}
	center := Point{X: s.explosionX, Y: s.explosionY}
	for _, p := range []Point{
		center,
		{X: center.X - 1, Y: center.Y},
		{X: center.X + 1, Y: center.Y},
		{X: center.X, Y: center.Y - 1},
		{X: center.X, Y: center.Y + 1},
		{X: center.X - 1, Y: center.Y - 1},
		{X: center.X + 1, Y: center.Y - 1},
		{X: center.X - 1, Y: center.Y + 1},
		{X: center.X + 1, Y: center.Y + 1},
	} {
		if p.X < 0 || p.Y < 0 || p.X >= width || p.Y >= height {
			continue
		}
		screen.SetContent(p.X, p.Y, tcell.RuneBullet, nil, tcell.StyleDefault.Foreground(fg))
	}
}

func (s *Spider) initSpider(width, height int) {
	s.active = true
	// Start at top of screen
	s.x = spiderRandRange(float64(width/4), float64(3*width/4))
	s.y = 0
	// Drop down to middle area
	s.dropTargetY = spiderRandRange(float64(height/4), float64(3*height/4))
	s.webState = "dropping"
	s.pauseTicks = 0
	s.legPhase = spiderRandRange(0, math.Pi*2)
	if s.Color == tcell.ColorDefault || s.Color == 0 {
		s.Color = randomSpiderColor()
	}
	s.webSegments = []Point{}
	s.webSpokes = []Point{}
	s.dropSilk = []Point{}
	s.webBuildStep = 0
	s.webIncomplete = false
}

func NewSpider(cfg SpiderConfig) *Spider {
	if cfg.InitialDelayMax <= 0 {
		cfg.InitialDelayMax = 60
	}
	return &Spider{
		Color:           cfg.Color,
		initialDelayMax: cfg.InitialDelayMax,
	}
}

func randomSpiderColor() tcell.Color {
	colors := []tcell.Color{
		color.Gray,
		color.DarkGray,
		color.Maroon,
		color.Brown,
		color.DarkRed,
	}
	return colors[rand.Intn(len(colors))]
}

func spiderRandRange(minV, maxV float64) float64 {
	return minV + rand.Float64()*(maxV-minV)
}

func ensureSpiderSeeded() {
	if !spiderRandSeeded {
		rand.Seed(time.Now().UnixNano())
		spiderRandSeeded = true
	}
}
