package schema

type PersonalConfig struct {
	User      UserInfo       `yaml:"user"`
	Overrides map[string]any `yaml:"overrides,omitempty"`
}

type UserInfo struct {
	Name       string `yaml:"name"`
	Email      string `yaml:"email"`
	GitHub     string `yaml:"github,omitempty"`
	SigningKey string `yaml:"signing_key,omitempty"`
}

func (p *PersonalConfig) ToMap() map[string]string {
	m := make(map[string]string)

	if p.User.Name != "" {
		m["user.name"] = p.User.Name
	}
	if p.User.Email != "" {
		m["user.email"] = p.User.Email
	}
	if p.User.GitHub != "" {
		m["user.github"] = p.User.GitHub
	}
	if p.User.SigningKey != "" {
		m["user.signing_key"] = p.User.SigningKey
	}

	return m
}

func (p *PersonalConfig) IsComplete() bool {
	return p.User.Name != "" && p.User.Email != ""
}
