import 'zone.js'; 
import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { App } from './app/app';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { provideRouter } from '@angular/router';
import { routes } from './app/app.routes';

bootstrapApplication(App, {
  providers: [
    provideHttpClient(withFetch()),
    provideRouter(routes), // Dodaj provideHttpClient() ovde
    // Dodaj ostale providere ako ih imaš
    ...appConfig.providers // Ako appConfig ima druge providere
  ],
})
  .catch((err) => console.error(err));
