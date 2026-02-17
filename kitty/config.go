package kitty

import "github.com/gdamore/tcell/v3"

type SnakeConfig struct {
	MaxLen          int
	Color           tcell.Color
	InitialDelayMax int
}

type SwayStringConfig struct {
	MinLen          int
	MaxLen          int
	Color           tcell.Color
	InitialDelayMax int
}

type ButterflyConfig struct {
	Color           tcell.Color
	InitialDelayMax int
}

type LaserConfig struct {
	Color           tcell.Color
	InitialDelayMax int
}

type SpiderConfig struct {
	Color           tcell.Color
	InitialDelayMax int
}

type KittyConfig struct {
	SnakeCount       int
	SnakeConfig      SnakeConfig
	SwayStringCount  int
	SwayStringConfig SwayStringConfig
	ButterflyCount   int
	ButterflyConfig  ButterflyConfig
	LaserCount       int
	LaserConfig      LaserConfig
	SpiderCount      int
	SpiderConfig     SpiderConfig
}

func DefaultSnakeConfig() SnakeConfig {
	return SnakeConfig{
		MaxLen:          10,
		Color:           tcell.ColorDefault,
		InitialDelayMax: 40,
	}
}

func DefaultSwayStringConfig() SwayStringConfig {
	return SwayStringConfig{
		MinLen:          18,
		MaxLen:          36,
		Color:           tcell.ColorDefault,
		InitialDelayMax: 40,
	}
}

func DefaultButterflyConfig() ButterflyConfig {
	return ButterflyConfig{
		Color:           tcell.ColorDefault,
		InitialDelayMax: 80,
	}
}

func DefaultLaserConfig() LaserConfig {
	return LaserConfig{
		Color:           tcell.ColorDefault,
		InitialDelayMax: 80,
	}
}

func DefaultSpiderConfig() SpiderConfig {
	return SpiderConfig{
		Color:           tcell.ColorDefault,
		InitialDelayMax: 60,
	}
}

func DefaultKittyConfig() KittyConfig {
	return KittyConfig{
		SnakeCount:       2,
		SnakeConfig:      DefaultSnakeConfig(),
		SwayStringCount:  2,
		SwayStringConfig: DefaultSwayStringConfig(),
		ButterflyCount:   1,
		ButterflyConfig:  DefaultButterflyConfig(),
		LaserCount:       1,
		LaserConfig:      DefaultLaserConfig(),
		SpiderCount:      1,
		SpiderConfig:     DefaultSpiderConfig(),
	}
}
