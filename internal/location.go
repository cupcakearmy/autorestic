package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/robfig/cron"
)

type LocationType string

const (
	TypeLocal    LocationType = "local"
	TypeVolume   LocationType = "volume"
	VolumePrefix string       = "volume:"
)

type HookArray = []string

type Hooks struct {
	Before HookArray `yaml:"before"`
	After  HookArray `yaml:"after"`
}

type Options map[string]map[string][]string

type Location struct {
	name    string   `yaml:",omitempty"`
	From    string   `yaml:"from,omitempty"`
	To      []string `yaml:"to,omitempty"`
	Hooks   Hooks    `yaml:"hooks,omitempty"`
	Cron    string   `yaml:"cron,omitempty"`
	Options Options  `yaml:"options,omitempty"`
}

func GetLocation(name string) (Location, bool) {
	l, ok := GetConfig().Locations[name]
	l.name = name
	return l, ok
}

func (l Location) validate(c *Config) error {
	if l.From == "" {
		return fmt.Errorf(`Location "%s" is missing "from" key`, l.name)
	}
	if l.getType() == TypeLocal {
		if from, err := GetPathRelativeToConfig(l.From); err != nil {
			return err
		} else {
			if stat, err := os.Stat(from); err != nil {
				return err
			} else {
				if !stat.IsDir() {
					return fmt.Errorf("\"%s\" is not valid directory for location \"%s\"", from, l.name)
				}
			}
		}
	}

	if len(l.To) == 0 {
		return fmt.Errorf(`Location "%s" has no "to" targets`, l.name)
	}
	// Check if backends are all valid
	for _, to := range l.To {
		_, ok := GetBackend(to)
		if !ok {
			return fmt.Errorf("invalid backend `%s`", to)
		}
	}
	return nil
}

func (l Location) getOptions(key string) []string {
	var options []string
	saved := l.Options[key]
	for k, values := range saved {
		for _, value := range values {
			options = append(options, fmt.Sprintf("--%s", k), value)
		}
	}
	return options
}

func ExecuteHooks(commands []string, options ExecuteOptions) error {
	if len(commands) == 0 {
		return nil
	}
	colors.Secondary.Println("\nRunning hooks")
	for _, command := range commands {
		colors.Body.Println("> " + command)
		out, err := ExecuteCommand(options, "-c", command)
		if err != nil {
			colors.Error.Println(out)
			return err
		}
		if VERBOSE {
			colors.Faint.Println(out)
		}
	}
	colors.Body.Println("")
	return nil
}

func (l Location) getType() LocationType {
	if strings.HasPrefix(l.From, VolumePrefix) {
		return TypeVolume
	}
	return TypeLocal
}

func (l Location) getVolumeName() string {
	return strings.TrimPrefix(l.From, VolumePrefix)
}

func (l Location) getPath() (string, error) {
	t := l.getType()
	switch t {
	case TypeLocal:
		if path, err := GetPathRelativeToConfig(l.From); err != nil {
			return "", err
		} else {
			return path, nil
		}
	case TypeVolume:
		return "/volume/" + l.name + "/" + l.getVolumeName(), nil
	}
	return "", fmt.Errorf("could not get path for location \"%s\"", l.name)
}

