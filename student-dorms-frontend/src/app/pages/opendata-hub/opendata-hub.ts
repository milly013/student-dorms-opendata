import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-opendata-hub',
  imports: [CommonModule,RouterLink],
  standalone: true,
  templateUrl: './opendata-hub.html',
  styleUrl: './opendata-hub.css'
})
export class OpenDataHubComponent {
  statistics = [
    {
      title: 'Prosječna cijena po gradu',
      description: 'Pregledajte prosječnu cijenu smještaja u različitim gradovima.',
      route: '/opendata/prices'
    },
    {
      title: 'Prosjecne ocjene.',
      description: 'Pregled prosjecnih ocjena domova.',
      route: '/opendata/avgRating'
    },
    {
      title: 'Trend popunjenosti',
      description: 'Pratite trend popunjenosti domova kroz vrijeme.',
      route: '/opendata/trends'
    },
    {
      title: 'Amenitiji po gradu',
      description: 'Lista svih dostupnih amenitija i njihov broj po gradu.',
      route: '/opendata/facilities'
    },
    {
      title: 'Rang lista gradova po slobodnim mjestima',
      description: 'Gradovi sa najviše slobodnih mjesta u domovima.',
      route: '/opendata/free-spots'
    }
  ];
}