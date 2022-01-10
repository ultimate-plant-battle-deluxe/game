package main

import "github.com/hajimehoshi/ebiten/v2"

type Plant struct {
	*Sprite
	Garden *Garden
	Kind PlantKind
	Slot int
	MouseOver bool
}

type PlantKind int
const (
	PlantFlowerBasic PlantKind = 0
)


func (p *Plant) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(p.ImageScale, p.ImageScale)
	geom.Translate(float64(p.Garden.X()) + (float64(p.Slot) * 70), float64(p.Garden.Y() - 70))
	screen.DrawImage(p.Image, &ebiten.DrawImageOptions{GeoM: geom})
}