package main

import (
	"fmt"
	"os"

	"github.com/200sc/viertris-templated/internal/scenes"
	"github.com/200sc/viertris-templated/internal/scenes/viertris"
	"github.com/oakmound/oak/v3"
)

func main() {
	oak.AddScene(scenes.Viertris, viertris.Scene)
	err := oak.Init(scenes.Viertris, func(c oak.Config) (oak.Config, error) {
		c.Title = "Viertris"
		return c, nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
