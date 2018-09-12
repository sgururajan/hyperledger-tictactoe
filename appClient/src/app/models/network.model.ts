export interface Network extends NamedObject {
    organizations: Organization[]
    channels: Channel[]
}

export interface Organization extends NamedObject {
    peers: Peer[]
}

export interface Peer extends NamedObject {
    url: string
}

export interface Channel extends NamedObject {
}

export interface NamedObject {
    name: string
}