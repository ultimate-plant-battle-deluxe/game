package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var scenery struct {
	clouds  []*ebiten.Image
	shadows []*ebiten.Image
}

type Cloud struct {
	Kind    int
	X       float64
	Y       float64
	ShadowY float64
	Size    float64
	Speed   float64
	Opacity float64
}

var clouds []*Cloud = []*Cloud{}

func LoadClouds() {
	scenery.clouds = make([]*ebiten.Image, 4)
	scenery.shadows = make([]*ebiten.Image, 4)
	var err error
	scenery.clouds[0], _, err = ebitenutil.NewImageFromFile("images/stages/scenery/cloud1.png")
	if err != nil {
		log.Fatal(err)
	}
	scenery.clouds[1], _, err = ebitenutil.NewImageFromFile("images/stages/scenery/cloud2.png")
	if err != nil {
		log.Fatal(err)
	}
	scenery.clouds[2], _, err = ebitenutil.NewImageFromFile("images/stages/scenery/cloud3.png")
	if err != nil {
		log.Fatal(err)
	}
	scenery.clouds[3], _, err = ebitenutil.NewImageFromFile("images/stages/scenery/cloud4.png")
	if err != nil {
		log.Fatal(err)
	}
	scenery.shadows[0], _, err = ebitenutil.NewImageFromFile("images/stages/scenery/cloud1shadow.png")
	if err != nil {
		log.Fatal(err)
	}
	scenery.shadows[1], _, err = ebitenutil.NewImageFromFile("images/stages/scenery/cloud2shadow.png")
	if err != nil {
		log.Fatal(err)
	}
	scenery.shadows[2], _, err = ebitenutil.NewImageFromFile("images/stages/scenery/cloud3shadow.png")
	if err != nil {
		log.Fatal(err)
	}
	scenery.shadows[3], _, err = ebitenutil.NewImageFromFile("images/stages/scenery/cloud4shadow.png")
	if err != nil {
		log.Fatal(err)
	}
}

func createCloud(init bool) *Cloud {
	kind := RandomInt(0, 3)
	x := float64(-scenery.clouds[kind].Bounds().Max.X)
	if init {
		x = float64(RandomInt(0, 1080))
	}
	return &Cloud{
		Kind:    kind,
		X:       x,
		Y:       float64(RandomInt(32, 400)),
		ShadowY: float64(RandomInt(0, 10)),
		Size:    float64(RandomInt(50, 100)),
		Speed:   float64(RandomInt(50, 150)) / 100,
		Opacity: float64(RandomInt(25, 50)) / 100}
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
		geom.Translate(cloud.X, cloud.Y)
		geom.Scale(cloud.Size/100, cloud.Size/100)
		color := ebiten.ColorM{}
		color.Translate(0, 0, 0, -(1 - cloud.Opacity))
		screen.DrawImage(scenery.clouds[cloud.Kind], &ebiten.DrawImageOptions{GeoM: geom, ColorM: color})

		sgeom := ebiten.GeoM{}
		sgeom.Translate(cloud.X, float64(screen.Bounds().Max.Y)-cloud.ShadowY)
		sgeom.Scale(cloud.Size/100, cloud.Size/100)
		scolor := ebiten.ColorM{}
		scolor.Translate(0, 0, 0, -(1 - cloud.Opacity*0.25))
		screen.DrawImage(scenery.shadows[cloud.Kind], &ebiten.DrawImageOptions{GeoM: sgeom, ColorM: scolor})
		if cloud.X > float64(screen.Bounds().Max.X) {
			clouds = clouds[1:]
		}
	}
}
