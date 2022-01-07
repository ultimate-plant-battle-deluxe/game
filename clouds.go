package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ultimate-plant-battle-deluxe/game/resources"
)

type Cloud struct {
	*Sprite
	Kind    int
}

var clouds []*Cloud = []*Cloud{}

func LoadClouds() {
}

func createCloud(init bool) *Cloud {
	kind := RandomInt(0, 3)
	x := float64(-resources.Images.Scenery.Clouds[kind].Bounds().Max.X)
	if init {
		x = float64(RandomInt(0, 1080))
	}
	return &Cloud{
		Kind:    kind,
		Sprite: &Sprite{
			X:       x,
			Y:       float64(RandomInt(32, 400)),
			ShadowY: float64(RandomInt(-10, 300)),
			Size:    float64(RandomInt(50, 100)),
			Speed:   float64(RandomInt(50, 150)) / 100,
			Opacity: float64(RandomInt(25, 50)) / 100,
		},
	}
}

var lastCloudTime time.Time
var nextTick time.Duration

func DrawClouds(screen *ebiten.Image) {
	if lastCloudTime.IsZero() {
		lastCloudTime = time.Now()
		startClouds := RandomInt(1, 3)
		for startClouds > 0 {
			clouds = append(clouds, createCloud(true))
			startClouds--
		}
	}
	if lastCloudTime.Add(nextTick).Before(time.Now()) {
		lastCloudTime = time.Now()
		nextTick = time.Duration(RandomInt(int(10000*time.Millisecond), int(20000*time.Millisecond)))
		clouds = append(clouds, createCloud(false))
	}
	for _, cloud := range clouds {
		cloud.X += cloud.Speed
		geom := ebiten.GeoM{}
		geom.Scale(cloud.Size/100, cloud.Size/100)
		geom.Translate(cloud.X, cloud.Y)
		color := ebiten.ColorM{}
		color.Translate(0, 0, 0, -(1 - cloud.Opacity))
		screen.DrawImage(resources.Images.Scenery.Clouds[cloud.Kind], &ebiten.DrawImageOptions{GeoM: geom, ColorM: color})

		sgeom := ebiten.GeoM{}
		sgeom.Scale(cloud.Size/100, cloud.Size/100)
		sgeom.Translate(cloud.X, float64(screen.Bounds().Max.Y)-cloud.ShadowY)
		scolor := ebiten.ColorM{}
		scolor.Translate(0, 0, 0, -(1 - cloud.Opacity*0.25))
		screen.DrawImage(resources.Images.Scenery.Shadows[cloud.Kind], &ebiten.DrawImageOptions{GeoM: sgeom, ColorM: scolor})
		if cloud.X > float64(screen.Bounds().Max.X) {
			clouds = clouds[1:]
		}
	}
}
