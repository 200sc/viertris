package main

import (
	"fmt"
	"os"
	"image"
	"image/draw"
	"image/color"
	"math/rand"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/scene"
	"github.com/oakmound/oak/v3/render"
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
			st := NewGameState(ctx, GameConfig{}) 
			ctx.DrawStack.Draw(st)
			populateTestBoard(st.GameBoard)
	// Game Scene:
	// -- Game Board 
	// -- Score / Level tracking 
	// -- Stored / Preview window 
	// -- Events:
	// --- EnterFrame: move current tris down 1 after duration 
	// --- Inputs:
	// ---- rotate current tris 
	// ---- drop current tris 
	// ---- store current tris 
	// ---- retrieve stored tris 
		},
	})
	// Init 
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
	ThisTris ActiveTris
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

type ActiveTris struct {
	X BoardDimension 
	Y BoardDimension 
	TrisKind
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
	Set [][]TrisKind
}

func NewGameBoard(cfg GameConfig) GameBoard {
	// TODO 
	const todoWidth = 8
	const todoHeight = 24 
	set := make([][]TrisKind, todoWidth)
	for i := range set {
		set[i] = make([]TrisKind, todoHeight)
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
	for x := 0; x < len(gb.Set); x++ {
		for y := 0; y < len(gb.Set[0]); y++ {
			if gb.Set[x][y] == KindNone {
				continue 
			}
			// todo: make optimized filled rect draw helper
			c := gb.Set[x][y].Color()
			drawRect(buff.(*image.RGBA), 
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

func populateTestBoard(gb GameBoard) {
	for x := 0; x < len(gb.Set); x++ {
		for y := 0; y < len(gb.Set[0]); y++ {
			k := rand.Intn(int(KindFinal-1)) + 1 
			gb.Set[x][y] = TrisKind(k)
		}
	}
}