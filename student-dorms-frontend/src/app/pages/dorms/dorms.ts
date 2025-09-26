import { Component, OnInit } from '@angular/core';
import { Dorm, DormService } from '../../services/dorm';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-dorms',
  imports: [CommonModule],
  standalone: true,
  templateUrl: './dorms.html',
  styleUrl: './dorms.css'
})
export class DormsComponent implements OnInit {
  dorms: Dorm[] = [];
  loading = true;
  error: string | null = null;

  constructor(private dormService: DormService, private router: Router) {}

  ngOnInit(): void {
    this.dormService.getAllDorms().subscribe({
      next: (data) => {
        this.dorms = data;
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load dorms';
        this.loading = false;
        console.error(err);
      }
    });
  }
  viewDorm(id: string) {
    this.router.navigate(['/dorms', id]);
  }
}
