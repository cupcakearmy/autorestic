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

func InstallRestic() error {
	installed := internal.CheckIfCommandIsCallable("restic")
	if installed {
		fmt.Println("restic already installed")
		return nil
	} else {
		body, err := dlJSON("https://api.github.com/repos/restic/restic/releases/latest")
		if err != nil {
			return err
		}
		ending := fmt.Sprintf("_%s_%s.bz2", runtime.GOOS, runtime.GOARCH)
		for _, asset := range body.Assets {
			if strings.HasSuffix(asset.Name, ending) {
				// Found
				fmt.Println(asset.Link)

				// Download archive
				resp, err := http.Get(asset.Link)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				// Uncompress
				bz := bzip2.NewReader(resp.Body)

				// Save binary
				file, err := os.Create(path.Join(INSTALL_PATH, "restic"))
				if err != nil {
					return err
				}
				file.Chmod(0755)
				defer file.Close()
				io.Copy(file, bz)

				fmt.Printf("Successfully installed restic under %s\n", INSTALL_PATH)
				return nil
			}
		}
		return errors.New("could not find right binary for your system, please install restic manually. https://bit.ly/2Y1Rzai")
	}
}

func upgradeRestic() error {
	out, err := internal.ExecuteCommand(internal.ExecuteOptions{
		Command: "restic",
	}, "self-update")
	fmt.Println(out)
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
	fmt.Println(current)

	body, err := dlJSON("https://api.github.com/repos/cupcakearmy/autorestic/releases/latest")
	if err != nil {
		return err
	}
	latest, err := semver.ParseTolerant(body.Tag)
	if err != nil {
		return err
	}
	if current.GT(latest) {

		fmt.Println("Updated autorestic")
	} else {
		fmt.Println("Already up to date")
	}
	return nil
}
