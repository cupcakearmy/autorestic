package internal

func RunCron() error {
	c := GetConfig()
	for _, l := range c.Locations {
		err := l.RunCron()
		if err != nil {
			return err
		}
	}
	return nil
}
