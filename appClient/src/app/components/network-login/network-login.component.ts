import {ChangeDetectorRef, Component, OnInit, Optional} from '@angular/core';
import {Select, Store} from "@ngxs/store";
import {Observable} from "rxjs/internal/Observable";
import {NetworkModel, SelectOrganization} from "../../state/network.state";
import {Organization} from "../../models/network.model";
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {MatDialogRef} from "@angular/material";
import {DialogResult} from "../../models/dialog-result.model";

@Component({
  selector: 'app-network-login',
  templateUrl: './network-login.component.html',
  styleUrls: ['./network-login.component.scss']
})
export class NetworkLoginComponent implements OnInit {

  @Select(state=> state.networkstate.organizations) orgsList$:Observable<NetworkModel>;
  orgFormControl:FormControl = new FormControl('', [Validators.required]);

  constructor(private store:Store, @Optional() private dialogRef: MatDialogRef<NetworkLoginComponent>) { }

  ngOnInit() {

  }

  onLogin() {
    if (!this.orgFormControl.valid)return;
    console.log(this.orgFormControl.value);
    this.store.dispatch(new SelectOrganization(this.orgFormControl.value));
    if (this.dialogRef) {
      this.dialogRef.close(<DialogResult>{ Success: true});
    }
  }

}
