// Heavily inspired (copied) by the restic build file
// https://github.com/restic/restic/blob/aa0faa8c7d7800b6ba7b11164fa2d3683f7f78aa/helpers/build-release-binaries/main.go#L225

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
	"strings"
	"github.com/cupcakearmy/autorestic/internal"
)

var DIR, _ = filepath.Abs("./dist")

var targets = map[string][]string{
	"darwin":  {"amd64", "arm64"},
	"freebsd": {"386", "amd64", "arm"},
	"linux":   {"386", "amd64", "arm", "arm64"},
	"netbsd":  {"386", "amd64"},
	"openbsd": {"386", "amd64"},
	"windows": {"386", "amd64"},
}

type buildOptions struct {
	Target, Arch, Version string
}

func build(options buildOptions, wg *sync.WaitGroup) {
	fmt.Printf("Building %s %s\n", options.Target, options.Arch)
	out := fmt.Sprintf("autorestic_%s_%s_%s", options.Version, options.Target, options.Arch)

	// append .exe for Windows
	if (options.Target == "windows") {
		out += ".exe"
	}

	out = path.Join(DIR, out)
	out, _ = filepath.Abs(out)
	fmt.Println(out)

	// Build
	{
		c := exec.Command("go", "build", "-o", out, "./main.go")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Env = append(os.Environ(),
			"CGO_ENABLED=0",
			"GOOS="+options.Target,
			"GOARCH="+options.Arch,
		)
		err := c.Run()
		if err != nil {
			panic(err)
		}
	}

	// Compress
	{
		var c *exec.Cmd
		switch options.Target {
		// use zip for Windows
		case "windows":
			zipFile := strings.TrimSuffix(out, ".exe") + ".zip"
			c = exec.Command("zip", "-j", "-q", "-X", zipFile, out)
		// use bzip2 for everything else
		default:
			c = exec.Command("bzip2", out)
		}

		c.Dir = DIR
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}

func main() {
	os.RemoveAll(DIR)
	v := internal.VERSION
	var wg sync.WaitGroup
	for target, archs := range targets {
		for _, arch := range archs {
			wg.Add(1)
			build(buildOptions{
				Target:  target,
				Arch:    arch,
				Version: v,
			}, &wg)
		}
	}
	wg.Wait()
}
