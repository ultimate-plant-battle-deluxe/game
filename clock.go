package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ultimate-plant-battle-deluxe/game/resources"
)

type Clock struct {
	*Sprite
}

var clock *Clock = &Clock{
	Sprite: &Sprite{
		Image: resources.Images.Clock.Face,
		X: 1560,
		Y: 120,
	},
}

func (c *Clock) Draw(screen *ebiten.Image) {
	handGeom := ebiten.GeoM{}
	handGeom.Translate(-float64(resources.Images.Clock.Hand.Bounds().Max.X)/2, -float64(resources.Images.Clock.Hand.Bounds().Max.Y)/2)
	handGeom.Rotate(float64(gameState.Time) * 30 * 2 * (math.Pi / 360))
	handGeom.Translate(float64(resources.Images.Clock.Hand.Bounds().Max.X)/2, float64(resources.Images.Clock.Hand.Bounds().Max.Y)/2)
	clock.Image = ebiten.NewImage(resources.Images.Clock.Face.Bounds().Max.X, resources.Images.Clock.Face.Bounds().Max.Y)
	clock.Image.DrawImage(resources.Images.Clock.Face, nil)
	clock.Image.DrawImage(resources.Images.Clock.Hand, &ebiten.DrawImageOptions{GeoM: handGeom})
	geom := ebiten.GeoM{}
	geom.Translate(1560, 150)
	screen.DrawImage(clock.Image, &ebiten.DrawImageOptions{GeoM: geom})
	hoursString := "Hours"
	if gameState.Time == 1 {
		hoursString = "Hour"
	}
	DrawText(screen, fmt.Sprintf("%v %s", gameState.Time, hoursString), 1450, 120)
}