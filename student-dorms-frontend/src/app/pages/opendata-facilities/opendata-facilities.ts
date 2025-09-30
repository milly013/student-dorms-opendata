import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Opendata, FacilitiesSummary } from '../../services/opendata';

@Component({
  selector: 'app-opendata-facilities',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './opendata-facilities.html',
  styleUrl: './opendata-facilities.css'
})
export class OpendataFacilities implements OnInit {
  facilitiesSummary: FacilitiesSummary = {};
  loading = true;
  error = '';

  constructor(private openDataService: Opendata) {}

  ngOnInit(): void {
    this.openDataService.getFacilitiesSummary().subscribe({
      next: (data) => {
        this.facilitiesSummary = data;
        this.loading = false;
      },
      error: () => {
        this.error = 'Neuspješno dohvaćanje podataka';
        this.loading = false;
      }
    });
  }

  // helper: lista svih amenitija za jedan grad
  getFacilitiesList(city: string): { name: string; count: number }[] {
    const cityFacilities = this.facilitiesSummary[city];
    return Object.entries(cityFacilities).map(([name, count]) => ({
      name,
      count
    }));
  }
}
