import { Component, OnInit, AfterViewInit } from '@angular/core';
import { Dorm, DormService } from '../../services/dorm';
import { PopularDorm, RequestService } from '../../services/request.service';
import { forkJoin } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-popular-dorms',
  imports: [CommonModule],
  standalone:true,
  templateUrl: './popular-dorms.html',
  styleUrls: ['./popular-dorms.css']
})
export class PopularDorms implements OnInit, AfterViewInit  {
  popularDorms: Dorm[] = [];

  constructor(
    private requestService: RequestService,
    private dormService: DormService
  ) {}
  ngAfterViewInit(): void {
    this.loadPopularDorms();
  }

  ngOnInit(): void {
    this.loadPopularDorms();
  }

  loadPopularDorms() {
    this.requestService.getPopularDorms().subscribe({
      next: (popular: PopularDorm[]) => {
        // Za svaki dormID pozovi DormService da dobiješ detalje
        const dormRequests = popular.map(d => this.dormService.getDormById(d.dorm_id));
        
        // forkJoin čeka sve HTTP pozive da se završe
        forkJoin(dormRequests).subscribe({
          next: (dorms: Dorm[]) => {
            // Sortiraj dormove prema originalnom redu popularnosti
            this.popularDorms = popular.map(p => dorms.find(d => d.id === p.dorm_id)!);
          },
          error: (err) => console.error('Failed to load dorm details', err)
        });
      },
      error: (err) => console.error('Failed to load popular dorms', err)
    });
  }
  exportToJSON() {
    const jsonData = JSON.stringify(this.popularDorms, null, 2);
    const blob = new Blob([jsonData], { type: 'application/json' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'popular-dorms.json';
    a.click();
    window.URL.revokeObjectURL(url);
  }

  // 📌 Export u CSV
  exportToCSV() {
    if (!this.popularDorms.length) return;

    const headers = ['Rank', 'Name', 'Address', 'City', 'Type', 'Capacity', 'Occupied'];
    const rows = this.popularDorms.map((dorm, i) => [
      i + 1,
      dorm.name,
      dorm.address,
      dorm.city,
      dorm.type,
      dorm.capacity,
      dorm.occupied
    ]);

    const csvContent =
      [headers, ...rows]
        .map(row => row.map(item => `"${item}"`).join(','))
        .join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'popular-dorms.csv';
    a.click();
    window.URL.revokeObjectURL(url);
  }
}