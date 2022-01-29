package viertris

import (
	"fmt"
	"image/color"
	"math/rand"
)

type ActiveTris struct {
	Board *GameBoard
	Rotation
	X BoardDimension
	Y BoardDimension
	TrisKind
}

// Offsets on an ActiveTris takes into account rotation
func (at *ActiveTris) Offsets() [4][2]int8 {
	return at.TestOffsets(at.Rotation)
}

func (at *ActiveTris) TestOffsets(rotation Rotation) [4][2]int8 {
	rawOff := at.TrisKind.Offsets()
	switch rotation {
	case Rotation270:
		for i, off := range rawOff {
			off[0], off[1] = off[1], -1*off[0]
			rawOff[i] = off
		}
		fallthrough
	case Rotation180:
		for i, off := range rawOff {
			off[0], off[1] = off[1], -1*off[0]
			rawOff[i] = off
		}
		fallthrough
	case Rotation90:
		for i, off := range rawOff {
			off[0], off[1] = off[1], -1*off[0]
			rawOff[i] = off
		}
		fallthrough
	case Rotation0:
		return rawOff
	default:
		panic(fmt.Sprintf("invalid rotation %v", at.Rotation))
	}
}

func (at *ActiveTris) RotateLeft() {
	newRotation := at.Rotation.RotateLeft()
	off := at.TestOffsets(newRotation)
	for _, o := range off {
		x := int(at.X) + int(o[0])
		y := int(at.Y) + int(o[1])
		if at.Board.IsSet(x, y) || at.Board.IsOffscreen(x, y) {
			return
		}
	}
	at.Rotation = newRotation
}

func (at *ActiveTris) RotateRight() {
	newRotation := at.Rotation.RotateRight()
	off := at.TestOffsets(newRotation)
	for _, o := range off {
		x := int(at.X) + int(o[0])
		y := int(at.Y) + int(o[1])
		if at.Board.IsSet(x, y) || at.Board.IsOffscreen(x, y) {
			return
		}
	}
	at.Rotation = newRotation
}

func (at *ActiveTris) MoveLeft() {
	minX := int(at.X)
	off := at.Offsets()
	for _, o := range off {
		x := int(at.X) + int(o[0])
		y := int(at.Y) + int(o[1])
		if at.Board.IsSet(x-1, y) {
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
		if at.Board.IsSet(x+1, y) {
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
	KindNone   TrisKind = iota
	KindT      TrisKind = iota
	KindLine   TrisKind = iota
	KindSquare TrisKind = iota
	KindZ      TrisKind = iota
	KindS      TrisKind = iota
	KindL      TrisKind = iota
	KindJ      TrisKind = iota
	KindFinal  TrisKind = iota
)

var kindColors = []color.RGBA{
	KindNone:   {},
	KindT:      {200, 0, 0, 255},
	KindLine:   {0, 200, 0, 255},
	KindSquare: {0, 0, 200, 255},
	KindZ:      {200, 200, 0, 255},
	KindS:      {200, 200, 200, 255},
	KindL:      {200, 0, 200, 255},
	KindJ:      {0, 200, 200, 255},
}

func (tk TrisKind) Color() color.RGBA {
	return kindColors[tk]
}

var kindOffsets = [][4][2]int8{
	KindT: {
		{0, 0},
		{-1, 0},
		{0, -1},
		{1, 0},
	},
	KindLine: {
		{0, 0},
		{0, -1},
		{0, 1},
		{0, 2},
	},
	KindSquare: {
		{0, 0},
		{0, 1},
		{1, 1},
		{1, 0},
	},
	KindS: {
		{0, 0},
		{1, 0},
		{0, 1},
		{-1, 1},
	},
	KindZ: {
		{0, 0},
		{-1, 0},
		{0, 1},
		{1, 1},
	},
	KindL: {
		{0, 0},
		{0, -1},
		{0, 1},
		{1, 1},
	},
	KindJ: {
		{0, 0},
		{0, -1},
		{0, 1},
		{-1, 1},
	},
}

func (tk TrisKind) Offsets() [4][2]int8 {
	return kindOffsets[tk]
}

func RandomKind() TrisKind {
	return TrisKind(rand.Intn(int(KindFinal-1)) + 1)
}
