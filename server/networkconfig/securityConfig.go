package networkconfig

// FabSecurityConfiguration - SecurityConfig
type FabSecurityConfiguration struct {
	config SecurityConfiguration
}

// newFabSecurityConfig - NewSecurityConfig
func newFabSecurityConfig(secConfig SecurityConfiguration) *FabSecurityConfiguration {
	return &FabSecurityConfiguration{
		config: secConfig,
	}
}

// IsSecurityEnabled - IsSecurityEnabled
func (m *FabSecurityConfiguration) IsSecurityEnabled() bool {
	return m.config.IsEnabled
}

func (m *FabSecurityConfiguration) SecurityAlgorithm() string {
	return m.config.Algorithm
}

func (m *FabSecurityConfiguration) SecurityLevel() int {
	return m.config.Level
}

func (m *FabSecurityConfiguration) SecurityProvider() string {
	return m.config.Provider
}

func (m *FabSecurityConfiguration) SoftVerify() bool {
	return m.config.IsSoftVerify
}

func (m *FabSecurityConfiguration) SecurityProviderLibPath() string {
	return m.config.ProviderLibPath
}

func (m *FabSecurityConfiguration) SecurityProviderPin() string {
	return m.config.ProviderPin
}

func (m *FabSecurityConfiguration) SecurityProviderLabel() string {
	return m.config.ProviderLabel
}

func (m *FabSecurityConfiguration) KeyStorePath() string {
	return m.config.KeyStoragePath
}
