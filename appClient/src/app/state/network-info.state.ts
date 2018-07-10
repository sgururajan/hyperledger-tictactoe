import { NetworksInfo } from "../models/network.model";
import { State, StateContext, Action, Store } from "@ngxs/store";
import { NetworkService } from "../modules/network-api/api/network.service";
import { Injector } from "@angular/core";
import { ActionDispatchHelperService } from "../utils/action-dispatch-helper.service";
import { tap } from "rxjs/operators";

export class NetworksInfoModel {
    networks: NetworksInfo[]
}

export class GetNetworksInfo {
    static readonly type = "[NetworksInfo] GetNetworksInfo"
}

@State<NetworksInfoModel>({
    name: "networkinfo",
    defaults: {
        networks: []
    }
})

export class NetworksInfoState {

    private _networkService: NetworkService

    constructor(private injector: Injector, private store: Store, private actionDispatcher: ActionDispatchHelperService) {
    }

    private get networkService() {
        if(!this._networkService) {
            this._networkService = this.injector.get(NetworkService)
        }
        return this._networkService;
    }
    
    @Action(GetNetworksInfo)
    getNetworksInfo(context: StateContext<NetworksInfoModel>, action: GetNetworksInfo) {
        let state = context.getState()
        // this.networkService.getNetworksInfo()
        return this.networkService.getNetworksInfo().pipe(tap(response => {
            console.log(response);
            return context.setState({
                networks: [...state.networks, ...response]
            });            
        }));
    }
}