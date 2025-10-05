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
}