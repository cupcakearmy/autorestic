package internal

func RunCron() error {
	c := GetConfig()
	for _, l := range c.Locations {
		if err := l.RunCron(); err != nil {
			return err
		}
	}
	return nil
}
