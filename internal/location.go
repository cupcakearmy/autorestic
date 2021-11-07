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
	"github.com/cupcakearmy/autorestic/internal/metadata"
	"github.com/robfig/cron"
)

type LocationType string

const (
	TypeLocal  LocationType = "local"
	TypeVolume LocationType = "volume"
)

type HookArray = []string

type Hooks struct {
	Dir     string    `yaml:"dir"`
	Before  HookArray `yaml:"before,omitempty"`
	After   HookArray `yaml:"after,omitempty"`
	Success HookArray `yaml:"success,omitempty"`
	Failure HookArray `yaml:"failure,omitempty"`
}

type Location struct {
	name    string   `yaml:",omitempty"`
	From    []string `yaml:"from,omitempty"`
	Type    string   `yaml:"type,omitempty"`
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

func (l Location) validate() error {
	if len(l.From) == 0 {
		return fmt.Errorf(`Location "%s" is missing "from" key`, l.name)
	}
	t, err := l.getType()
	if err != nil {
		return err
	}
	switch t {
	case TypeLocal:
		for _, path := range l.From {
			if from, err := GetPathRelativeToConfig(path); err != nil {
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
	case TypeVolume:
		if len(l.From) > 1 {
			return fmt.Errorf(`location "%s" has more than one docker volume`, l.name)
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

func (l Location) ExecuteHooks(commands []string, options ExecuteOptions) error {
	if len(commands) == 0 {
		return nil
	}
	if l.Hooks.Dir != "" {
		if dir, err := GetPathRelativeToConfig(l.Hooks.Dir); err != nil {
			return err
		} else {
			options.Dir = dir
		}
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

func (l Location) getType() (LocationType, error) {
	t := strings.ToLower(l.Type)
	if t == "" || t == "local" {
		return TypeLocal, nil
	} else if t == "volume" {
		return TypeVolume, nil
	}
	return "", fmt.Errorf("invalid location type \"%s\"", l.Type)
}

func buildTag(parts ...string) string {
	parts = append([]string{"ar"}, parts...)
	return strings.Join(parts, ":")
}

func (l Location) getLocationTags() string {
	return buildTag("location", l.name)
}

func (l Location) Backup(cron bool, specificBackend string) []error {
	var errors []error
	var backends []string
	colors.PrimaryPrint("  Backing up location \"%s\"  ", l.name)
	t, err := l.getType()
	if err != nil {
		errors = append(errors, err)
		return errors
	}
	cwd, _ := GetPathRelativeToConfig(".")
	options := ExecuteOptions{
		Command: "bash",
		Dir:     cwd,
		Envs: map[string]string{
			"AUTORESTIC_LOCATION": l.name,
		},
	}

	if err := l.validate(); err != nil {
		errors = append(errors, err)
		goto after
	}

	// Hooks
	if err := l.ExecuteHooks(l.Hooks.Before, options); err != nil {
		errors = append(errors, err)
		goto after
	}

	// Backup
	if specificBackend == "" {
		backends = l.To
	} else {
		if l.hasBackend(specificBackend) {
			backends = []string{specificBackend}
		} else {
			errors = append(errors, fmt.Errorf("backup location \"%s\" has no backend \"%s\"", l.name, specificBackend))
			return errors
		}
	}
	for i, to := range backends {
		backend, _ := GetBackend(to)
		colors.Secondary.Printf("Backend: %s\n", backend.name)
		env, err := backend.getEnv()
		if err != nil {
			errors = append(errors, err)
			continue
		}

		cmd := []string{"backup"}
		cmd = append(cmd, combineOptions("backup", l, backend)...)
		if cron {
			cmd = append(cmd, "--tag", buildTag("cron"))
		}
		cmd = append(cmd, "--tag", l.getLocationTags())
		backupOptions := ExecuteOptions{
			Envs: env,
		}

		var out string
		switch t {
		case TypeLocal:
			for _, from := range l.From {
				path, err := GetPathRelativeToConfig(from)
				if err != nil {
					errors = append(errors, err)
					goto after
				}
				cmd = append(cmd, path)
			}
			out, err = ExecuteResticCommand(backupOptions, cmd...)
		case TypeVolume:
			ok := CheckIfVolumeExists(l.From[0])
			if !ok {
				errors = append(errors, fmt.Errorf("volume \"%s\" does not exist", l.From[0]))
				continue
			}
			cmd = append(cmd, "/data")
			out, err = backend.ExecDocker(l, cmd)
		}
		if err != nil {
			errors = append(errors, err)
			continue
		}

		md := metadata.ExtractMetadataFromBackupLog(out)
		mdEnv := metadata.MakeEnvFromMetadata(&md)
		for k, v := range mdEnv {
			options.Envs[k+"_"+fmt.Sprint(i)] = v
			options.Envs[k+"_"+strings.ToUpper(backend.name)] = v
		}
		if VERBOSE {
			colors.Faint.Println(out)
		}
	}

	// After hooks
	if err := l.ExecuteHooks(l.Hooks.After, options); err != nil {
		errors = append(errors, err)
	}

after:
	var commands []string
	if len(errors) > 0 {
		commands = l.Hooks.Failure
	} else {
		commands = l.Hooks.Success
	}
	if err := l.ExecuteHooks(commands, options); err != nil {
		errors = append(errors, err)
	}

	if len(errors) == 0 {
		colors.Success.Println("Done")
	}
	return errors
}

func (l Location) Forget(prune bool, dry bool) error {
	colors.PrimaryPrint("Forgetting for location \"%s\"", l.name)

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
		cmd := []string{"forget", "--tag", l.getLocationTags()}
		if prune {
			cmd = append(cmd, "--prune")
		}
		if dry {
			cmd = append(cmd, "--dry-run")
		}
		cmd = append(cmd, combineOptions("forget", l, backend)...)
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

func (l Location) Restore(to, from string, force bool, snapshot string) error {
	if from == "" {
		from = l.To[0]
	} else if !l.hasBackend(from) {
		return fmt.Errorf("invalid backend: \"%s\"", from)
	}

	to, err := filepath.Abs(to)
	if err != nil {
		return err
	}

	if snapshot == "" {
		snapshot = "latest"
	}

	colors.PrimaryPrint("Restoring location \"%s\"", l.name)
	backend, _ := GetBackend(from)
	colors.Secondary.Printf("Restoring %s@%s → %s\n", snapshot, backend.name, to)

	t, err := l.getType()
	if err != nil {
		return err
	}
	switch t {
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
		err = backend.Exec([]string{"restore", "--target", to, "--tag", l.getLocationTags(), snapshot})
	case TypeVolume:
		_, err = backend.ExecDocker(l, []string{"restore", "--target", "/", "--tag", l.getLocationTags(), snapshot})
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
		l.Backup(true, "")
	} else {
		if !CRON_LEAN {
			colors.Body.Printf("Skipping \"%s\", not due yet.\n", l.name)
		}
	}
	return nil
}
