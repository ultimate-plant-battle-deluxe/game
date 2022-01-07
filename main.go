package main

import (
	"encoding/json"
	"fmt"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"math/rand"

	b64 "encoding/base64"

	"github.com/golang-jwt/jwt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	uuid "github.com/satori/go.uuid"
	"github.com/ultimate-plant-battle-deluxe/game/resources"
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

type Garden struct {
	Plants []*Plant
}

var Items []*Item
var Gardens []*Garden

func AddItem(kind ItemKind) {
	var img *ebiten.Image
	if kind == ItemWater {
		img = resources.Images.Items.Water
	}
	if kind == ItemLeaf {
		img = resources.Images.Items.Leaf
	}
	if kind == ItemSeedsBasic {
		img = resources.Images.Items.Seeds.Basic
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

func (g *Garden) AddPlant(kind PlantKind) {
	var img *ebiten.Image
	if kind == PlantFlowerBasic {
		img = resources.Images.Plants.Flowers.Basic
	}
	g.Plants = append(g.Plants, &Plant{
		Garden: g,
		Slot: len(g.Plants) + 1,
		Kind: kind,
		Sprite: &Sprite{
			Image: img,
		},
		MouseOver: false,
	})
}

func (s *Sprite) In(x, y int) bool {
	width := s.Image.Bounds().Max.X
	height := s.Image.Bounds().Max.Y
	return x >= int(s.X) && x <= int(s.X) + width && y >= int(s.Y) && y <= int(s.Y) + height
}

var fontInner font.Face
var fontOutline font.Face

type GameState struct {
	Id uuid.UUID `json:"id"`
	Time int `json:"time"`
	Items []ItemKind `json:"items"`
}

var gameState *GameState = &GameState{}

func init() {
	ebiten.SetMaxTPS(60)
	rand.Seed(time.Now().UnixNano())
	resources.Init()
	LoadClouds()

	Items = []*Item{}
	Gardens = []*Garden{{}}

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

	resp, err := http.Get("http://localhost:8080/v1/start")
	if err != nil {
		log.Fatalln(err)
	}
	token, err := jwt.Parse(resp.Header.Get("X-Token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("yolo"), nil
	})

	claims := token.Claims.(jwt.MapClaims);
	sDec, _ := b64.StdEncoding.DecodeString(claims["state"].(string))
	json.Unmarshal([]byte(sDec), gameState)

	for _, item := range gameState.Items {
		if item == ItemWater {
			AddItem(ItemWater)
		}
		if item == ItemLeaf {
			AddItem(ItemLeaf)
		}
		if item == ItemSeedsBasic {
			AddItem(ItemSeedsBasic)
		}
	}
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

	// Stage
	screen.DrawImage(resources.Images.Stages.Basic, nil)

	// Gardens
	for idx, garden := range Gardens {
		geom := ebiten.GeoM{}
		geom.Translate(float64((idx + 1) * 200), 800)
		screen.DrawImage(resources.Images.Patches.Dirt, &ebiten.DrawImageOptions{GeoM: geom})
		for _, plant := range garden.Plants {
			plant.Draw(screen)
		}
	}

	DrawClouds(screen)

	// Items
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
