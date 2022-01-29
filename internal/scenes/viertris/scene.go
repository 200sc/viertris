package viertris

import (
	"fmt"
	"time"

	"github.com/200sc/viertris/internal/buildinfo"
	"github.com/200sc/viertris/internal/scenes"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/scene"
)

// TODO: menu scene to start game, pick level, pick settings
var Scene = scene.Scene{
	Start: func(ctx *scene.Context) {
		cfg, ok := ctx.SceneInput.(GameConfig)
		if !ok {
			fmt.Println("no game config passed to scene")
		}

		st := NewGameState(ctx, cfg)
		ctx.DrawStack.Draw(st)
		var keyRepeatDuration = 70 * time.Millisecond
		var fallDuration = 700 * time.Millisecond
		dropAt := time.Now().Add(fallDuration)

		st.SetTrisActive(RandomKind())
		st.NextTris = RandomKind()

		paused := false
		gameOver := false
		canSwap := true

		keyRepeat := time.Now().Add(keyRepeatDuration)
		ctx.EventHandler.GlobalBind(event.Enter, func(_ event.CID, payload interface{}) int {
			if paused || gameOver {
				return 0
			}
			tileDone := false
			//enter := payload.(event.EnterPayload)
			if time.Now().After(dropAt) {
				tileDone = st.ActiveTris.MoveDown()
				dropAt = time.Now().Add(fallDuration)
			} else if time.Now().After(keyRepeat) {
				if ctx.KeyState.IsDown(key.A) {
					st.ActiveTris.MoveLeft()
					keyRepeat = time.Now().Add(keyRepeatDuration)
				} else if ctx.KeyState.IsDown(key.D) {
					st.ActiveTris.MoveRight()
					keyRepeat = time.Now().Add(keyRepeatDuration)
				} else if ctx.KeyState.IsDown(key.S) {
					tileDone = st.ActiveTris.MoveDown()
					dropAt = time.Now().Add(fallDuration)
					keyRepeat = time.Now().Add(keyRepeatDuration / 2)
				} else if ctx.KeyState.IsDown(key.Q) {
					st.ActiveTris.RotateLeft()
					keyRepeat = time.Now().Add(keyRepeatDuration * 2)
				} else if ctx.KeyState.IsDown(key.E) {
					st.ActiveTris.RotateRight()
					keyRepeat = time.Now().Add(keyRepeatDuration * 2)
				} else if canSwap && ctx.KeyState.IsDown(key.Enter) {
					canSwap = false
					if st.StoredTris == KindNone {
						st.StoredTris = st.ActiveTris.TrisKind
						st.SetTrisActive(st.NextTris)
						st.NextTris = RandomKind()
					} else {
						toStore := st.ActiveTris.TrisKind
						st.SetTrisActive(st.StoredTris)
						st.StoredTris = toStore
					}
				}
			}
			if tileDone {
				canSwap = true
				clears := st.GameBoard.PlaceActiveTile()
				collidedAtStart := st.SetTrisActive(st.NextTris)
				if collidedAtStart {
					gameOver = true
					for y, row := range st.Set {
						for x, kind := range row {
							if kind != KindNone {
								st.Set[y][x] = KindFinal
							}
						}
					}
					// restart binding
					ctx.EventHandler.GlobalBind(key.Down+key.R, func(c event.CID, i interface{}) int {
						ctx.Window.GoToScene(scenes.Viertris)
						return 0
					})
				}
				st.NextTris = RandomKind()

				st.Clears += uint64(clears)
				switch clears {
				case 1:
					st.Score += 1
				case 2:
					st.Score += 4
				case 3:
					st.Score += 9
				case 4:
					st.Score += 16
				}
				lastLevel := st.Level
				st.Level = st.Clears / 10
				if lastLevel < st.Level {
					fallDuration = time.Duration(float64(fallDuration) * .85)
				}
			}
			return 0
		})
		var dropAtDelta time.Duration
		ctx.EventHandler.GlobalBind(key.Down+key.P, func(c event.CID, i interface{}) int {
			paused = !paused
			if paused {
				dropAtDelta = time.Until(dropAt)
			} else {
				dropAt = time.Now().Add(dropAtDelta)
			}
			return 0
		})
		ctx.EventHandler.GlobalBind(key.Down+key.Spacebar, func(c event.CID, i interface{}) int {
			var tileDone bool
			for !tileDone {
				tileDone = st.ActiveTris.MoveDown()
			}
			dropAt = time.Now()
			return 0
		})
		if buildinfo.AreCheatsEnabled() {
			ctx.EventHandler.GlobalBind(key.Down+key.L, func(c event.CID, i interface{}) int {
				st.ActiveTris.TrisKind = KindLine
				return 0
			})
			ctx.EventHandler.GlobalBind(key.Down+key.One, func(c event.CID, i interface{}) int {
				st.Clears += 10
				st.Level = st.Clears / 10
				fallDuration = time.Duration(float64(fallDuration) * .85)
				return 0
			})
		}
	},
}
