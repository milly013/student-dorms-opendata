import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Dorm , DormService  } from '../../services/dorm';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';


@Component({
  selector: 'app-dorm-detail',
  templateUrl: './dorm-detail.html',
  standalone: true,
  imports:[CommonModule,FormsModule]
})
export class DormDetailComponent implements OnInit {
  dorm?: Dorm;
  newRating: number = 0;
  message: string | null = null;

  constructor(
    private route: ActivatedRoute,
    private dormsService: DormService
  ) {}

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadDorm(id);
    }
  }

  loadDorm(id: string): void {
    this.dormsService.getDormById(id).subscribe((data) => {
      this.dorm = data;
    });
  }

  submitRating(): void {
    if (!this.dorm || this.newRating < 1 || this.newRating > 5) {
      this.message = 'Molimo unesite ocjenu između 1 i 5.';
      return;
    }

    this.dormsService.addRating(this.dorm.id, this.newRating).subscribe({
      next: () => {
        this.message = '✅ Ocjena je uspješno dodata!';
        this.newRating = 0;
        this.loadDorm(this.dorm!.id); // ponovo učitaj dom da se vidi nova prosječna ocjena
      },
      error: (err) => {
        console.error(err);
        this.message = err.status === 401
          ? '❌ Niste autorizovani. Prijavite se ponovo.'
          : '❌ Greška pri dodavanju ocjene.';
      }
    });
  }
}
