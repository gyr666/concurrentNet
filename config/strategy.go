package config

type GetFromFileStrategy struct {
}

func (g *GetFromFileStrategy)Fill(c *Config) error{
	return nil
}

type GetFromDefaultStrategy struct {
}

func (g *GetFromDefaultStrategy)Fill(c *Config) error{
	return nil
}