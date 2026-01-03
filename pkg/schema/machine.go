package schema

type Machine struct {
	Machine  MachineInfo    `yaml:"machine"`
	Inherits []string       `yaml:"inherits,omitempty"`
	Settings map[string]any `yaml:"settings,omitempty"`
	Features []string       `yaml:"features,omitempty"`
}

type MachineInfo struct {
	Hostname    string `yaml:"hostname"`
	Description string `yaml:"description,omitempty"`
}

func (m *Machine) GetSetting(key string) any {
	if m.Settings == nil {
		return nil
	}
	return m.Settings[key]
}

func (m *Machine) GetStringSetting(key string) string {
	val := m.GetSetting(key)
	if val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

func (m *Machine) GetIntSetting(key string) int {
	val := m.GetSetting(key)
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}

func (m *Machine) GetBoolSetting(key string) bool {
	val := m.GetSetting(key)
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

func (m *Machine) HasFeature(feature string) bool {
	for _, f := range m.Features {
		if f == feature {
			return true
		}
	}
	return false
}

func (m *Machine) HasInherit(profile string) bool {
	for _, p := range m.Inherits {
		if p == profile {
			return true
		}
	}
	return false
}
