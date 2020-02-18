package config

type GetFromFileStrategy struct {
}

func (g *GetFromFileStrategy) Fill(c *Config) error {
	return nil
}
