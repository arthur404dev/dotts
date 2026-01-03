package schema

type RepoConfig struct {
	Name            string   `yaml:"name"`
	Author          string   `yaml:"author,omitempty"`
	Description     string   `yaml:"description,omitempty"`
	Version         string   `yaml:"version,omitempty"`
	DefaultMachine  string   `yaml:"default_machine,omitempty"`
	Features        []string `yaml:"features,omitempty"`
	MinDottsVersion string   `yaml:"min_dotts_version,omitempty"`
}

func (c *RepoConfig) HasFeature(feature string) bool {
	for _, f := range c.Features {
		if f == feature {
			return true
		}
	}
	return false
}
