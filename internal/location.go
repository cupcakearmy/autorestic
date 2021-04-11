package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

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
	for _, command := range commands {
		out, err := ExecuteCommand(options, "-c", command)
		fmt.Println(out)
		return err
	}
	return nil
}

func (l Location) forEachBackend(fn func(ExecuteOptions) error) error {
	from, err := GetPathRelativeToConfig(l.From)
	if err != nil {
		return err
	}
	for _, to := range l.To {
		backend, _ := GetBackend(to)
		env, err := backend.getEnv()
		if err != nil {
			return nil
		}
		options := ExecuteOptions{
			Command: "bash",
			Envs:    env,
			Dir:     from,
		}
		if err := fn(options); err != nil {
			return err
		}
	}
	return nil
}

func (l Location) Backup() error {
	return l.forEachBackend(func(options ExecuteOptions) error {
		if err := ExecuteHooks(l.Hooks.Before, options); err != nil {
			return nil
		}

		flags := l.getOptions("backup")
		cmd := []string{"backup"}
		cmd = append(cmd, flags...)
		cmd = append(cmd, ".")
		out, err := ExecuteResticCommand(options, cmd...)
		fmt.Println(out)
		if err != nil {
			return err
		}

		if err := ExecuteHooks(l.Hooks.After, options); err != nil {
			return nil
		}
		return nil
	})
}

func (l Location) Forget(prune bool, dry bool) error {
	return l.forEachBackend(func(options ExecuteOptions) error {
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
		fmt.Println(out)
		if err != nil {
			return err
		}
		return nil
	})
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
	fmt.Printf("Restoring location to %s using %s.\n", to, from)

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
	last := lock.GetCron("test")
	fmt.Println(last)
	next := schedule.Next(time.Unix(last, 0))
	fmt.Println(next)
	now := time.Now()
	if now.After(next) {
		fmt.Println("Running")
		lock.SetCron("test", now.Unix())
	} else {
		fmt.Println("Not due yet")
	}
	return nil
}
