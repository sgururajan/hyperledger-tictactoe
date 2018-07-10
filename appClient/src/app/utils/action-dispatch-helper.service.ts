import { Injectable } from '@angular/core';
import { Store, Actions, ofActionSuccessful } from '@ngxs/store';

@Injectable({
  providedIn: 'root'
})
export class ActionDispatchHelperService {

  constructor(private store: Store, private actions$:Actions) { }

  dispatchAndSubscribe(actionObject:any, callback:Function) {
    const subscription = this.actions$.pipe(ofActionSuccessful(actionObject)).subscribe(()=>{
      callback()
      subscription.unsubscribe()
    });
    this.store.dispatch(actionObject);
  }
}
