package internal

func RunCron() error {
	c := GetConfig()
	for name, l := range c.Locations {
		l.name = name
		if err := l.RunCron(); err != nil {
			return err
		}
	}
	return nil
}
