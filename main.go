package main

import (
	"encoding/json"
	"fmt"
	_ "image/png"
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
)

type Sprite struct {
	Image *ebiten.Image
	ImageScale float64
	X       float64
	Y       float64
	ShadowY float64
	Size    float64
	Speed   float64
	Opacity float64
}

type Garden struct {
	Slot int
	Water int
	Plants []*Plant
	Highlighted bool
}

func (g *Garden) X() int {
	return (g.Slot + 1) * 400
}

func (g *Garden) Y() int {
	return 800
}

func (g *Garden) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(float64(g.X()), float64(g.Y()))
	colorm := ebiten.ColorM{}
	if g.Highlighted {
		colorm.ChangeHSV(0, 100, 100)
	}
	screen.DrawImage(resources.Images.Patches.Dirt, &ebiten.DrawImageOptions{GeoM: geom, ColorM: colorm})
	for _, plant := range g.Plants {
		plant.Draw(screen)
	}

	// TODO: Enable icon in draw text?
	waterGeom := ebiten.GeoM{}
	waterGeom.Scale(0.25, 0.25)
	waterGeom.Translate(float64(g.X() + 170), float64(g.Y() + 165))
	screen.DrawImage(resources.Images.Items.Water, &ebiten.DrawImageOptions{GeoM: waterGeom})
	DrawText(screen, fmt.Sprint(g.Water), g.X() + 100, g.Y() + 250)
}

var Items []*Item
var Gardens []*Garden

func Api(path string) {
	req, err := http.NewRequest("GET", "http://localhost:8080/v1/" + path, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("x-token", gameToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	token, err := jwt.Parse(resp.Header.Get("X-Token"), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("yolo"), nil
	})
	gameToken = token.Raw
	claims := token.Claims.(jwt.MapClaims);
	sDec, _ := b64.StdEncoding.DecodeString(claims["state"].(string))
	gameState = &GameState{}
	json.Unmarshal([]byte(sDec), gameState)
	Items = []*Item{}
	ApplyGameState()
}

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
	if kind == ItemTrowel {
		img = resources.Images.Items.Trowel
	}
	Items = append(Items, &Item{
		Slot: len(Items) + 1,
		Kind: kind,
		Sprite: &Sprite{
			Image: img,
			ImageScale: 0.5,
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
			ImageScale: 0.5,
		},
		MouseOver: false,
	})
}

func (s *Sprite) In(x, y int) bool {
	scale := s.ImageScale
	if scale == 0 {
		scale = 1
	}
	width := int(float64(s.Image.Bounds().Max.X) * scale)
	height := int(float64(s.Image.Bounds().Max.Y) * scale)
	return x >= int(s.X) && x <= int(s.X) + width && y >= int(s.Y) && y <= int(s.Y) + height
}

func (s *Sprite) Draw() *ebiten.Image {
	w, h := s.Image.Bounds().Max.X, s.Image.Bounds().Max.Y
	scale := s.ImageScale
	if scale == 0 {
		scale = 1
	}
	w = int(float64(w) * scale)
	h = int(float64(h) * scale)
	img := ebiten.NewImage(w, h)
	geom := ebiten.GeoM{}
	geom.Scale(scale, scale)
	img.DrawImage(s.Image, &ebiten.DrawImageOptions{GeoM: geom})
	return img;
}

type GameState struct {
	Id uuid.UUID `json:"id"`
	Time int `json:"time"`
	Items []ItemKind `json:"items"`
	Gardens []Garden `json:"gardens"`
}

var gameState *GameState = &GameState{}
var gameToken string

func ApplyGameState() {
	for _, item := range gameState.Items {
		AddItem(ItemKind(item))
	}

	Gardens = []*Garden{}
	for idx, garden := range gameState.Gardens {
		Gardens = append(Gardens, &Garden{
			Slot: idx,
			Plants: []*Plant{},
			Water: garden.Water,
		})

		for _, plant := range garden.Plants {
			Gardens[idx].AddPlant(plant.Kind)
		}
	}
}

func init() {
	ebiten.SetMaxTPS(60)
	rand.Seed(time.Now().UnixNano())
	resources.Init()
	LoadClouds()

	Items = []*Item{}

	Api("start")
}

type Game struct{}

var draggingItem *Item
var hoveringItem *Item
var hoveringGarden *Garden
var dragOffsetX float64
var dragOffsetY float64

var highlightStage bool

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

		if clock.In(cursorX, cursorY) {
			Api("roll")
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if hoveringGarden != nil {
			if draggingItem != nil {
				switch draggingItem.Kind {
				case ItemWater:
					Api(fmt.Sprintf("water?gardenId=%v", hoveringGarden.Slot))
					break
				case ItemSeedsBasic:
					Api(fmt.Sprintf("plant?gardenId=%v", hoveringGarden.Slot))
				}
			}
			hoveringGarden.Highlighted = false
			hoveringGarden = nil
		}
		if highlightStage {
			Api("garden")
		}
		if draggingItem != nil {
			draggingItem.OnMouseUp()
			draggingItem = nil
		}
	}
	
	if draggingItem != nil && draggingItem.MouseDown {
		draggingItem.X = float64(cursorX) - dragOffsetX
		draggingItem.Y = float64(cursorY) - dragOffsetY
		
	}
	if draggingItem != nil && draggingItem.Kind == ItemTrowel && cursorY >= 680  {
		highlightStage = true
	} else {
		highlightStage = false
	}

	if draggingItem != nil {
		for _, garden := range Gardens {
			if cursorX >= garden.X() && cursorX <= garden.X() + resources.Images.Patches.Dirt.Bounds().Max.X &&
				cursorY >= garden.Y() && cursorY <= garden.Y() + resources.Images.Patches.Dirt.Bounds().Max.Y {
				garden.Highlighted = true
				hoveringGarden = garden
			} else if garden.Highlighted {
				garden.Highlighted = false
				hoveringGarden = nil
			}
		}
	}
	return nil
}

func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Stage
	screen.DrawImage(resources.Images.Stages.Basic.Night, nil)
	
	sceneColor := ebiten.ColorM{}
	sceneColor.Translate(0, 0, 0, (float64(gameState.Time) / 10)-1)
	screen.DrawImage(resources.Images.Stages.Basic.Day, &ebiten.DrawImageOptions{ColorM: sceneColor})
	
	if highlightStage {
		screen.DrawImage(resources.Images.Stages.Basic.Highlight, nil)
	}
	
	// Gardens
	for _, garden := range Gardens {
		garden.Draw(screen)
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

	clock.Draw(screen)
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
