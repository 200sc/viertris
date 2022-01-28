//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build namespace for general build tooling.
type Build mg.Namespace

// All distros should be built.
func (Build) All() error {
	os.Chdir("build")
	return sh.RunV("go", "run", "build.go", "-win", "-js", "-nix", "-osx", "-v")
}

// Js is a convience for the js build run pattern.
func (Build) Js() error {
	os.Chdir("build")
	toRun := []string{"run", "build.go", "-js", "-win=false"}
	fmt.Println("Running: go ", toRun)
	err := sh.RunV("go", toRun...)
	os.Chdir("..")
	return err
}
