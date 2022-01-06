package main

import (
	_ "image/png"
	"io/ioutil"
	"log"
	"time"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/warent/plant-friends/resources"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Sprite struct {
	Image *ebiten.Image
	X       float64
	Y       float64
	ShadowY float64
	Size    float64
	Speed   float64
	Opacity float64
}

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
	ItemLeaf ItemKind = 0
	ItemWater = 1
)

var Items []*Item

func AddItem(kind ItemKind) {
	var img *ebiten.Image
	if kind == ItemWater {
		img = resources.Images.Items.Water
	}
	if kind == ItemLeaf {
		img = resources.Images.Items.Leaf
	}
	Items = append(Items, &Item{
		Slot: len(Items) + 1,
		Kind: kind,
		Sprite: &Sprite{
			Image: img,
		},
		MouseOver: false,
	})
	Items[len(Items)-1].Snap()
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

var defaultItemY int = 400

func (i *Item) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	if i.MouseOver {
		i.pulse(&geom)
	}
	geom.Translate(i.X, i.Y)
	screen.DrawImage(i.Image, &ebiten.DrawImageOptions{GeoM: geom})
}

func (i *Item) pulse(geom *ebiten.GeoM) {
	if i.tween == nil {
		i.tween = gween.New(1, 1.5, 30, ease.InOutSine)
	}
	size, isFinished := i.tween.Update(1)
	i.Size = float64(size)
	if isFinished {
		if i.tweenDirection == 0 {
			i.tween = gween.New(1.5, 1, 30, ease.InOutSine)
			i.tweenDirection = 1
		} else {
			i.tween = gween.New(1, 1.5, 30, ease.InOutSine)
			i.tweenDirection = 0
		}
	}
	
	// xOffset := ((i.X * size) - i.X) / 2
	imageWidth := float64(resources.Images.Items.Water.Bounds().Max.X)
	imageHeight := float64(resources.Images.Items.Water.Bounds().Max.X)
	geom.Scale(float64(size), float64(size))
	geom.Translate(
		-(((imageWidth*float64(size)) - imageWidth)/2) + float64(size),
		-((imageHeight*float64(size)) - imageHeight)/2)
}

func (s *Sprite) In(x, y int) bool {
	width := s.Image.Bounds().Max.X
	height := s.Image.Bounds().Max.Y
	return x >= int(s.X) && x <= int(s.X) + width && y >= int(s.Y) && y <= int(s.Y) + height
}

var fontInner font.Face
var fontOutline font.Face

func init() {
	ebiten.SetMaxTPS(60)
	rand.Seed(time.Now().UnixNano())
	resources.Init()
	LoadClouds()

	Items = []*Item{}
	AddItem(ItemWater)
	AddItem(ItemWater)
	AddItem(ItemLeaf)

	fontFile, _ := ioutil.ReadFile("static/fonts/roboto.ttf")
	ff, _ := opentype.Parse(fontFile)
	fontInner, _ = opentype.NewFace(ff, &opentype.FaceOptions{
		Size:    64,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	fontOutline, _ = opentype.NewFace(ff, &opentype.FaceOptions{
		Size:    74,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

type Game struct{}

var draggingItem *Item
var hoveringItem *Item
var dragOffsetX float64
var dragOffsetY float64

func (g *Game) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()

	if draggingItem == nil {
		for _, item := range Items {
			mouseOver := item.In(cursorX, cursorY)
			if mouseOver && !item.MouseOver {
				item.OnMouseOver()
				hoveringItem = item
			} else if !mouseOver && item.MouseOver && !item.MouseDown {
				item.OnMouseOut()
				hoveringItem = nil
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, item := range Items {
			mouseOver := item.In(cursorX, cursorY)
			if mouseOver {
				item.OnMouseDown()
				draggingItem = item
				dragOffsetX = float64(cursorX) - item.X
				dragOffsetY = float64(cursorY) - item.Y
				break
			}
		} 
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, item := range Items {
			mouseOver := item.In(cursorX, cursorY)
			if mouseOver && item.MouseDown {
				item.OnMouseUp()
				draggingItem = nil
				break
			}
		} 
	}

	if draggingItem != nil && draggingItem.MouseDown {
		draggingItem.X = float64(cursorX) - dragOffsetX
		draggingItem.Y = float64(cursorY) - dragOffsetY
	}
	return nil
}

func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(resources.Images.Stages.Basic, nil)

	geom := ebiten.GeoM{}
	geom.Translate(200, 800)
	screen.DrawImage(resources.Images.Patches.Dirt, &ebiten.DrawImageOptions{GeoM: geom})

	geom2 := ebiten.GeoM{}
	geom2.Translate(300, 740)
	screen.DrawImage(resources.Images.Flowers.Basic, &ebiten.DrawImageOptions{GeoM: geom2})
	
	DrawClouds(screen)
	for _, item := range Items {
		if item == draggingItem {
			continue
		}
		item.Draw(screen)
	}

	
	if hoveringItem != nil && hoveringItem.Kind == ItemWater {
		geom := ebiten.GeoM{}
		geom.Translate(20, 20)
		screen.DrawImage(resources.Images.Contexts.OneWater, &ebiten.DrawImageOptions{GeoM: geom})
	}
	
	if draggingItem != nil {
		draggingItem.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}

func main() {
	ebiten.SetWindowSize(960, 540)
	ebiten.SetWindowTitle("Ultimate Plant Battle Deluxe")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
