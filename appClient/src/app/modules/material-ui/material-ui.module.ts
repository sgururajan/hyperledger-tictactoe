import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { FlexLayoutModule } from "@angular/flex-layout";
import { LayoutModule } from "@angular/cdk/layout";
import {Material_UI_Modules} from "./material-ui.header";

@NgModule({
  imports: [
    BrowserAnimationsModule,
    FlexLayoutModule,
    LayoutModule,
    ...Material_UI_Modules,
  ],
  exports: [
    BrowserAnimationsModule,
    FlexLayoutModule,
    LayoutModule,
    ...Material_UI_Modules,
  ]
})
export class MaterialUiModule { }
