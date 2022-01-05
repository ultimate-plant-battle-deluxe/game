package main

import (
	_ "image/png"
	"log"
	"time"

	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

var stages struct {
	basic *ebiten.Image
}
var flowers struct {
	basic *ebiten.Image
}
var items struct {
	water *ebiten.Image
}

// Sprite represents an image.
type Sprite struct {
	image *ebiten.Image
	x     int
	y     int
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (s *Sprite) In(x, y int) bool {
	// Check the actual color (alpha) value at the specified position
	// so that the result of In becomes natural to users.
	//
	// Note that this is not a good manner to use At for logic
	// since color from At might include some errors on some machines.
	// As this is not so important logic, it's ok to use it so far.
	return s.image.At(x-s.x, y-s.y).(color.RGBA).A > 0
}

func init() {
	ebiten.SetMaxTPS(60)
	rand.Seed(time.Now().UnixNano())
	var err error
	stages.basic, _, err = ebitenutil.NewImageFromFile("images/stages/basic.png")
	if err != nil {
		log.Fatal(err)
	}
	flowers.basic, _, err = ebitenutil.NewImageFromFile("images/flowers/basic/basic.png")
	if err != nil {
		log.Fatal(err)
	}
	items.water, _, err = ebitenutil.NewImageFromFile("images/items/water.png")
	if err != nil {
		log.Fatal(err)
	}
	LoadClouds()
}

type Game struct{}

func (g *Game) Update() error {
	// fmt.Println(ebiten.CursorPosition())
	return nil
}

func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

var tween = gween.New(1, 1.5, 45, ease.InOutSine)
var tweenDirection int = 0

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(stages.basic, nil)
	screen.DrawImage(flowers.basic, nil)
	DrawClouds(screen)
	size, isFinished := tween.Update(1)
	if isFinished {
		if tweenDirection == 0 {
			tween = gween.New(1.5, 1, 45, ease.InOutSine)
			tweenDirection = 1
		} else {
			tween = gween.New(1, 1.5, 45, ease.InOutSine)
			tweenDirection = 0
		}
	}
	geom := ebiten.GeoM{}
	geom.Scale(float64(size), float64(size))
	geom.Translate((-(float64(items.water.Bounds().Max.X) - float64(items.water.Bounds().Max.X)/float64(size)))+100, (-(float64(items.water.Bounds().Max.Y) - float64(items.water.Bounds().Max.Y)/float64(size)))+100)
	screen.DrawImage(items.water, &ebiten.DrawImageOptions{GeoM: geom})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

func main() {
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("Super Plant Friends")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
