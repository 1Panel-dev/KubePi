import {Routes} from '@angular/router';
import {AppComponent} from "./app.component";
import {TerminalComponent} from "./terminal/terminal.component";


export const routes: Routes = [
  {
    path: '', component: AppComponent, children: [
      {path: '', redirectTo: 'app', pathMatch: 'full'},
      {path: 'app', component: TerminalComponent},
      {path: '*', redirectTo: '', pathMatch: 'full'},
    ]
  },

]
