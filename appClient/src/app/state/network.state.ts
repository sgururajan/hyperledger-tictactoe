import { State, Action, StateContext } from "@ngxs/store";
import { Network } from "../models/network.model";

export class NetworksModel {
    networks: Network[]
    currentNetwork?: Network
}

export class GetNetworks {
    static readonly type="[Networks] GetNetwork";
}

@State<NetworksModel> ({
    name: "networks",
    defaults:{
        networks: []
    }
})

export class NetworksState {
    
    @Action(GetNetworks)
    getNetworks(context: StateContext<NetworksModel>, action: GetNetworks) {

    }
}
