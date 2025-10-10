import { provideRouter, Routes } from '@angular/router';
import { LoginComponent } from './pages/login/login';
import { DormsComponent } from './pages/dorms/dorms';
import { DormDetailComponent } from './pages/dorm-detail/dorm-detail';
import { OpenDataHubComponent } from './pages/opendata-hub/opendata-hub';
import { OpendataPrices } from './pages/opendata-prices/opendata-prices';
import { OpendataFreeSpots } from './pages/opendata-free-spots/opendata-free-spots';
import { OpendataFacilities } from './pages/opendata-facilities/opendata-facilities';
import { OpendataAvgratingimplements } from './pages/opendata-avgrating/opendata-avgrating';
import { Requests } from './pages/requests/requests';
import { PopularDorms } from './pages/popular-dorms/popular-dorms';
import { MoveOut } from './pages/move-out/move-out';
import { CreateDorm } from './pages/create-dorm/create-dorm';
import { Register } from './pages/register/register';


export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'dorms', component: DormsComponent },
  { path: 'dorms/:id', component: DormDetailComponent },
  { path: 'opendata', component: OpenDataHubComponent },
  { path: 'opendata/prices', component: OpendataPrices },
  { path: 'opendata/free-spots', component: OpendataFreeSpots},
  { path: 'opendata/facilities', component: OpendataFacilities},
  { path: 'opendata/avgRating', component: OpendataAvgratingimplements},
  { path: 'opendata/popularDorms', component: PopularDorms},
  { path: 'requests', component: Requests},
  { path: 'move-out', component: MoveOut},
  { path: 'create-dorm', component: CreateDorm},
  { path: 'register', component: Register}

];