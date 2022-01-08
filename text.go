package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/ultimate-plant-battle-deluxe/game/resources"
)

func DrawText(dst *ebiten.Image, str string, x, y int) {
	text.Draw(dst, str, resources.Fonts.Becak.Solid, x, y, color.White)
	text.Draw(dst, str, resources.Fonts.Becak.Outline, x + 4, y + 4, color.Black)
}