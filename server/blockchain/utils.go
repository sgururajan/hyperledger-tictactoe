package blockchain

/*func getDefaultOrdererEndpoint(config *networkconfig.FabNetworkConfiguration) (string, bool) {
	for _,v:= range config.OrderersInfo {
		if v.IsPrimary {
			return v.Endpoint, true
		}
	}

	if len(config.OrderersInfo) > 0 {
		return config.OrderersInfo[0].Endpoint, true
	}

	return "", false
}

func getDefaultOrgInfo(config *networkconfig.FabNetworkConfiguration) (networkconfig.OrgInfo, bool) {
	for _, v:= range config.OrgsInfo {
		if v.IsPrimary {
			return v, true
		}
	}

	if len(config.OrgsInfo) > 0 {
		return config.OrgsInfo[0], true
	}

	return networkconfig.OrgInfo{}, false
}

func getDefaultOrgEndpoint(config *networkconfig.FabNetworkConfiguration) (string, bool) {
	for _, v:= range config.OrgsInfo {
		if v.IsPrimary {
			return v.Endpoint, true
		}
	}

	if len(config.OrgsInfo) > 0 {
		return config.OrgsInfo[0].Endpoint, true
	}

	return "", false
}



func getDefaultEndrosingPeerEndpointForOrganization(config *networkconfig.FabNetworkConfiguration, orgName string) (string, bool) {
	for _,v:= range config.PeersInfo {
		if v.OrgName == orgName && v.IsEndrosingPeer {
			return v.Endpoint, true
		}
	}

	// return the first peer for the organization
	for _,v:= range config.PeersInfo {
		if v.OrgName==orgName {
			return v.Endpoint, true
		}
	}

	return "", false
}*/


