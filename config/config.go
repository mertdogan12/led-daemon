package config

import "github.com/namsral/flag"

type Config struct{}

func (c *Config) Init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.String(flag.DefaultConfigFlagname, "", "~/.config/led/config.conf")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	return nil
}
