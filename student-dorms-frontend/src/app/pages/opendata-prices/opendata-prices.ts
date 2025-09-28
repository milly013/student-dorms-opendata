import { Component, OnInit } from '@angular/core';
import { CityDorms, DormPrice, Opendata } from '../../services/opendata';
import { CommonModule } from '@angular/common';
import { Dorm } from '../../services/dorm';

@Component({
  selector: 'app-opendata-prices',
  imports: [CommonModule],
  standalone: true,
  templateUrl: './opendata-prices.html',
  styleUrl: './opendata-prices.css'
})
export class OpendataPrices implements OnInit{

  constructor(private openDataService: Opendata) {}
  
  cityDorms: CityDorms[] = [];
  loading = true;
  error = '';

  ngOnInit(): void {
    // Dohvati sve dormove
    this.openDataService.getDorms().subscribe({
      next: (dorms: Dorm[]) => {
        const cityMap: { [city: string]: DormPrice[] } = {};

        // Grupisanje domova po gradu
        dorms.forEach(d => {
          if (!cityMap[d.city]) {
            cityMap[d.city] = [];
          }
          cityMap[d.city].push({ name: d.name, price: d.price });
        });

        // Formiranje niza CityDorms sa avgPrice i listom dormova
        this.cityDorms = Object.entries(cityMap).map(([city, dorms]) => {
          const totalPrice = dorms.reduce((sum, d) => sum + d.price, 0);
          const avgPrice = dorms.length ? totalPrice / dorms.length : 0;
          return { city, dorms, averagePrice: avgPrice };
        });

        this.loading = false;
      },
      error: (err) => {
        this.error = 'Neuspješno dohvaćanje podataka';
        this.loading = false;
      }
    });
  }

}
