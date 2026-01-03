package schema

type Profile struct {
	Name        string         `yaml:"name,omitempty"`
	Description string         `yaml:"description,omitempty"`
	Inherits    []string       `yaml:"inherits,omitempty"`
	Configs     []string       `yaml:"configs,omitempty"`
	Packages    []string       `yaml:"packages,omitempty"`
	Settings    map[string]any `yaml:"settings,omitempty"`
	Scripts     ProfileScripts `yaml:"scripts,omitempty"`
}

type ProfileScripts struct {
	PreInstall  []string `yaml:"pre_install,omitempty"`
	PostInstall []string `yaml:"post_install,omitempty"`
	PreUpdate   []string `yaml:"pre_update,omitempty"`
	PostUpdate  []string `yaml:"post_update,omitempty"`
}

func (p *Profile) HasConfig(name string) bool {
	for _, c := range p.Configs {
		if c == name {
			return true
		}
	}
	return false
}

func (p *Profile) HasPackageGroup(name string) bool {
	for _, pkg := range p.Packages {
		if pkg == name {
			return true
		}
	}
	return false
}

func (p *Profile) GetSetting(key string) any {
	if p.Settings == nil {
		return nil
	}
	return p.Settings[key]
}

func (p *Profile) GetStringListSetting(key string) []string {
	val := p.GetSetting(key)
	if val == nil {
		return nil
	}

	if list, ok := val.([]string); ok {
		return list
	}

	if list, ok := val.([]any); ok {
		result := make([]string, 0, len(list))
		for _, item := range list {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}

	return nil
}
