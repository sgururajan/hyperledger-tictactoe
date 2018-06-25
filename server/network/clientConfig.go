package network

// ClientConfig - ClientConfig
//type ClientConfig struct {
//	Organization    string
//	Logging         api.LoggingType
//	CryptoConfig    msp.CCType
//	CredentialStore msp.CredentialStoreType
//	TLSCerts        endpoint.MutualTLSConfig
//	TLSKey          []byte
//	TLSCert         []byte
//}
//
//type CAConfig struct {
//	URL        string
//	TLSCACerts endpoint.MutualTLSConfig
//	Registrar  msp.EnrollCredentials
//	CAName     string
//}
//
//func NewClientConfig() *ClientConfig {
//	return &ClientConfig{}
//}
//
//func (m *ClientConfig) Client() *msp.ClientConfig {
//	return &msp.ClientConfig{
//		Organization:    strings.ToLower(m.Organization),
//		Logging:         m.Logging,
//		CryptoConfig:    m.CryptoConfig,
//		CredentialStore: m.CredentialStore,
//		TLSKey:          m.TLSCerts.Client.Key.Bytes(),
//		TLSCert:         m.TLSCerts.Client.Cert.Bytes(),
//	}
//}
//
//func (m *ClientConfig) getCAConfig(networkConfig *fab.NetworkConfig, org string) (*msp.CAConfig, bool) {
//	if len(networkConfig.Organizations[strings.ToLower(org)].CertificateAuthorities) == 0 {
//		return nil, false
//	}
//
//	organization := networkConfig.Organizations[strings.ToLower(org)]
//	certAuthorityName := organization.CertificateAuthorities[0]
//	if certAuthorityName == "" {
//		return nil, false
//	}
//
//	return nil, true
//}
