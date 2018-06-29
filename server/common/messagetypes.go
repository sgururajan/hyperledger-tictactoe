package common

type CreateChannelRequest struct {
	ChannelName string `json:"channelName"`
	OrganizationNames []string `json:"organizationNames"`
	ConsortiumName string `json:"consortiumName"`
	AnchorPeers map[string][]string `json:"anchorPeers"` //anchor peers for each organization
}
