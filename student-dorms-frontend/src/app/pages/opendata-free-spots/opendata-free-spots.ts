import { Component, OnInit } from '@angular/core';
import { Opendata } from '../../services/opendata';
import { NgChartsModule } from 'ng2-charts';
import { ChartOptions, ChartType, ChartData } from 'chart.js';

export interface DormFreeSpots {
  name: string;
  freeSpots: number;
}

export interface CityFreeSpots {
  city: string;
  dorms: DormFreeSpots[];
  totalFreeSpots: number;
}

import { CommonModule } from '@angular/common';
import { Dorm } from '../../services/dorm';

@Component({
  selector: 'app-opendata-free-spots',
  imports: [CommonModule, NgChartsModule],
  standalone: true,
  templateUrl: './opendata-free-spots.html',
  styleUrl: './opendata-free-spots.css'
})
export class OpendataFreeSpots implements OnInit {
  cityDorms: CityFreeSpots[] = [];
  loading = true;
  error = '';

  // Procenat popunjenosti po gradu
  occupancy: Record<string, number> = {};

  // Pie chart config (isti za sve)
  pieChartType = 'pie' as const; 
  pieChartLegend = true;
  pieChartOptions: ChartOptions<'pie'> = {
    responsive: true,
    plugins: {
      legend: { position: 'top' },
      tooltip: { enabled: true }
    }
  };

  constructor(private openDataService: Opendata) {}

  ngOnInit(): void {
    this.openDataService.getDorms().subscribe({
      next: (dorms: Dorm[]) => {
        const cityMap: { [city: string]: DormFreeSpots[] } = {};

        // Grupisanje domova po gradu i računanje slobodnih mesta
        dorms.forEach(d => {
          if (!cityMap[d.city]) cityMap[d.city] = [];
          cityMap[d.city].push({ name: d.name, freeSpots: d.capacity - d.occupied });
        });

        // Formiranje niza sa ukupnim slobodnim mestima po gradu
        this.cityDorms = Object.entries(cityMap).map(([city, dorms]) => {
          const totalFreeSpots = dorms.reduce((sum, d) => sum + d.freeSpots, 0);
          return { city, dorms, totalFreeSpots };
        });

        this.loading = false;
      },
      error: () => {
        this.error = 'Neuspješno dohvaćanje podataka';
        this.loading = false;
      }
    });
    this.openDataService.getOccupancyPerCity().subscribe({
      next: (data) => {
        this.occupancy = data;
      },
      error: () => {
        this.error = 'Greška pri učitavanju popunjenosti';
      }
    });
  }
  getPieChartData(city: string): ChartData<'pie', number[], string | string[]> {
    const percent = this.occupancy[city] || 0;
    return {
      labels: ['Popunjeno', 'Slobodno'],
      datasets: [
        {
          data: [percent, 100 - percent]
        }
      ]
    };
  }
  exportJSON(): void {
    const dataStr = JSON.stringify(this.cityDorms, null, 2);
    const blob = new Blob([dataStr], { type: 'application/json' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'free_spots.json';
    a.click();
    window.URL.revokeObjectURL(url);
  }

  exportCSV(): void {
    let csv = 'City,Dorm,FreeSpots,TotalFreeSpots\n';
    this.cityDorms.forEach(city => {
      city.dorms.forEach(dorm => {
        csv += `${city.city},${dorm.name},${dorm.freeSpots},${city.totalFreeSpots}\n`;
      });
    });

    const blob = new Blob([csv], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'free_spots.csv';
    a.click();
    window.URL.revokeObjectURL(url);
  }
}
