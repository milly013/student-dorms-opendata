import { provideRouter, Routes } from '@angular/router';
import { LoginComponent } from './pages/login/login';
import { DormsComponent } from './pages/dorms/dorms';
import { DormDetailComponent } from './pages/dorm-detail/dorm-detail';
import { OpenDataHubComponent } from './pages/opendata-hub/opendata-hub';


export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'dorms', component: DormsComponent },
  { path: 'dorms/:id', component: DormDetailComponent },
  { path: 'opendata', component: OpenDataHubComponent },

];