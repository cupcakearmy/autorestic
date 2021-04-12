package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cupcakearmy/autorestic/internal/colors"
	"github.com/cupcakearmy/autorestic/internal/lock"
	"github.com/robfig/cron"
)

type HookArray = []string

type Hooks struct {
	Before HookArray `mapstructure:"before"`
	After  HookArray `mapstructure:"after"`
}

type Options map[string]map[string][]string

type Location struct {
	Name    string   `mapstructure:"name"`
	From    string   `mapstructure:"from"`
	To      []string `mapstructure:"to"`
	Hooks   Hooks    `mapstructure:"hooks"`
	Cron    string   `mapstructure:"cron"`
	Options Options  `mapstructure:"options"`
}

func GetLocation(name string) (Location, bool) {
	c := GetConfig()
	for _, b := range c.Locations {
		if b.Name == name {
			return b, true
		}
	}
	return Location{}, false
}

func (l Location) validate(c *Config) error {
	if l.Name == "" {
		return fmt.Errorf(`Location is missing name`)
	}
	if l.From == "" {
		return fmt.Errorf(`Location "%s" is missing "from" key`, l.Name)
	}
	if len(l.To) == 0 {
		return fmt.Errorf(`Location "%s" has no "to" targets`, l.Name)
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
		if VERBOSE {
			colors.Faint.Println(out)
		}
		if err != nil {
			return err
		}
	}
	colors.Body.Println("")
	return nil
}

func (l Location) Backup() error {
	colors.PrimaryPrint("  Backing up location \"%s\"  ", l.Name)
	from, err := GetPathRelativeToConfig(l.From)
	if err != nil {
		return err
	}
	options := ExecuteOptions{
		Command: "bash",
		Dir:     from,
	}
	if err := ExecuteHooks(l.Hooks.Before, options); err != nil {
		return nil
	}
	for _, to := range l.To {
		backend, _ := GetBackend(to)
		colors.Secondary.Printf("Backend: %s\n", backend.Name)
		env, err := backend.getEnv()
		if err != nil {
			return nil
		}
		options := ExecuteOptions{
			Command: "restic",
			Dir:     from,
			Envs:    env,
		}
		flags := l.getOptions("backup")
		cmd := []string{"backup"}
		cmd = append(cmd, flags...)
		cmd = append(cmd, ".")
		out, err := ExecuteResticCommand(options, cmd...)
		if VERBOSE {
			colors.Faint.Println(out)
		}
		if err != nil {
			return err
		}
	}
	if err := ExecuteHooks(l.Hooks.After, options); err != nil {
		return nil
	}
	colors.Success.Println("Done")
	return err
}

func (l Location) Forget(prune bool, dry bool) error {
	colors.PrimaryPrint("Forgetting for location \"%s\"", l.Name)

	from, err := GetPathRelativeToConfig(l.From)
	if err != nil {
		return err
	}
	for _, to := range l.To {
		backend, _ := GetBackend(to)
		colors.Secondary.Printf("For backend \"%s\"\n", backend.Name)
		env, err := backend.getEnv()
		if err != nil {
			return nil
		}
		options := ExecuteOptions{
			Command: "bash",
			Envs:    env,
			Dir:     from,
		}
		flags := l.getOptions("forget")
		cmd := []string{"forget", "--path", options.Dir}
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
	colors.PrimaryPrint("Restoring location \"%s\"", l.Name)
	colors.Secondary.Println("Restoring lastest snapshot")
	colors.Body.Printf("%s → %s.\n", from, to)

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

	backend, _ := GetBackend(from)
	resolved, err := GetPathRelativeToConfig(l.From)
	if err != nil {
		return nil
	}
	err = backend.Exec([]string{"restore", "--target", to, "--path", resolved, "latest"})
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
	last := time.Unix(lock.GetCron(l.Name), 0)
	next := schedule.Next(last)
	now := time.Now()
	if now.After(next) {
		lock.SetCron(l.Name, now.Unix())
		l.Backup()
	} else {
		colors.Body.Printf("Skipping \"%s\", not due yet.\n", l.Name)
	}
	return nil
}
