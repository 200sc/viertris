package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/200sc/viertris/internal/scenes"
	"github.com/200sc/viertris/internal/scenes/viertris"
	"github.com/oakmound/oak/v3"
)

func main() {

	rand.Seed(time.Now().Unix())

	oak.AddScene(scenes.Viertris, viertris.Scene)
	err := oak.Init(scenes.Viertris, func(c oak.Config) (oak.Config, error) {
		c.Title = "Viertris"
		c.BatchLoad = false
		return c, nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