func (l Location) Backup(cron bool) error {
	colors.PrimaryPrint("  Backing up location \"%s\"  ", l.name)
	t := l.getType()
	options := ExecuteOptions{
		Command: "bash",
	}

	if t == TypeLocal {
		dir, _ := GetPathRelativeToConfig(l.From)
		options.Dir = dir
	}

	// Hooks
	if err := ExecuteHooks(l.Hooks.Before, options); err != nil {
		return err
	}

	// Backup
	for _, to := range l.To {
		backend, _ := GetBackend(to)
		colors.Secondary.Printf("Backend: %s\n", backend.name)
		env, err := backend.getEnv()
		if err != nil {
			return nil
		}

		flags := l.getOptions("backup")
		cmd := []string{"backup"}
		cmd = append(cmd, flags...)
		if cron {
			cmd = append(cmd, "--tag", "cron")
		}
		cmd = append(cmd, ".")
		backupOptions := ExecuteOptions{
			Dir:  options.Dir,
			Envs: env,
		}

		var out string

		switch t {
		case TypeLocal:
			out, err = ExecuteResticCommand(backupOptions, cmd...)
		case TypeVolume:
			out, err = backend.ExecDocker(l, cmd)
		}
		if err != nil {
			colors.Error.Println(out)
			return err
		}
		if VERBOSE {
			colors.Faint.Println(out)
		}
	}

	// After hooks
	if err := ExecuteHooks(l.Hooks.After, options); err != nil {
		return err
	}
	colors.Success.Println("Done")
	return nil
}

func (l Location) Forget(prune bool, dry bool) error {
	colors.PrimaryPrint("Forgetting for location \"%s\"", l.name)

	path, err := l.getPath()
	if err != nil {
		return err
	}

	for _, to := range l.To {
		backend, _ := GetBackend(to)
		colors.Secondary.Printf("For backend \"%s\"\n", backend.name)
		env, err := backend.getEnv()
		if err != nil {
			return nil
		}
		options := ExecuteOptions{
			Envs: env,
		}
		flags := l.getOptions("forget")
		cmd := []string{"forget", "--path", path}
		if prune {
			cmd = append(cmd, "--prune")
		}
		if dry {
			cmd = append(cmd, "--dry-run")
		}
		cmd = append(cmd, flags...)
		out, err := ExecuteResticCommand(options, cmd...)
		if VERBOSE {
			colors.Faint.Println(out)
		}
		if err != nil {
			return err
		}
	}
	colors.Success.Println("Done")
	return nil
}

func (l Location) hasBackend(backend string) bool {
	for _, b := range l.To {
		if b == backend {
			return true
		}
	}
	return false
}

func (l Location) Restore(to, from string, force bool) error {
	if from == "" {
		from = l.To[0]
	} else if !l.hasBackend(from) {
		return fmt.Errorf("invalid backend: \"%s\"", from)
	}

	to, err := filepath.Abs(to)
	if err != nil {
		return err
	}
	colors.PrimaryPrint("Restoring location \"%s\"", l.name)

	backend, _ := GetBackend(from)
	path, err := l.getPath()
	if err != nil {
		return nil
	}
	colors.Secondary.Println("Restoring lastest snapshot")
	colors.Body.Printf("%s → %s.\n", from, path)
	switch l.getType() {
	case TypeLocal:
		// Check if target is empty
		if !force {
			notEmptyError := fmt.Errorf("target %s is not empty", to)
			_, err = os.Stat(to)
			if err == nil {
				files, err := ioutil.ReadDir(to)
				if err != nil {
					return err
				}
				if len(files) > 0 {
					return notEmptyError
				}
			} else {
				if !os.IsNotExist(err) {
					return err
				}
			}
		}
		err = backend.Exec([]string{"restore", "--target", to, "--path", path, "latest"})
	case TypeVolume:
		_, err = backend.ExecDocker(l, []string{"restore", "--target", ".", "--path", path, "latest"})
	}
	if err != nil {
		return err
	}
	colors.Success.Println("Done")
	return nil
}

func (l Location) RunCron() error {
	if l.Cron == "" {
		return nil
	}

	schedule, err := cron.ParseStandard(l.Cron)
	if err != nil {
		return err
	}
	last := time.Unix(lock.GetCron(l.name), 0)
	next := schedule.Next(last)
	now := time.Now()
	if now.After(next) {
		lock.SetCron(l.name, now.Unix())
		l.Backup(true)
	} else {
		if !CRON_LEAN {
			colors.Body.Printf("Skipping \"%s\", not due yet.\n", l.name)
		}
	}
	return nil
}
