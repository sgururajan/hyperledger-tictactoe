import { Component, OnInit } from '@angular/core';
import { Store, Select, Actions, ofActionSuccessful } from '@ngxs/store';
import { Observable } from 'rxjs';
import { NetworksInfoModel, GetNetworksInfo } from '../../state/network-info.state';
import { NetworkService } from '../../modules/network-api/api/network.service';

@Component({
  selector: 'app-network-list',
  templateUrl: './network-list.component.html',
  styleUrls: ['./network-list.component.scss']
})
export class NetworkListComponent implements OnInit {

  @Select(store=> store.networkinfo.networks) networks$:Observable<NetworksInfoModel>;

  constructor(private store: Store, private actions$: Actions) {
    // console.log(this.networkService);
    // this.networkService.getNetworksInfo().subscribe(response=>console.log(response)); 
    this.actions$.pipe(ofActionSuccessful(GetNetworksInfo)).subscribe((res)=>{
      console.log("successfully loaded networks info");
      console.log(res);
    });
  }

  ngOnInit() {
    console.log("initializing networklist component");
    this.store.dispatch(new GetNetworksInfo())
  }

}
