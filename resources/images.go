package resources

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var Images struct {
	Stages struct {
		Basic *ebiten.Image
	}
	Plants struct {
		Flowers struct {
			Basic *ebiten.Image
		}
	}
	Items struct {
		Water *ebiten.Image
		Leaf *ebiten.Image
		Seeds struct {
			Basic *ebiten.Image
		}
	}
	Scenery struct {
		Clouds  []*ebiten.Image
		Shadows []*ebiten.Image
	}	
	Contexts struct {
		OneWater *ebiten.Image
	}
	Patches struct {
		Dirt *ebiten.Image
	}
}


func Init() {
	var err error
	Images.Stages.Basic, _, err = ebitenutil.NewImageFromFile("static/images/stages/basic.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Plants.Flowers.Basic, _, err = ebitenutil.NewImageFromFile("static/images/flowers/basic/basic.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Items.Water, _, err = ebitenutil.NewImageFromFile("static/images/items/water.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Items.Leaf, _, err = ebitenutil.NewImageFromFile("static/images/items/leaf.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Items.Seeds.Basic, _, err = ebitenutil.NewImageFromFile("static/images/items/seeds/basic.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Clouds = make([]*ebiten.Image, 4)
	Images.Scenery.Shadows = make([]*ebiten.Image, 4)
	Images.Scenery.Clouds[0], _, err = ebitenutil.NewImageFromFile("static/images/stages/scenery/cloud1.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Clouds[1], _, err = ebitenutil.NewImageFromFile("static/images/stages/scenery/cloud2.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Clouds[2], _, err = ebitenutil.NewImageFromFile("static/images/stages/scenery/cloud3.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Clouds[3], _, err = ebitenutil.NewImageFromFile("static/images/stages/scenery/cloud4.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Shadows[0], _, err = ebitenutil.NewImageFromFile("static/images/stages/scenery/cloud1shadow.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Shadows[1], _, err = ebitenutil.NewImageFromFile("static/images/stages/scenery/cloud2shadow.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Shadows[2], _, err = ebitenutil.NewImageFromFile("static/images/stages/scenery/cloud3shadow.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Shadows[3], _, err = ebitenutil.NewImageFromFile("static/images/stages/scenery/cloud4shadow.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Contexts.OneWater, _, err = ebitenutil.NewImageFromFile("static/images/contexts/1water.png")
	if err != nil {
		log.Fatal(err)
	}
	Images.Patches.Dirt, _, err = ebitenutil.NewImageFromFile("static/images/patches/dirt1.png")
	if err != nil {
		log.Fatal(err)
	}
}