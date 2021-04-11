// Heavily inspired (copied) by the restic build file
// https://github.com/restic/restic/blob/aa0faa8c7d7800b6ba7b11164fa2d3683f7f78aa/helpers/build-release-binaries/main.go#L225

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/cupcakearmy/autorestic/internal"
)

var DIR, _ = filepath.Abs("./dist")

var targets = map[string][]string{
	"darwin":  {"amd64"},
	"freebsd": {"386", "amd64", "arm"},
	"linux":   {"386", "amd64", "arm", "arm64"},
	"netbsd":  {"386", "amd64"},
	"openbsd": {"386", "amd64"},
}

type buildOptions struct {
	Target, Arch, Version string
}

func build(options buildOptions) error {
	fmt.Printf("Building %s %s\n", options.Target, options.Arch)
	out := fmt.Sprintf("autorestic_%s_%s_%s", options.Version, options.Target, options.Arch)
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
			return err
		}
	}

	// Compress
	{
		c := exec.Command("bzip2", out)
		c.Dir = DIR
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	os.RemoveAll(DIR)
	v := internal.VERSION
	for target, archs := range targets {
		for _, arch := range archs {
			err := build(buildOptions{
				Target:  target,
				Arch:    arch,
				Version: v,
			})
			if err != nil {
				panic(err)
			}
		}
	}
}
