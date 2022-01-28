package main

// Build cross-compiles packages on set of
// OS and architecture pairs.

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	// Defaults
	osxPairs = [][2]string{
		// I think this has to be actually run on osx
		{"darwin", "amd64"},
		{"darwin", "arm64"},
	}
	linuxPairs = [][2]string{
		{"linux", "386"},
		{"linux", "amd64"},
		{"linux", "arm"},
		{"linux", "arm64"},
	}
	winPairs = [][2]string{
		{"windows", "386"},
		{"windows", "amd64"},
		{"windows", "arm"},
		{"windows", "arm64"},
	}
	jsPairs = [][2]string{
		{"js", "wasm"},
	}
	// End Defaults
	android = [][2]string{
		{"android", "arm"},
		{"android", "arm64"},
	}

	osArchPairs [][2]string

	archPairLDFlags = map[[2]string]string{
		{"windows", "386"}:   "-H=windowsgui",
		{"windows", "amd64"}: "-H=windowsgui",
	}

	outputName string
	verbose    bool
	useosx     bool
	usewin     bool
	uselinux   bool
	usedroid   bool
	useall     bool
	usejs      bool
	help       bool
)

func init() {
	flag.BoolVar(&verbose, "v", true, "print build commands as they are run")
	flag.StringVar(&outputName, "o", "viertris", "output executable name")
	flag.BoolVar(&useosx, "osx", true, "build darwin executables")
	flag.BoolVar(&uselinux, "nix", true, "build linux exectuables")
	flag.BoolVar(&usewin, "win", true, "build windows exectuables")
	flag.BoolVar(&usedroid, "android", false, "build android executables")
	flag.BoolVar(&usejs, "js", true, "build js executables")
	flag.BoolVar(&useall, "all", false, "build all executables")
	flag.BoolVar(&help, "h", false, "prints usage")
}

const (
	packageName  = "github.com/200sc/viertris"
	buildinfoPkg = packageName + "/internal/buildinfo"
)

func main() {
	if help {
		fmt.Println("Usage: go run build.go <flags> -pkg <package>")
		return
	}
	flag.Parse()
	if useall {
		useosx = true
		usewin = true
		usedroid = true
		usewin = true
		usejs = true
	}

	lastCommit, err := exec.Command("git", "rev-list", "-1", "HEAD").CombinedOutput()
	if err != nil {
		fmt.Println(string(lastCommit))
		panic(err)
	}

	buildVersion, err := os.ReadFile(filepath.Join("..", "version.json"))
	if err != nil {
		panic(err)
	}

	var buildInfoLDFlags = []string{
		"-X " + buildinfoPkg + ".CheatsEnabled=false",
		"-X " + buildinfoPkg + ".BuildTime=" + time.Now().Format(time.RFC3339),
		"-X " + buildinfoPkg + ".BuildCommit=" + string(lastCommit),
		"-X " + buildinfoPkg + ".BuildVersion=" + string(buildVersion),
	}

	if useosx {
		osArchPairs = append(osArchPairs, osxPairs...)
	}
	if uselinux {
		osArchPairs = append(osArchPairs, linuxPairs...)
	}
	if usedroid {
		osArchPairs = append(osArchPairs, android...)
	}
	if usewin {
		osArchPairs = append(osArchPairs, winPairs...)
	}
	if usejs {
		osArchPairs = append(osArchPairs, jsPairs...)
	}

	workDir := os.TempDir()

	var eg errgroup.Group

	for _, pair := range osArchPairs {
		pair := pair
		eg.Go(func() error {
			buildDir := pair[0] + "-" + pair[1]
			os.MkdirAll(buildDir, 0777)
			buildName := filepath.Join(buildDir, outputName+"."+pair[0]+"."+pair[1])
			if pair[0] == "windows" {
				buildName += ".exe"
			}
			if pair[1] == "wasm" {
				buildName += ".wasm"
			}
			toRun := []string{"build", "-o", buildName}
			ldFlags := "-ldflags="
			allLDFlags := buildInfoLDFlags
			if flags, ok := archPairLDFlags[pair]; ok {
				allLDFlags = append(allLDFlags, flags)
			}
			ldFlags += strings.Join(allLDFlags, " ")
			toRun = append(toRun, packageName)

			env := os.Environ()
			env = append(env, []string{
				"GOOS=" + pair[0],
				"GOARCH=" + pair[1],
				"GOTMPDIR=" + workDir,
			}...)
			if verbose {
				fmt.Println("Running: go ", toRun)
			}
			cmd := exec.Command("go", toRun...)
			cmd.Env = env

			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Target", pair, "errored:")
				fmt.Println(err)
			}
			if verbose && len(out) != 0 {
				fmt.Printf("%s\n", string(out))
			}
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
