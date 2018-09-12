import { Component, OnInit, Input } from '@angular/core';
import { BlockInfo } from '../../models/block.model';

@Component({
  selector: 'app-block',
  templateUrl: './block.component.html',
  styleUrls: ['./block.component.scss']
})
export class BlockComponent implements OnInit {

  @Input()block: BlockInfo;

  constructor() { }

  ngOnInit() {
  }

}
