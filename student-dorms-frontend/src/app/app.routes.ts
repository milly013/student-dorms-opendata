import { provideRouter, Routes } from '@angular/router';
import { LoginComponent } from './pages/login/login';
import { bootstrapApplication } from '@angular/platform-browser';
import { provideHttpClient } from '@angular/common/http';
import { App } from './app';
import { Dorms } from './pages/dorms/dorms';

export const routes: Routes = [
  { path: '', redirectTo: 'login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'dorms', component: Dorms },
];