import { BlockInfo } from "../models/block.model";
import { State, Store, StateContext, Action } from "@ngxs/store";

import * as _ from 'lodash';

export interface BlockStateModel {
    blockList: BlockInfo[]
}

export class UpdateBlockList {
    static readonly type="[BlockState] UpdateBlockList";
    constructor(public blocks:BlockInfo[]){}
}

@State({
    name: "blockstate",
    defaults: {
        blockList:[]
    }
})

export class BlockState {
    constructor (private store:Store){}

    @Action(UpdateBlockList)
    updateBlockList(context:StateContext<BlockStateModel>, action: UpdateBlockList) {
        const state = context.getState();
        console.log(`block list state: `, state.blockList);
        //const modifiedState = {...state.blockList, ...action.blocks};
        console.log(`action input: `, action.blocks);
        action.blocks.forEach(b=> state.blockList.push(b));
        console.log(`modified state: `, state);
        return context.patchState({
            blockList: _.unionBy(state.blockList, action.blocks, 'blockNumber')
        });
    }
}