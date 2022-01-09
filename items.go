package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/ultimate-plant-battle-deluxe/game/resources"
)

type Item struct {
	*Sprite
	Kind ItemKind
	Slot int
	MouseOver bool
	MouseDown bool
	tween *gween.Tween
	tweenDirection int
}
func (i *Item) OnMouseDown() {
	i.MouseDown = true
}

func (i *Item) OnMouseUp() {
	i.MouseDown = false
	i.Snap()
}

type ItemKind int
const (
	ItemLeaf ItemKind = iota
	ItemWater
	ItemSeedsBasic
	ItemTrowel
)

func (i *Item) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	if i.MouseOver {
		i.pulse(&geom)
	}
	geom.Translate(i.X, i.Y)
	screen.DrawImage(i.Sprite.Draw(), &ebiten.DrawImageOptions{GeoM: geom})
}

func (i *Item) pulse(geom *ebiten.GeoM) {
	if i.tween == nil {
		i.tween = gween.New(1, 1.25, 30, ease.InOutSine)
	}
	size, isFinished := i.tween.Update(1)
	i.Size = float64(size)
	if isFinished {
		if i.tweenDirection == 0 {
			i.tween = gween.New(1.25, 1, 30, ease.InOutSine)
			i.tweenDirection = 1
		} else {
			i.tween = gween.New(1, 1.25, 30, ease.InOutSine)
			i.tweenDirection = 0
		}
	}
	
	// xOffset := ((i.X * size) - i.X) / 2
	imageWidth := float64(resources.Images.Items.Water.Bounds().Max.X) * i.Sprite.ImageScale
	imageHeight := float64(resources.Images.Items.Water.Bounds().Max.X)  * i.Sprite.ImageScale
	geom.Scale(float64(size), float64(size))
	geom.Translate(
		-(((imageWidth*float64(size)) - imageWidth)/2) + float64(size),
		-((imageHeight*float64(size)) - imageHeight)/2)
}

func (i *Item) Snap() {
	i.X = float64(275 + (i.Slot * 300))
	i.Y = 400
}

func (i *Item) OnMouseOver() {
	i.MouseOver = true
}
func (i *Item) OnMouseOut() {
	i.MouseOver = false
	i.tween.Reset()
}