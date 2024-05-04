// Heavily inspired (copied) by the restic build file
// https://github.com/restic/restic/blob/aa0faa8c7d7800b6ba7b11164fa2d3683f7f78aa/helpers/build-release-binaries/main.go#L225

package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/cseitz-forks/autorestic/internal"
)

var DIR, _ = filepath.Abs("./dist")

var targets = map[string][]string{
	// "aix":     {"ppc64"}, // Not supported by fsnotify
	"darwin":  {"amd64", "arm64"},
	"freebsd": {"386", "amd64", "arm"},
	"linux":   {"386", "amd64", "arm", "arm64", "ppc64le", "mips", "mipsle", "mips64", "mips64le", "s390x"},
	"netbsd":  {"386", "amd64"},
	"openbsd": {"386", "amd64"},
	// "windows": {"386", "amd64"}, // Not supported by autorestic
	"solaris": {"amd64"},
}

type buildOptions struct {
	Target, Arch, Version string
}

const (
	CHECKSUM_MD5     = "MD5SUMS"
	CHECKSUM_SHA_1   = "SHA1SUMS"
	CHECKSUM_SHA_256 = "SHA256SUMS"
)

type Checksums struct {
	filename, md5, sha1, sha256 string
}

func writeChecksums(checksums *[]Checksums) {
	FILE_MD5, _ := os.OpenFile(path.Join(DIR, CHECKSUM_MD5), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer FILE_MD5.Close()
	FILE_SHA1, _ := os.OpenFile(path.Join(DIR, CHECKSUM_SHA_1), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer FILE_SHA1.Close()
	FILE_SHA256, _ := os.OpenFile(path.Join(DIR, CHECKSUM_SHA_256), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	defer FILE_SHA256.Close()

	for _, checksum := range *checksums {
		fmt.Fprintf(FILE_MD5, "%s %s\n", checksum.md5, checksum.filename)
		fmt.Fprintf(FILE_SHA1, "%s %s\n", checksum.sha1, checksum.filename)
		fmt.Fprintf(FILE_SHA256, "%s %s\n", checksum.sha256, checksum.filename)
	}
}

func build(options buildOptions, wg *sync.WaitGroup, checksums *[]Checksums) {
	defer wg.Done()

	fmt.Printf("Building: %s %s\n", options.Target, options.Arch)
	out := fmt.Sprintf("autorestic_%s_%s_%s", options.Version, options.Target, options.Arch)
	out = path.Join(DIR, out)

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
		c := exec.Command("bzip2", out)
		c.Dir = DIR
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			panic(err)
		}
	}

	// Checksum
	{
		file := out + ".bz2"
		content, _ := ioutil.ReadFile(file)
		*checksums = append(*checksums, Checksums{
			filename: path.Base(file),
			md5:      fmt.Sprintf("%x", md5.Sum(content)),
			sha1:     fmt.Sprintf("%x", sha1.Sum(content)),
			sha256:   fmt.Sprintf("%x", sha256.Sum256(content)),
		})
	}

	fmt.Printf("Built: %s\n", path.Base(out))
}

func main() {
	os.RemoveAll(DIR)
	v := internal.VERSION
	checksums := []Checksums{}

	// Build async
	var wg sync.WaitGroup
	for target, archs := range targets {
		for _, arch := range archs {
			wg.Add(1)
			go build(buildOptions{
				Target:  target,
				Arch:    arch,
				Version: v,
			}, &wg, &checksums)
		}
	}
	wg.Wait()
	writeChecksums(&checksums)
}
