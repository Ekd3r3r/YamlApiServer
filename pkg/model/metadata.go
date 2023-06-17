package model

type Maintainer struct {
	Name  string `yaml:"name" valid:"ascii"`
	Email string `yaml:"email" valid:"email"`
}

type Metadata struct {
	Title       string       `yaml:"title" valid:"ascii"`
	Version     string       `yaml:"version" valid:"semver"`
	Maintainers []Maintainer `yaml:"maintainers" valid:"-"`
	Company     string       `yaml:"company" valid:"ascii"`
	Website     string       `yaml:"website" valid:"url"`
	Source      string       `yaml:"source" valid:"url"`
	License     string       `yaml:"license" valid:"ascii"`
	Description string       `yaml:"description" valid:"ascii"`
}
