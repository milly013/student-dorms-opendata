import { Component, OnInit } from '@angular/core';
import { Opendata, OpenDorm } from '../../services/opendata';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-opendata-avgrating',
  imports: [CommonModule],
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

}