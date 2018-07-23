package viewModel

type PeerViewModel struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type ChannelViewModel struct {
	Name string `json:"name"`
}

type OrganizationViewModel struct {
	Name string `json:"name"`
	Peers []PeerViewModel `json:"peers"`
}

type NetworkViewModel struct {
	Name string `json:"name"`
	Organizations []OrganizationViewModel `json:"organizations"`
	Channels []ChannelViewModel `json:"channels"`
}
