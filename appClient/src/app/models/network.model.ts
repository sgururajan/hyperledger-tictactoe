export interface NetworksInfo {
    name: string
    noOfOrganizations: number
    noOfPeers: number
    noOfChannels: number
    noOfBlocks: number
}

export interface Network {
    name: string
    organizations: string[]
    channels: string[]
}