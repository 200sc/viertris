package main

import (
	"fmt"
	"os"
	"image/color"

	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/scene"
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
		End: func() (string, *scene.Result) {
			return GameSceneName, nil
		},
	})
	oak.AddScene(GameSceneName, scene.Scene{
		Start: func(ctx *scene.Context) {
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
		End: func() (string, *scene.Result) {
			return MenuSceneName, nil 
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

type GameState struct {
	GameBoard 
	GameConfig
	Clears uint32
	ThisTris ActiveTris
	Level uint16
	NextTris TrisKind 
	StoredTris TrisKind 
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

type BoardDimension uint8 