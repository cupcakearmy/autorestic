package bins

import (
	"compress/bzip2"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
)

const INSTALL_PATH = "/usr/local/bin"

type GithubReleaseAsset struct {
	Name string `json:"name"`
	Link string `json:"browser_download_url"`
}
type GithubRelease struct {
	Tag    string               `json:"tag_name"`
	Assets []GithubReleaseAsset `json:"assets"`
}

func dlJSON(url string) (GithubRelease, error) {
	var parsed GithubRelease
	resp, err := http.Get(url)
	if err != nil {
		return parsed, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return parsed, err

	}
	json.Unmarshal(body, &parsed)
	return parsed, nil
}

func Uninstall(restic bool) error {
	if err := os.Remove(path.Join(INSTALL_PATH, "autorestic")); err != nil {
		colors.Error.Println(err)
	}
	if restic {
		if err := os.Remove(path.Join(INSTALL_PATH, "restic")); err != nil {
			colors.Error.Println(err)
		}
	}
	return nil
}

func downloadAndInstallAsset(body GithubRelease, name string) error {
	ending := fmt.Sprintf("_%s_%s.bz2", runtime.GOOS, runtime.GOARCH)
	for _, asset := range body.Assets {
		if strings.HasSuffix(asset.Name, ending) {
			// Download archive
			colors.Faint.Println("Downloading:", asset.Link)
			resp, err := http.Get(asset.Link)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Uncompress
			bz := bzip2.NewReader(resp.Body)

			// Save binary
			file, err := os.Create(path.Join(INSTALL_PATH, name))
			if err != nil {
				return err
			}
			file.Chmod(0755)
			defer file.Close()
			io.Copy(file, bz)

			colors.Success.Printf("Successfully installed '%s' under %s\n", name, INSTALL_PATH)
			return nil
		}
	}
	return errors.New("could not find right binary for your system, please install restic manually. https://bit.ly/2Y1Rzai")
}

func InstallRestic() error {
	installed := internal.CheckIfCommandIsCallable("restic")
	if installed {
		colors.Body.Println("restic already installed")
		return nil
	} else {
		if body, err := dlJSON("https://api.github.com/repos/restic/restic/releases/latest"); err != nil {
			return err
		} else {
			return downloadAndInstallAsset(body, "restic")
		}
	}
}

func upgradeRestic() error {
	out, err := internal.ExecuteCommand(internal.ExecuteOptions{
		Command: "restic",
	}, "self-update")
	colors.Faint.Println(out)
	return err
}

func Upgrade(restic bool) error {
	// Upgrade restic
	if restic {
		InstallRestic()
		upgradeRestic()
	}

	// Upgrade self
	current, err := semver.ParseTolerant(internal.VERSION)
	if err != nil {
		return err
	}
	body, err := dlJSON("https://api.github.com/repos/cupcakearmy/autorestic/releases/latest")
	if err != nil {
		return err
	}
	latest, err := semver.ParseTolerant(body.Tag)
	if err != nil {
		return err
	}
	if current.LT(latest) {
		downloadAndInstallAsset(body, "autorestic")
		colors.Success.Println("Updated autorestic")
	} else {
		colors.Body.Println("Already up to date")
	}
	return nil
}
