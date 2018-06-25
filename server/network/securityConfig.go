package network

// SecurityConfig - SecurityConfig
type SecurityConfig struct {
	config SecurityConfiguration
}

// NewSecurityConfig - NewSecurityConfig
func NewSecurityConfig(secConfig SecurityConfiguration) *SecurityConfig {
	return &SecurityConfig{
		config: secConfig,
	}
}

func (m *SecurityConfig) IsSecurityEnabled() bool {
	return m.config.IsEnabled
}

func (m *SecurityConfig) SecurityAlgorithm() string {
	return m.config.Algorithm
}

func (m *SecurityConfig) SecurityLevel() int {
	return m.config.Level
}

func (m *SecurityConfig) SecurityProvider() string {
	return m.config.Provider
}

func (m *SecurityConfig) SoftVerify() bool {
	return m.config.IsSoftVerify
}

func (m *SecurityConfig) SecurityProviderLibPath() string {
	return m.config.ProviderLibPath
}

func (m *SecurityConfig) SecurityProviderPin() string {
	return m.config.ProviderPin
}

func (m *SecurityConfig) SecurityProviderLabel() string {
	return m.config.ProviderLabel
}

func (m *SecurityConfig) KeyStorePath() string {
	return m.config.KeyStoragePath
}
