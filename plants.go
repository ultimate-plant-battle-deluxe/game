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
	geom.Translate(p.X + (float64(p.Slot) * 50), p.Y)
	screen.DrawImage(p.Image, &ebiten.DrawImageOptions{GeoM: geom})
}