import {Routes} from '@angular/router';
import {AppComponent} from "./app.component";
import {TerminalComponent} from "./terminal/terminal.component";
import {LoggingComponent} from "./logging/logging.component";


export const routes: Routes = [
  {
    path: '', component: AppComponent, children: [
      {path: '', redirectTo: 'app', pathMatch: 'full'},
      {path: 'app', component: TerminalComponent},
      {path: 'logging', component: LoggingComponent},
      {path: '*', redirectTo: '', pathMatch: 'full'},
    ]
  },

]
