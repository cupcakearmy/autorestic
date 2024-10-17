package internal

import (
	"errors"
	"fmt"
)

func RunCron() error {
	c := GetConfig()
	var errs []error
	for name, l := range c.Locations {
		l.name = name
		if err := l.RunCron(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("Encountered errors during cron process:\n%w", errors.Join(errs...))
	}
	return nil
}
