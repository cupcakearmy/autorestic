package bins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/cupcakearmy/autorestic/internal"
	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/flags"
)

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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return parsed, err

	}
	json.Unmarshal(body, &parsed)
	return parsed, nil
}

func Uninstall(restic bool) error {
	if err := os.Remove(path.Join(installPath(), exeName("autorestic"))); err != nil {
		return err
	}
	if restic {
		if err := os.Remove(path.Join(installPath(), exeName("restic"))); err != nil {
			return err
		}
	}
	return nil
}

func downloadAndInstallAsset(body GithubRelease, name string) error {
	ending := fmt.Sprintf("_%s_%s.%s", runtime.GOOS, runtime.GOARCH, formatName)
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
			bz, err := decompress(resp)
			if err != nil {
				return err
			}

			// Save to tmp file in the same directory as the install directory
			// Linux does not support overwriting the file that is currently running
			// But it can delete the old one and a new one can be moved in its place.
			tmp, err := os.CreateTemp(installPath(), "autorestic-")
			if err != nil {
				return err
			}
			defer tmp.Close()
			if err := tmp.Chmod(0755); err != nil {
				return err
			}
			if _, err := io.Copy(tmp, bz); err != nil {
				return err
			}
			// On Windows the rename below will fail if the file remains open.
			tmp.Close()

			to := path.Join(installPath(), exeName(name))
			defer os.Remove(tmp.Name()) // Cleanup temporary file after thread exits

			mode := os.FileMode(0755)
			if originalBin, err := os.Lstat(to); err == nil {
				mode = originalBin.Mode()
				err := os.Remove(to)
				if err != nil {
					return err
				}
			}

			err = os.Rename(tmp.Name(), to)
			if err != nil {
				return err
			}

			err = os.Chmod(to, mode)
			if err != nil {
				return err
			}

			colors.Success.Printf("Successfully installed '%s' under %s\n", name, installPath())
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
	_, _, err := internal.ExecuteCommand(internal.ExecuteOptions{
		Command: exeName(flags.RESTIC_BIN),
	}, "self-update")
	return err
}

func Upgrade(restic bool) error {
	// Upgrade restic
	if restic {
		if err := InstallRestic(); err != nil {
			return err
		}
		if err := upgradeRestic(); err != nil {
			return err
		}
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
		if err := downloadAndInstallAsset(body, "autorestic"); err != nil {
			return err
		}
		colors.Success.Println("Updated autorestic")
	} else {
		colors.Body.Println("Already up to date")
	}
	return nil
}
