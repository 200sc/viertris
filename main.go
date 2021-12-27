package main

import (
	"fmt"
	"os"
	"image"
	"image/draw"
	"image/color"
	"image/gif"
	"math/rand"
	"sort"
	"time"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/alg/intgeom"
)

const MenuSceneName = "menu"
const GameSceneName = "game"

func main() {
	// Pseudocode:
	// Initialize scenes:
	// Menu Scene:
	// TODO Oak v4: make scene an interface? 
	oak.AddScene(MenuSceneName, scene.Scene{
		Start: func(ctx *scene.Context) {
			ctx.Window.GoToScene(GameSceneName)
			// TODO:
			// -- Start game button -> goto game scene
			// -- Online multiplayer (yeah right)
			// -- High Scores 
			// -- Options button -> replace button set with
			// --- Master volume
			// --- Music volume
			// --- SFX volume
			// --- window size 
			// --- enable window resize + autoscaling 
			// --- keymapping 
			// --- Back 
			// -- Exit button
			// - Controllable via keyboard, joystick, mouse 
		},
	})
	oak.AddScene(GameSceneName, scene.Scene{
		Start: func(ctx *scene.Context) {
			//recordGif(ctx)
			rand.Seed(time.Now().Unix())
			st := NewGameState(ctx, GameConfig{}) 
			ctx.DrawStack.Draw(st)
			//populateTestBoard(st.GameBoard)
			const keyRepeatDuration = 70 * time.Millisecond
			const todoFallDuration = 700 * time.Millisecond 
			dropAt := time.Now().Add(todoFallDuration)
			

			st.ActiveTris = ActiveTris{
				Board: &st.GameBoard,
				X: 5,
				Y: 0, 
				TrisKind: RandomKind(),
			}

			keyRepeat := time.Now().Add(keyRepeatDuration)
			ctx.EventHandler.GlobalBind(event.Enter, func(_ event.CID, payload interface{}) int {
				tileDone := false 
				//enter := payload.(event.EnterPayload)
				if time.Now().After(dropAt) {
					tileDone = st.ActiveTris.MoveDown()
					dropAt = time.Now().Add(todoFallDuration)
				}
				if time.Now().After(keyRepeat) {
					if ctx.KeyState.IsDown(key.A) {
						st.ActiveTris.MoveLeft()
						keyRepeat = time.Now().Add(keyRepeatDuration)
					}
					if ctx.KeyState.IsDown(key.D) {
						st.ActiveTris.MoveRight()
						keyRepeat = time.Now().Add(keyRepeatDuration)
					}
					if ctx.KeyState.IsDown(key.S) {
						tileDone = st.ActiveTris.MoveDown()
						dropAt = time.Now().Add(todoFallDuration)
						keyRepeat = time.Now().Add(keyRepeatDuration/2)
					}
					if ctx.KeyState.IsDown(key.Q) {
						st.ActiveTris.Rotation = st.ActiveTris.RotateLeft()
						keyRepeat = time.Now().Add(keyRepeatDuration)
					}
					if ctx.KeyState.IsDown(key.E) {
						st.ActiveTris.Rotation = st.ActiveTris.RotateRight()
						keyRepeat = time.Now().Add(keyRepeatDuration)
					}
				}
				if tileDone {
					st.GameBoard.SetActiveTile()
					st.ActiveTris = ActiveTris{
						Board: &st.GameBoard,
						X: 5,
						Y: 0, 
						TrisKind: RandomKind(),
					}
				}
				return 0
			})
			ctx.EventHandler.GlobalBind(key.Down+key.R, event.Empty(func() {
				recordGif(ctx)
			}))
	// Game Scene:
	// -- Score / Level tracking 
	// -- Stored / Preview window 
	// ---- store current tris 
	// ---- retrieve stored tris 
		},
	})
	err := oak.Init(MenuSceneName, func(c oak.Config) (oak.Config, error) {
		c.Title = "Viertris"
		return c, nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 

var _ render.Renderable = &GameState{}

func NewGameState(ctx *scene.Context, cfg GameConfig) *GameState {
	return &GameState{
		LayeredPoint: render.NewLayeredPoint(0,0,0),
		ctx: ctx, 
		w: ctx.Window.Width(),
		h: ctx.Window.Height(),
		GameBoard: NewGameBoard(cfg),
	}
	// TODO: on ctx window size change, update w and h 
}

type GameState struct {
	render.LayeredPoint
	GameBoard 
	GameConfig
	ctx *scene.Context 
	w, h int 
	Clears uint32
	Level uint16
	NextTris TrisKind 
	StoredTris TrisKind 
}

const boardRatio = 0.7
const buffer = 5
const boardX = buffer
const boardY = buffer

func (gs *GameState) Draw(buff draw.Image, _ float64, _ float64) {
	
	
	// board outline 
	boardW := int(boardRatio * float64(gs.w))
	gs.GameBoard.draw(buff, boardW, gs.h)
	

	// sidebar
	// preview / next outlines 

	const sidebarRatio = 1 - boardRatio
	const tilePreviewRatio = 0.3
	sidebarW := int(sidebarRatio * float64(gs.w))
	previewH := int(tilePreviewRatio * float64(gs.h))
	drawRect(
		buff.(*image.RGBA), 
		intgeom.Point2{
			boardX+boardW,
			boardY,
		},
		intgeom.Point2{
			sidebarW - buffer*2,
			previewH - buffer*2,
		},
		color.RGBA{255,255,255,255})

	// score outline

	const scoreRatio = 1 - tilePreviewRatio
	scoreH := int(scoreRatio * float64(gs.h))
	drawRect(
		buff.(*image.RGBA), 
		intgeom.Point2{
			boardX+int(boardRatio * float64(gs.w)),
			boardY+previewH,
		},
		intgeom.Point2{
			sidebarW - buffer*2,
			scoreH - buffer*2,
		},
		color.RGBA{255,255,255,255})

}

func (gs *GameState) GetDims() (int, int) {
	return gs.w, gs.h
}

type GameConfig struct {
	// todo
}

type Rotation uint8

const (
	Rotation0 Rotation = iota 
	Rotation90 Rotation = iota 
	Rotation180 Rotation = iota 
	Rotation270 Rotation = iota 
	RotationMax Rotation = iota 
)

func (r Rotation) RotateRight() Rotation {
	r++
	if r > Rotation270 {
		return Rotation0
	}
	return r
}

func (r Rotation) RotateLeft() Rotation {
	r--
	if r == 255 {
		return Rotation270
	}
	return r 
}

type ActiveTris struct {
	Board *GameBoard
	Rotation 
	X BoardDimension 
	Y BoardDimension 
	TrisKind
}

// Offsets on an ActiveTris takes into account rotation
func (at *ActiveTris) Offsets() [4][2]int8{
	rawOff := at.TrisKind.Offsets()
	switch at.Rotation {
	case Rotation270:
		for i, off := range rawOff {
			off[0], off[1] = off[1], -1 * off[0]
			rawOff[i] = off 
		}
		fallthrough 
	case Rotation180:
		for i, off := range rawOff {
			off[0], off[1] = off[1], -1 * off[0]
			rawOff[i] = off 
		}
		fallthrough
	case Rotation90:
		for i, off := range rawOff {
			off[0], off[1] = off[1], -1 * off[0]
			rawOff[i] = off 
		}
		fallthrough
	case Rotation0:
		return rawOff
	default:
		panic(fmt.Sprintf("invalid rotation", at.Rotation))
	}
}

func (at *ActiveTris) MoveLeft() {
	minX := int(at.X) 
	off := at.Offsets()
	for _, o := range off {
		x := int(at.X) + int(o[0])
		y := int(at.Y) + int(o[1])
		if at.Board.IsSet(x-1,y) {
			return 
		}
		if x < minX {
			minX = x 
		}
	}
	if minX > 0 {
		at.X--
	}
}

func (at *ActiveTris) MoveRight() {
	maxX := int(at.X) 
	off := at.Offsets()
	for _, o := range off {
		x := int(at.X) + int(o[0])
		y := int(at.Y) + int(o[1])
		if at.Board.IsSet(x+1,y) {
			return 
		}
		if x > maxX {
			maxX = x 
		}
	}
	if maxX < int(at.Board.Width-1) {
		at.X++
	}
}

func (at *ActiveTris) MoveDown() bool {
	maxY := int(at.Y)
	off := at.Offsets()
	for _, o := range off {
		y := int(at.Y) + int(o[1])
		if y > maxY {
			maxY = y
		}
	}
	if maxY <= int(at.Board.Height-1) {
		placed := at.Board.CheckIfTileIsPlaced()
		if !placed {
			at.Y++
		}
		return placed 
	}
	return false 
}

type TrisKind uint8

const (
	KindNone TrisKind = iota 
	KindT TrisKind = iota
	KindLine TrisKind = iota
	KindSquare TrisKind = iota
	KindZ TrisKind = iota
	KindL TrisKind = iota
	KindJ TrisKind = iota
	KindFinal TrisKind = iota 
)

var kindColors = []color.RGBA{
	KindNone: {}, 
	KindT: {200,0,0,255},
	KindLine: {0,200,0,255},
	KindSquare: {0,0,200,255},
	KindZ: {200, 200, 0, 255},
	KindL: {200, 0, 200, 255},
	KindJ: {0, 200, 200, 255},
}

func (tk TrisKind) Color() color.RGBA {
	return kindColors[tk]
}

var kindOffsets = [][4][2]int8{
	KindT: [4][2]int8{
		{0,0},
		{-1,0},
		{0,-1},
		{1,0},
	},
	KindLine: [4][2]int8{
		{0,0},
		{0,-1},
		{0,1},
		{0,2},
	},
	KindSquare: [4][2]int8{
		{0,0},
		{0,1},
		{1,1},
		{1,0},
	},
	KindZ: [4][2]int8{
		{0,0},
		{-1,0},
		{0,1},
		{1,1},
	},
	KindL: [4][2]int8{
		{0,0},
		{0,-1},
		{0,1},
		{1,1},
	},
	KindJ: [4][2]int8{
		{0,0},
		{0,-1},
		{0,1},
		{-1,1},
	},
} 

func (tk TrisKind) Offsets() [4][2]int8{
	return kindOffsets[tk]
}

type GameBoard struct {
	Width BoardDimension 
	Height BoardDimension 
	// Unexpected! Y x X! Because that makes it easier to clear lines! 
	Set [][]TrisKind
	ActiveTris ActiveTris
}

func NewGameBoard(cfg GameConfig) GameBoard {
	// TODO 
	const todoWidth = 10
	const todoHeight = 24 
	set := make([][]TrisKind, todoHeight)
	for y := range set {
		set[y] = make([]TrisKind, todoWidth)
	}
	return GameBoard{
		Width: todoWidth,
		Height: todoHeight,
		Set: set,
	}
}

func (gb *GameBoard) draw(buff draw.Image, w, h int) {
	const cellBuffer = 1 
	w -= buffer*2 
	h -= buffer*2 
	drawRect(
		buff.(*image.RGBA), 
		intgeom.Point2{boardX,boardY},
		intgeom.Point2{
			w,
			h,
		},
		color.RGBA{255,255,255,255})
	cellW := w / int(gb.Width)
	cellH := h / int(gb.Height)
	for x := 0; x < int(gb.Width); x++ {
		for y := 0; y < int(gb.Height); y++ {
			if gb.Set[y][x] == KindNone {
				continue 
			}
			// todo: make optimized filled rect draw helper
			c := gb.Set[y][x].Color()
			drawFilledRect(buff.(*image.RGBA), 
				intgeom.Point2{
					boardX + (x*cellW) + cellBuffer,
					boardY + (y*cellH) + cellBuffer,
				},
				intgeom.Point2{
					cellW-cellBuffer*2,
					cellH-cellBuffer*2,
				}, c,
			)
		}
	}
	activeOff := gb.ActiveTris.Offsets()
	activeC := gb.ActiveTris.Color()
	for _, off := range activeOff {
		x := int(gb.ActiveTris.X) + int(off[0])
		y := int(gb.ActiveTris.Y) + int(off[1]) 
		if x >= 0 && x < int(gb.Width) {
			if y >= 0 && y < int(gb.Height) {
				drawFilledRect(buff.(*image.RGBA), 
					intgeom.Point2{
						boardX + (x*cellW) + cellBuffer,
						boardY + (y*cellH) + cellBuffer,
					},
					intgeom.Point2{
						cellW-cellBuffer*2,
						cellH-cellBuffer*2,
					}, activeC,
				)
			}
		}
	}
}

type BoardDimension uint8 

func drawRect(buff *image.RGBA, pos, dims intgeom.Point2, c color.RGBA) {
	render.DrawLine(buff, 
		pos.X(), pos.Y(), 
		pos.X()+dims.X(), pos.Y(), 
		c)
	render.DrawLine(buff, 
		pos.X()+dims.X(), pos.Y(),
		pos.X()+dims.X(), pos.Y()+dims.Y(), 
		c)
	render.DrawLine(buff, 
		pos.X()+dims.X(), pos.Y()+dims.Y(), 
		pos.X(), pos.Y()+dims.Y(), 
		c)
	render.DrawLine(buff, 
		pos.X(), pos.Y()+dims.Y(), 
		pos.X(), pos.Y(), 
		c)
}

func drawFilledRect(buff *image.RGBA, pos, dims intgeom.Point2, c color.RGBA) {
	draw.Draw(buff, image.Rect(pos.X(),pos.Y(),pos.X()+dims.X(),pos.Y()+dims.Y()),
		image.NewUniform(c), image.Point{}, draw.Over)
}

func populateTestBoard(gb GameBoard) {
	for x := 0; x < int(gb.Width); x++ {
		for y := 0; y < int(gb.Height); y++ {
			gb.Set[y][x] = RandomKind()
		}
	}
}

func RandomKind() TrisKind {
	return TrisKind(rand.Intn(int(KindFinal-1)) + 1)
}

func (gb *GameBoard) CheckIfTileIsPlaced() (placed bool) {
	// is there a set tile beneath any tile of the current active tile,
	// it is placed. Do not move it, change the active tile.
	// TODO: game over state 
	// TODO: line clears 
	activeOff := gb.ActiveTris.Offsets()
	for _, off := range activeOff {
		x := int(gb.ActiveTris.X) + int(off[0])
		y := int(gb.ActiveTris.Y) + int(off[1]) 
		if x >= 0 && x < int(gb.Width) {
			if y >= 0 && y < int(gb.Height) {
				if y == int(gb.Height)-1 {
					return true 
				}
				if gb.IsSet(x,y+1) {
					return true 
				}
			}
		}
	}
	return false 
}

func (gb *GameBoard) SetActiveTile() {
	activeOff := gb.ActiveTris.Offsets()
	allY := map[int]struct{}{}
	for _, off := range activeOff {
		x := int(gb.ActiveTris.X) + int(off[0])
		y := int(gb.ActiveTris.Y) + int(off[1]) 
		if x >= 0 && x < int(gb.Width){
			if y >= 0 && y < int(gb.Height) {
				gb.Set[y][x] = gb.ActiveTris.TrisKind	
				allY[y] = struct{}{}			
			}
		}
	} 
	orderedY := []int{} 
	for y := range allY {
		orderedY = append(orderedY, y)
	}
	sort.Slice(orderedY, func(i, j int) bool {
		return i < j 
	})
	for _, y := range orderedY {
		gb.ClearFullLines(y)
	}
}

func (gb *GameBoard) ClearFullLines(y int) {
	for x := 0; x < int(gb.Width); x++ {
		if gb.Set[y][x] == KindNone {
			return 
		}
	}
	gb.Set = append(gb.Set[:y], gb.Set[y+1:]...)
	firstRow := make([]TrisKind, gb.Width)
	gb.Set = append([][]TrisKind{firstRow}, gb.Set...)
}

func (gb *GameBoard) IsSet(x, y int) bool {
	if x < 0 {
		return false 
	}
	if y < 0 {
		return false // ??
	}
	if x >= int(gb.Width) {
		return false 
	}
	if y >= int(gb.Height) {
		return false 
	}
	return gb.Set[y][x] != KindNone
}

func recordGif(ctx *scene.Context) {
	stop := ctx.Window.(*oak.Window).RecordGIF(2)

	go ctx.DoAfter(20 * time.Second, func() {
		g := stop()
		f, err := os.Create("demo.gif")
		if err == nil {
			err = gif.EncodeAll(f, g)
			if err != nil {
				fmt.Println(err)
			}
			f.Close()
		} else {
			fmt.Println(err)
		}
	})
}