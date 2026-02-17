package kitty

import (
	"github.com/gdamore/tcell/v3"
)

type KittyPlayThing interface {
	Update(s tcell.Screen)
	Draw(s tcell.Screen)
}
