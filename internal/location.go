package internal

import (
	"fmt"
)

type HookArray = []string

type Hooks struct {
	Before HookArray `mapstructure:"before"`
	After  HookArray `mapstructure:"after"`
}

type Options map[string]map[string][]string

type Location struct {
	From    string   `mapstructure:"from"`
	To      []string `mapstructure:"to"`
	Hooks   Hooks    `mapstructure:"hooks"`
	Cron    string   `mapstructure:"cron"`
	Options Options  `mapstructure:"options"`
}

func (l Location) validate(c Config) error {
	// Check if backends are all valid
	for _, to := range l.To {
		_, ok := c.Backends[to]
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

func (l Location) Backup() error {
	c := GetConfig()
	from := GetPathRelativeToConfig(l.From)
	for _, to := range l.To {
		backend := c.Backends[to]
		options := ExecuteOptions{
			Command: "bash",
			Envs:    backend.getEnv(),
			Dir:     from,
		}

		if err := ExecuteHooks(l.Hooks.Before, options); err != nil {
			return nil
		}
		{
			flags := l.getOptions("backup")
			cmd := []string{"backup"}
			cmd = append(cmd, flags...)
			cmd = append(cmd, ".")
			out, err := ExecuteResticCommand(options, cmd...)
			fmt.Println(out)
			if err != nil {
				return err
			}
		}
		if err := ExecuteHooks(l.Hooks.After, options); err != nil {
			return nil
		}
	}
	return nil
}
