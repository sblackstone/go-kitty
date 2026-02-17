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
}

func (k *Kitty) EventLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-k.s.EventQ():
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					return
				}
			}
		}
	}
}

func (k *Kitty) Play(ctx context.Context) {
	//k.objects = append(k.objects, &BouncySquare{X1: 0, Len: 2, Vx: 1, Vy: 1})
	k.objects = append(k.objects, &Snake{MaxLen: 10})
	for {
		select {
		case <-ctx.Done():
			return
		default:
			k.s.Clear()
			for _, o := range k.objects {
				o.Update(k.s)
			}
			for _, o := range k.objects {
				o.Draw(k.s)
			}
			k.s.Show()
			time.Sleep(100 * time.Millisecond)
		}

	}
}

func (k *Kitty) Start(ctx context.Context) {
	go k.EventLoop(ctx)
	k.Play(ctx)
}

func New() (*Kitty, error) {

	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}

	s.SetStyle(DEFAULT_STYLE)

	width, height := s.Size()

	return &Kitty{
		screenWidth:  width,
		screenHeight: height,
		s:            s,
	}, nil
}
