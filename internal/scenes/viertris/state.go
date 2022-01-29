package viertris

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

var _ render.Renderable = &GameState{}

const fontSize = 20

func NewGameState(ctx *scene.Context, cfg GameConfig) *GameState {
	gs := &GameState{
		LayeredPoint: render.NewLayeredPoint(0, 0, 0),
		ctx:          ctx,
		w:            ctx.Window.Width(),
		h:            ctx.Window.Height(),
		GameBoard:    NewGameBoard(cfg),
	}
	// TODO: font config
	fnt := render.DefaultFont()
	fnt, _ = fnt.RegenerateWith(func(fg render.FontGenerator) render.FontGenerator {
		fg.Size = fontSize
		return fg
	})
	gs.ClearsText = NewUint64Text("Clears: ", fnt, &gs.Clears)
	gs.ScoreText = NewUint64Text("Score:  ", fnt, &gs.Score)
	gs.LevelText = NewUint64Text("Level:  ", fnt, &gs.Level)

	return gs
	// TODO: on ctx window size change, update w and h
}

type GameState struct {
	render.LayeredPoint
	GameBoard
	GameConfig
	ctx  *scene.Context
	w, h int

	Clears uint64
	Score  uint64
	Level  uint64

	ClearsText *render.Text
	LevelText  *render.Text
	ScoreText  *render.Text

	NextTris   TrisKind
	StoredTris TrisKind
}

const boardRatio = 0.7
const buffer = 5
const boardX = buffer
const boardY = buffer

func (gs *GameState) SetTrisActive(kind TrisKind) (gameOver bool) {
	gs.ActiveTris = ActiveTris{
		Board:    &gs.GameBoard,
		X:        gs.Width / 2,
		Y:        0,
		TrisKind: kind,
	}
	activeOff := gs.ActiveTris.Offsets()
	for _, off := range activeOff {
		x := int(gs.GameBoard.ActiveTris.X) + int(off[0])
		y := int(gs.GameBoard.ActiveTris.Y) + int(off[1])
		if gs.GameBoard.IsSet(x, y) {
			return true
		}
	}
	return false
}

func (gs *GameState) Draw(buff draw.Image, _ float64, _ float64) {

	// board outline
	boardW := int(boardRatio * float64(gs.w))
	gs.GameBoard.draw(buff, boardW, gs.h)

	// sidebar
	// preview / next outlines

	const sidebarRatio = 1 - boardRatio
	const tilePreviewRatio = 0.4
	sidebarW := int(sidebarRatio * float64(gs.w))
	previewH := int(tilePreviewRatio * float64(gs.h))
	drawRect(
		buff.(*image.RGBA),
		intgeom.Point2{
			boardX + boardW,
			boardY,
		},
		intgeom.Point2{
			sidebarW - buffer*2,
			previewH/2 - buffer*2,
		},
		color.RGBA{255, 255, 255, 255})

	drawRect(
		buff.(*image.RGBA),
		intgeom.Point2{
			boardX + boardW,
			boardY + previewH/2,
		},
		intgeom.Point2{
			sidebarW - buffer*2,
			previewH/2 - buffer*2,
		},
		color.RGBA{255, 255, 255, 255})

	// stored

	cellW := boardW / int(gs.GameBoard.Width)
	cellH := gs.h / int(gs.GameBoard.Height)

	if gs.StoredTris != KindNone {
		// TODO: magic numbers
		storedX := boardX + boardW + (cellW * 2)
		storedY := boardY + cellH + cellBuffer
		activeOff := gs.StoredTris.Offsets()
		activeC := gs.StoredTris.Color()
		for _, off := range activeOff {
			x := int(off[0])
			y := int(off[1])
			drawFilledRect(buff.(*image.RGBA),
				intgeom.Point2{
					storedX + (x * cellW) + cellBuffer,
					storedY + (y * cellH) + cellBuffer,
				},
				intgeom.Point2{
					cellW - cellBuffer*2,
					cellH - cellBuffer*2,
				}, activeC,
			)
		}
	}

	// next

	if gs.NextTris != KindNone {
		// TODO: magic numbers
		storedX := boardX + boardW + (cellW * 2)
		storedY := boardY + (cellH * 6) + cellBuffer
		activeOff := gs.NextTris.Offsets()
		activeC := gs.NextTris.Color()
		for _, off := range activeOff {
			x := int(off[0])
			y := int(off[1])
			drawFilledRect(buff.(*image.RGBA),
				intgeom.Point2{
					storedX + (x * cellW) + cellBuffer,
					storedY + (y * cellH) + cellBuffer,
				},
				intgeom.Point2{
					cellW - cellBuffer*2,
					cellH - cellBuffer*2,
				}, activeC,
			)
		}
	}

	// score outline

	const scoreRatio = 1 - tilePreviewRatio
	scoreH := int(scoreRatio * float64(gs.h))
	scoreX := boardX + int(boardRatio*float64(gs.w))
	scoreY := boardY + previewH
	drawRect(
		buff.(*image.RGBA),
		intgeom.Point2{
			scoreX,
			scoreY,
		},
		intgeom.Point2{
			sidebarW - buffer*2,
			scoreH - buffer*2,
		},
		color.RGBA{255, 255, 255, 255})

	scoreX += cellBuffer * 8
	scoreY += cellBuffer * 4
	gs.ScoreText.Draw(buff, float64(scoreX), float64(scoreY))

	// TODO: y delta
	scoreY += fontSize + 4
	gs.ClearsText.Draw(buff, float64(scoreX), float64(scoreY))

	scoreY += fontSize + 4
	gs.LevelText.Draw(buff, float64(scoreX), float64(scoreY))

}

func (gs *GameState) GetDims() (int, int) {
	return gs.w, gs.h
}
