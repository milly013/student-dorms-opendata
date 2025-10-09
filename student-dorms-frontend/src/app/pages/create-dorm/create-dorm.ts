import { Component } from '@angular/core';
import { Dorm, DormService } from '../../services/dorm';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-create-dorm',
  imports: [CommonModule, FormsModule],
  standalone: true,
  templateUrl: './create-dorm.html',
  styleUrl: './create-dorm.css'
})
export class CreateDorm {
dorm: Dorm = {
    id: '',
    name: '',
    address: '',
    city: '',
    capacity: 0,
    price: 0,
    occupied: 0,
    type: 'muški',
    amenities: [],
    description: '',
    ratings: [],
    comments: [],
    preferences: [],
  };
  newAmenity: string = '';
  message: string | null = null;

  constructor(private dormService: DormService) {}

  addAmenity(): void {
    if (this.newAmenity.trim() && !this.dorm.amenities.includes(this.newAmenity.trim())) {
      this.dorm.amenities.push(this.newAmenity.trim());
      this.newAmenity = '';
    }
  }

  removeAmenity(index: number): void {
    this.dorm.amenities.splice(index, 1);
  }

  submitDorm(): void {
    this.dormService.createDorm(this.dorm).subscribe({
      next: () => {
        this.message = '✅ Dom je uspešno kreiran!';
        this.dorm = {
          id: '',
          name: '',
          address: '',
          city: '',
          capacity: 0,
          price: 0,
          occupied: 0,
          type: 'muški',
          amenities: [],
          description: '',
          ratings: [],
          comments: [],
          preferences: [],
        };
      },
      error: (err) => {
        console.error(err);
        this.message = err.status === 401
          ? '❌ Niste autorizovani.'
          : '❌ Greška pri kreiranju doma.';
      }
    });
  }
}

