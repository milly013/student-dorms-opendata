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
  exportToJSON(): void { 
  const jsonData = JSON.stringify(this.cityDorms, null, 2); 
  const blob = new Blob([jsonData], { type: 'application/json' }); 
  const url = window.URL.createObjectURL(blob); 
  const a = document.createElement('a'); 
  a.href = url; 
  a.download = 'city-dorms.json'; 
  a.click(); 
  window.URL.revokeObjectURL(url); 
} 
  exportToCSV(): void { 
    let csv = 'City,Dorm,Price,Average Price\n'; 
    this.cityDorms.forEach(cityData => { 
    cityData.dorms.forEach(dorm => { 
    csv += `${cityData.city},${dorm.name},${dorm.price},${cityData.averagePrice}\n`;
    }); 
  }); 
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' }); 
  const url = window.URL.createObjectURL(blob); 
  const a = document.createElement('a'); 
  a.href = url; 
  a.download = 'city-dorms.csv'; 
  a.click(); 
  window.URL.revokeObjectURL(url); 
}

}
