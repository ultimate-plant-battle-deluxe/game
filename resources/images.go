package resources

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var Images struct {
	Stages struct {
		Basic struct {
			Day *ebiten.Image
			Night *ebiten.Image
			Highlight *ebiten.Image
		}
	}
	Plants struct {
		Flowers struct {
			Basic *ebiten.Image
		}
	}
	Items struct {
		Water *ebiten.Image
		Leaf *ebiten.Image
		Trowel *ebiten.Image
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
	Clock struct {
		Face *ebiten.Image
		Hand *ebiten.Image
	}
}

//go:embed static/images/stages/basic.png
var embedStageBasic []byte
//go:embed static/images/stages/basic-night.png
var embedStageBasicNight []byte
//go:embed static/images/stages/basic-highlight.png
var embedStageBasicHighlight []byte

//go:embed static/images/flowers/basic/basic.png
var embedflowersBasicBasic []byte

//go:embed static/images/items/water.png
var embedItemsWater []byte
//go:embed static/images/items/leaf.png
var embedItemsLeaf []byte
//go:embed static/images/items/trowel.png
var embedItemsTrowel []byte
//go:embed static/images/items/seeds/basic.png
var embedItemsSeedsBasic []byte
//go:embed static/images/stages/scenery/cloud1.png
var embedStagesSceneryCloud1 []byte
//go:embed static/images/stages/scenery/cloud2.png
var embedStagesSceneryCloud2 []byte
//go:embed static/images/stages/scenery/cloud3.png
var embedStagesSceneryCloud3 []byte
//go:embed static/images/stages/scenery/cloud4.png
var embedStagesSceneryCloud4 []byte
//go:embed static/images/stages/scenery/cloud1shadow.png
var embedStagesSceneryCloud1shadow []byte
//go:embed static/images/stages/scenery/cloud2shadow.png
var embedStagesSceneryCloud2shadow []byte
//go:embed static/images/stages/scenery/cloud3shadow.png
var embedStagesSceneryCloud3shadow []byte
//go:embed static/images/stages/scenery/cloud4shadow.png
var embedStagesSceneryCloud4shadow []byte
//go:embed static/images/contexts/1water.png
var embedContexts1water []byte
//go:embed static/images/patches/dirt1.png
var embedPatchesDirt1 []byte
//go:embed static/images/clock/face.png
var embedClockFace []byte
//go:embed static/images/clock/hand.png
var embedClockHand []byte

var Fonts struct {
	Becak struct {
		Outline font.Face
		Solid font.Face
	}
}
//go:embed static/fonts/becak.ttf
var embedFontBecak []byte
//go:embed static/fonts/becak-outline.ttf
var embedFontBecakOutline []byte

func Init() {
	var err error

	Images.Stages.Basic.Day, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStageBasic))
	if err != nil {
		log.Fatal(err)
	}
	Images.Stages.Basic.Night, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStageBasicNight))
	if err != nil {
		log.Fatal(err)
	}
	Images.Stages.Basic.Highlight, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStageBasicHighlight))
	if err != nil {
		log.Fatal(err)
	}
	Images.Plants.Flowers.Basic, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedflowersBasicBasic))
	if err != nil {
		log.Fatal(err)
	}
	Images.Items.Water, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedItemsWater))
	if err != nil {
		log.Fatal(err)
	}
	Images.Items.Leaf, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedItemsLeaf))
	if err != nil {
		log.Fatal(err)
	}
	Images.Items.Trowel, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedItemsTrowel))
	if err != nil {
		log.Fatal(err)
	}
	Images.Items.Seeds.Basic, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedItemsSeedsBasic))
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Clouds = make([]*ebiten.Image, 4)
	Images.Scenery.Shadows = make([]*ebiten.Image, 4)
	Images.Scenery.Clouds[0], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStagesSceneryCloud1))
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Clouds[1], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStagesSceneryCloud2))
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Clouds[2], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStagesSceneryCloud3))
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Clouds[3], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStagesSceneryCloud4))
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Shadows[0], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStagesSceneryCloud1shadow))
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Shadows[1], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStagesSceneryCloud2shadow))
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Shadows[2], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStagesSceneryCloud3shadow))
	if err != nil {
		log.Fatal(err)
	}
	Images.Scenery.Shadows[3], _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedStagesSceneryCloud4shadow))
	if err != nil {
		log.Fatal(err)
	}
	Images.Contexts.OneWater, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedContexts1water))
	if err != nil {
		log.Fatal(err)
	}
	Images.Patches.Dirt, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedPatchesDirt1))
	if err != nil {
		log.Fatal(err)
	}
	Images.Clock.Face, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedClockFace))
	if err != nil {
		log.Fatal(err)
	}
	Images.Clock.Hand, _, err = ebitenutil.NewImageFromReader(bytes.NewReader(embedClockHand))
	if err != nil {
		log.Fatal(err)
	}

	ff, _ := opentype.Parse(embedFontBecak)
	Fonts.Becak.Solid, _ = opentype.NewFace(ff, &opentype.FaceOptions{
		Size:    128,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	ff, _ = opentype.Parse(embedFontBecakOutline)
	Fonts.Becak.Outline, _ = opentype.NewFace(ff, &opentype.FaceOptions{
		Size:    128,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}