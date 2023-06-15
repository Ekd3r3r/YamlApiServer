package model

type Maintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type Metadata struct {
	Title       string       `yaml:"title"`
	Version     string       `yaml:"version"`
	Maintainers []Maintainer `yaml:"maintainers"`
	Company     string       `yaml:"company"`
	Website     string       `yaml:"website"`
	Source      string       `yaml:"source"`
	License     string       `yaml:"license"`
	Description string       `yaml:"description"`
}
