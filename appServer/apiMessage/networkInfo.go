package apiMessage

type NetworkInfo struct {
	Name              string `json:"name"`
	NoOfOrganizations int    `json:"noOfOrganizations"`
	NoOfPeers         int    `json:"noOfPeers"`
	NoOfChannels      int    `json:"noOfChannels"`
	NoOfBlocks        int    `json:"noOfBlocks"`
}
