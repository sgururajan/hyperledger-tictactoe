import {Injectable} from '@angular/core';
import {MatDialog, MatDialogConfig, MatDialogRef} from "@angular/material";
import {DialogResult} from "../../models/dialog-result.model";


@Injectable({
  providedIn: 'root'
})
export class ModalDialogService {

  constructor(private dialog:MatDialog) {

  }

  async open(component) : Promise<DialogResult> {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose=true;
    dialogConfig.autoFocus=true;

    let dialogRef = this.dialog.open(component);
    let result = await dialogRef.afterClosed().toPromise();

    console.log("right after open");
    console.log(result);
    return result;
    //return this.dialog.open(component, dialogConfig);
    // let dialogRef = this.dialog.open(component, dialogConfig);
    // dialogRef.afterClosed().subscribe(dresult=> {
    //   console.log("dialog closed");
    //   console.log("dialog result");
    //   console.log(dresult);
    //   result = dresult;
    // });
    // console.log(result);
    // return Observable.create(result);
  }
}
