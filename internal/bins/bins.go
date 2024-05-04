package bins

import (
	"compress/bzip2"
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
	"github.com/cseitz-forks/autorestic/internal"
	"github.com/cseitz-forks/autorestic/internal/colors"
	"github.com/cseitz-forks/autorestic/internal/flags"
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return parsed, err

	}
	json.Unmarshal(body, &parsed)
	return parsed, nil
}

func Uninstall(restic bool) error {
	if err := os.Remove(path.Join(INSTALL_PATH, "autorestic")); err != nil {
		return err
	}
	if restic {
		if err := os.Remove(path.Join(INSTALL_PATH, "restic")); err != nil {
			return err
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

			// Save to tmp file in the same directory as the install directory
			// Linux does not support overwriting the file that is currently being running
			// But it can be delete the old one and a new one moved in its place.
			tmp, err := os.CreateTemp(INSTALL_PATH, "autorestic-")
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

			to := path.Join(INSTALL_PATH, name)
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
	_, _, err := internal.ExecuteCommand(internal.ExecuteOptions{
		Command: flags.RESTIC_BIN,
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
	body, err := dlJSON("https://api.github.com/repos/cseitz-forks/autorestic/releases/latest")
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
