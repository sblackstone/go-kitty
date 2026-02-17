package kitty

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var DEFAULT_STYLE = tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

type Point struct {
	X int
	Y int
}

type KittyString struct {
	MaxLen int
	Point  Point
}

type Kitty struct {
	screenWidth  int
	screenHeight int
	s            tcell.Screen
	objects      []KittyPlayThing
	config       KittyConfig
}

func (k *Kitty) EventLoop(ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		ev := <- k.s.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				cancel()
				return
			}
		case *tcell.EventInterrupt:
			return
		}
	}
}

func (k *Kitty) Play(ctx context.Context) {
	//k.objects = append(k.objects, &BouncySquare{X1: 0, Len: 2, Vx: 1, Vy: 1})
	k.objects = k.objects[:0]
	cfg := k.config
	if cfg.SnakeCount < 0 {
		cfg.SnakeCount = 0
	}
	if cfg.SwayStringCount < 0 {
		cfg.SwayStringCount = 0
	}
	if cfg.ButterflyCount < 0 {
		cfg.ButterflyCount = 0
	}
	if cfg.LaserCount < 0 {
		cfg.LaserCount = 0
	}
	if cfg.SpiderCount < 0 {
		cfg.SpiderCount = 0
	}
	for i := 0; i < cfg.SnakeCount; i++ {
		k.objects = append(k.objects, NewSnake(cfg.SnakeConfig))
	}
	for i := 0; i < cfg.SwayStringCount; i++ {
		k.objects = append(k.objects, NewSwayString(cfg.SwayStringConfig))
	}
	for i := 0; i < cfg.ButterflyCount; i++ {
		k.objects = append(k.objects, NewButterfly(cfg.ButterflyConfig))
	}
	for i := 0; i < cfg.LaserCount; i++ {
		k.objects = append(k.objects, NewLaserPointer(cfg.LaserConfig))
	}
	for i := 0; i < cfg.SpiderCount; i++ {
		k.objects = append(k.objects, NewSpider(cfg.SpiderConfig))
	}


	for {
		select {
		case <-ctx.Done():
			return
		default:
			k.s.Clear()
			for _, o := range k.objects {
				o.Update(k.s)
			}
			k.handleLaserHits()
			for _, o := range k.objects {
				o.Draw(k.s)
			}
			k.s.Show()
			time.Sleep(55 * time.Millisecond)
		}

	}
}

func (k *Kitty) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer k.s.Fini()
	go k.EventLoop(ctx, cancel)
	k.Play(ctx)
}

func New(config KittyConfig) (*Kitty, error) {

	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}

	s.SetStyle(DEFAULT_STYLE)

	width, height := s.Size()

	if (config == KittyConfig{}) {
		config = DefaultKittyConfig()
	}

	return &Kitty{
		screenWidth:  width,
		screenHeight: height,
		s:            s,
		config:       config,
	}, nil
}

func (k *Kitty) handleLaserHits() {
	width, height := k.s.Size()
	lasers := make([]Point, 0)
	for _, o := range k.objects {
		if l, ok := o.(*LaserPointer); ok {
			if p, ok := l.Position(width, height); ok {
				lasers = append(lasers, p)
			}
		}
	}
	if len(lasers) == 0 {
		return
	}
	for _, o := range k.objects {
		b, ok := o.(*Butterfly)
		if !ok {
			continue
		}
		bx, by, ok := b.HitPoint(width, height)
		if !ok {
			continue
		}
		for _, p := range lasers {
			if absInt(p.X-bx) <= 1 && absInt(p.Y-by) <= 1 {
				b.Hit(bx, by)
				break
			}
		}
	}
	for _, o := range k.objects {
		s, ok := o.(*Spider)
		if !ok {
			continue
		}
		sx, sy, ok := s.HitPoint(width, height)
		if !ok {
			continue
		}
		for _, p := range lasers {
			if absInt(p.X-sx) <= 1 && absInt(p.Y-sy) <= 1 {
				s.Hit(sx, sy)
				break
			}
		}
	}
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
