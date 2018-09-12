import { Network, Organization } from "../models/network.model";
import {State, Action, StateContext, Selector} from "@ngxs/store";
import { NetworkService } from "../api-service/api/network.service";
import { Injector } from "@angular/core";
import {tap} from "rxjs/operators";

export interface NetworkModel {
    networks: Network[],
    currentNetwork?: Network,
    currentOrganization?: Organization,
    organizations?: Organization[]
}

export class GetNetworks {
    static readonly type="[NetworkState] GetNetworks";
}

export class SelectNetwork {
  static readonly type="[NetworkState] SelectNetwork";
  constructor(public networkName:String){}
}

export class SetCurrentNetworkOrganizations {
  static readonly type="[NetworkState] SetCurrentNetworkOrganizations";
}

export class SelectOrganization {
  static readonly type="[NetworkState] SelectOrganization";
  constructor(public org:Organization){}
}

@State<NetworkModel>({
    name: 'networkstate',
    defaults: {
        networks: [],
        currentNetwork: undefined,
        currentOrganization: undefined,
        organizations: []
    }
})
export class NetworkState {
  @Selector() static currentOrganization(state:NetworkModel): Organization {
    return state.currentOrganization;
  }

  @Selector() static currentNetwork(state:NetworkModel): Network {
    return state.currentNetwork;
  }

  private _networkService: NetworkService;

    constructor(private injector: Injector){

    }

    private get networkService() {
        if (!this._networkService) {
            this._networkService = this.injector.get(NetworkService);
        }
        return this._networkService;
    }

    @Action(GetNetworks)
    getNetworks(context: StateContext<NetworkModel>, action: GetNetworks) {
        console.log("getNetworks");
        return this.networkService.getNetworks().pipe(tap(response=>{
          console.log(response);
            return context.setState({
                networks: response,
                currentNetwork: undefined,
                currentOrganization: undefined,
                organizations: []
            });
        }));
    }

    @Action(SelectNetwork)
    selectNetwork(context: StateContext<NetworkModel>, action:SelectNetwork) {
      console.log("SelectNetwork");
      if(!action.networkName || action.networkName==="")return;
      let state = context.getState();

      console.log(state.networks.find(x=> x.name===action.networkName));
      return context.patchState({
        currentNetwork: state.networks.find(n=>n.name===action.networkName),
      })
    }

    @Action(SetCurrentNetworkOrganizations)
    setCurrentNetworkOrganizations(contex: StateContext<NetworkModel>, action: SetCurrentNetworkOrganizations) {
      console.log("setCurrentNetworkOrganizations");
      let state = contex.getState();
      if(!state.currentNetwork) return;
      return contex.patchState({
        organizations: state.currentNetwork.organizations,
      });
    }

    @Action(SelectOrganization)
    selectOrganization(context: StateContext<NetworkModel>, action: SelectOrganization) {
      console.log('selecting current organization: ' + action.org.name);
      return context.patchState({
        currentOrganization: action.org
      })
    }
}
