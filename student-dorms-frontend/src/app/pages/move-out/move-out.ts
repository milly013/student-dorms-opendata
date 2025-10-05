import { Component, OnInit } from '@angular/core';
import { AuthService } from '../../services/auth';
import { RequestService } from '../../services/request.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-move-out',
  imports: [CommonModule,FormsModule],
  standalone:true,
  templateUrl: './move-out.html',
  styleUrl: './move-out.css'
})
export class MoveOut implements OnInit {
  dormId: string = '';
  roomType: string = '';
  message: string = '';
  error: string = '';
  studentId: string | null = null;

  constructor(
    private authService: AuthService,
    private requestService: RequestService
  ) {}

  ngOnInit(): void {
    // Dohvati ID korisnika iz tokena
    this.studentId = this.authService.getUserId();
  }

  submitMoveOut(): void {
    if (!this.studentId) {
      this.error = 'Korisnik nije pronađen!';
      return;
    }

    if (!this.dormId || !this.roomType) {
      this.error = 'Molimo popunite Dorm ID i Tip sobe.';
      return;
    }

    this.requestService.createMoveOutRequest({
      student_id: this.studentId,
      dorm_id: this.dormId,
      room_type: this.roomType
    }).subscribe({
      next: (res) => {
        this.message = 'Zahtjev za iseljenje je uspješno poslan.';
        this.error = '';
      },
      error: (err) => {
        this.error = 'Došlo je do greške prilikom slanja zahtjeva.';
        this.message = '';
        console.error(err);
      }
    });
  }
}