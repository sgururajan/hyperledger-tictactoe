import { Component, OnInit } from '@angular/core';
import { Store, Select } from '@ngxs/store';
import { Observable } from 'rxjs';
import { BlockStateModel } from '../../state/block.state';

@Component({
  selector: 'app-block-list',
  templateUrl: './block-list.component.html',
  styleUrls: ['./block-list.component.scss']
})
export class BlockListComponent implements OnInit {

  @Select(state=> state.blockstate.blockList) blockList$:Observable<BlockStateModel>;

  constructor(private store:Store) { }

  ngOnInit() {
    console.log('intializing blocklist component');
    this.blockList$.subscribe(resp=>console.log(`blocklist response: `, resp));    
  }

}
