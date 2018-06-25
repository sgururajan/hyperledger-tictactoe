package network

// IdentityConfig - IdentityConfig
type IdentityConfig struct {
	clientConfig ClientConfiguration
	fabClientConfig *FabClientConfig
}

func NewIdentityConfig(cConfig ClientConfiguration) *IdentityConfig {
	return &IdentityConfig{
		clientConfig:cConfig,
		fabClientConfig:NewFabClientConfig(cConfig),
	}
}

