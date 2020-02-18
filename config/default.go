package config

type GetFromDefaultStrategy struct {
}

func (g *GetFromDefaultStrategy) Fill(c *Config) error {
	return nil
}
