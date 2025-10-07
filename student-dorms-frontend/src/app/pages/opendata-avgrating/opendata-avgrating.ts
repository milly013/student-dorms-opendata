import { Component, OnInit } from '@angular/core';
import { Opendata, OpenDorm } from '../../services/opendata';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-opendata-avgrating',
  imports: [CommonModule],
  standalone:true,
  templateUrl: './opendata-avgrating.html',
  styleUrl: './opendata-avgrating.css'
})
export class OpendataAvgratingimplements implements OnInit {

  dorms: OpenDorm[] = [];
  loading = true;
  error: string | null = null;

  constructor(private opendataService: Opendata) {}

  ngOnInit(): void {
    this.opendataService.getDorms().subscribe({
  next: (data) => {
    this.dorms = data.map(d => ({
      ...d,
      average_rating: (d as any).average_rating ?? 0,
      occupancy_rate: (d as any).occupancy_rate ?? 0,
      comments_count: (d as any).comments_count ?? 0,
      tags: (d as any).tags ?? []
      }));
      this.loading = false;
    },
    error: (err) => {
      this.error = "Failed to load dorms";
      console.error(err);
      this.loading = false;
    }
  });
  }
   exportToJSON() {
    const dataStr = JSON.stringify(this.dorms, null, 2);
    const blob = new Blob([dataStr], { type: 'application/json' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'dorms.json';
    a.click();
    window.URL.revokeObjectURL(url);
  }

  exportToCSV() {
    if (!this.dorms.length) return;

    const header = Object.keys(this.dorms[0]).join(',');
    const rows = this.dorms.map(d => Object.values(d).join(','));
    const csvContent = [header, ...rows].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'dorms.csv';
    a.click();
    window.URL.revokeObjectURL(url);
  }

}