package viertris

import (
	"image"
	"image/color"
	"image/draw"
	"sort"

	"github.com/oakmound/oak/v3/alg/intgeom"
)

type GameBoard struct {
	Width  BoardDimension
	Height BoardDimension
	// Unexpected! Y x X! Because that makes it easier to clear lines!
	Set        [][]TrisKind
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
		Width:  todoWidth,
		Height: todoHeight,
		Set:    set,
	}
}

const cellBuffer = 1

func (gb *GameBoard) draw(buff draw.Image, w, h int) {
	w -= buffer * 2
	h -= buffer * 2
	drawRect(
		buff.(*image.RGBA),
		intgeom.Point2{boardX, boardY},
		intgeom.Point2{
			w,
			h,
		},
		color.RGBA{255, 255, 255, 255})
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
					boardX + (x * cellW) + cellBuffer,
					boardY + (y * cellH) + cellBuffer,
				},
				intgeom.Point2{
					cellW - cellBuffer*2,
					cellH - cellBuffer*2,
				}, c,
			)
		}
	}
	// TODO: combine with drawing of stored / next?
	activeOff := gb.ActiveTris.Offsets()
	activeC := gb.ActiveTris.Color()
	for _, off := range activeOff {
		x := int(gb.ActiveTris.X) + int(off[0])
		y := int(gb.ActiveTris.Y) + int(off[1])
		if x >= 0 && x < int(gb.Width) {
			if y >= 0 && y < int(gb.Height) {
				drawFilledRect(buff.(*image.RGBA),
					intgeom.Point2{
						boardX + (x * cellW) + cellBuffer,
						boardY + (y * cellH) + cellBuffer,
					},
					intgeom.Point2{
						cellW - cellBuffer*2,
						cellH - cellBuffer*2,
					}, activeC,
				)
			}
		}
	}
}

func (gb *GameBoard) CheckIfTileIsPlaced() (placed bool) {
	// is there a set tile beneath any tile of the current active tile,
	// it is placed. Do not move it, change the active tile.
	// TODO: game over state
	activeOff := gb.ActiveTris.Offsets()
	for _, off := range activeOff {
		x := int(gb.ActiveTris.X) + int(off[0])
		y := int(gb.ActiveTris.Y) + int(off[1])
		if x >= 0 && x < int(gb.Width) {
			if y >= 0 && y < int(gb.Height) {
				if y == int(gb.Height)-1 {
					return true
				}
				if gb.IsSet(x, y+1) {
					return true
				}
			}
		}
	}
	return false
}

func (gb *GameBoard) PlaceActiveTile() (clears uint32, gameOver bool) {
	activeOff := gb.ActiveTris.Offsets()
	allY := map[int]struct{}{}
	for _, off := range activeOff {
		x := int(gb.ActiveTris.X) + int(off[0])
		y := int(gb.ActiveTris.Y) + int(off[1])
		if x >= 0 && x < int(gb.Width) {
			if y >= 0 && y < int(gb.Height) {
				gb.Set[y][x] = gb.ActiveTris.TrisKind
				allY[y] = struct{}{}
			} else if y > int(gb.Height) {
				gameOver = true
			}
		}
	}
	orderedY := []int{}
	for y := range allY {
		orderedY = append(orderedY, y)
	}
	sort.Slice(orderedY, func(i, j int) bool {
		return orderedY[i] < orderedY[j]
	})

	clears = 0
	for _, y := range orderedY {
		if gb.ClearFullLines(y) {
			clears++
		}
	}
	return clears, gameOver
}

func (gb *GameBoard) ClearFullLines(y int) (cleared bool) {
	for x := 0; x < int(gb.Width); x++ {
		if gb.Set[y][x] == KindNone {
			return false
		}
	}
	gb.Set = append(gb.Set[:y], gb.Set[y+1:]...)
	firstRow := make([]TrisKind, gb.Width)
	gb.Set = append([][]TrisKind{firstRow}, gb.Set...)
	return true
}

func (gb *GameBoard) IsSet(x, y int) bool {
	if gb.IsOffscreen(x, y) {
		return false
	}
	return gb.Set[y][x] != KindNone
}

func (gb *GameBoard) IsOffscreen(x, y int) bool {
	if x < 0 {
		return true
	}
	if y < 0 {
		return true // ??
	}
	if x >= int(gb.Width) {
		return true
	}
	if y >= int(gb.Height) {
		return true
	}
	return false
}
