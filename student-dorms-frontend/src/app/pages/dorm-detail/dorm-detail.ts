import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Dorm, DormService } from '../../services/dorm';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RequestService } from '../../services/request.service';
import { AuthService } from '../../services/auth';


@Component({
  selector: 'app-dorm-detail',
  templateUrl: './dorm-detail.html',
  standalone: true,
  imports: [CommonModule, FormsModule]
})
export class DormDetailComponent implements OnInit {
  dorm?: Dorm;
  newRating: number = 0;
  message: string | null = null;
  roomType: string = 'single';
  requestMessage: string | null = null;
  moveOutMessage: string | null = null;
  userBelongsToDormMap: { [dormId: string]: boolean } = {};

  constructor(
    private route: ActivatedRoute,
    private dormsService: DormService,
    private requestService: RequestService,
    public authService: AuthService

  ) { }

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadDorm(id);
    }

    const userId = this.authService.getUserId();
    if (!userId) return;

    this.authService.getUserInfo(userId).subscribe(user => {
      if (this.dorm) {
        this.userBelongsToDormMap[this.dorm.id] = user.dorm_id === this.dorm.id;
      }
    });
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
        this.loadDorm(this.dorm!.id);
      },
      error: (err) => {
        console.error(err);
        this.message = err.status === 401
          ? '❌ Niste autorizovani. Prijavite se ponovo.'
          : '❌ Greška pri dodavanju ocjene.';
      }
    });
  }
  submitRequest(): void {
    if (!this.dorm) {
      this.requestMessage = 'Greška: dom nije učitan.';
      return;
    }

    const studentId = this.authService.getUserId();
    if (!studentId) {
      this.requestMessage = 'Molimo se prijavite prije slanja zahtjeva.';
      return;
    }
    const data = {
      student_id: studentId,
      dorm_id: this.dorm.id,
      room_type: this.roomType
    };

    this.requestService.createRequest(data).subscribe({
      next: (res) => {
        this.requestMessage = '✅ Zahtjev za useljenje je uspješno poslan!';
        this.roomType = 'single'; // reset forme
      },
      error: (err) => {
        console.error(err);
        this.requestMessage = err.status === 401
          ? '❌ Niste autorizovani. Prijavite se ponovo.'
          : '❌ Greška pri slanju zahtjeva.';
      }
    });
  }
  submitMoveOut(): void {
    if (!this.dorm) {
      this.moveOutMessage = 'Greška: dom nije učitan.';
      return;
    }

    const studentId = this.authService.getUserId();
    if (!studentId) {
      this.moveOutMessage = 'Molimo se prijavite prije slanja zahtjeva.';
      return;
    }

    const data = {
      student_id: studentId,
      dorm_id: this.dorm.id,
      room_type: '',
      requestType: 'move_in'
    };

    this.requestService.createMoveOutRequest(data).subscribe({
      next: (res) => {
        this.moveOutMessage = '✅ Zahtjev za iseljenje je uspješno poslan!';
      },
      error: (err) => {
        console.error(err);
        this.moveOutMessage = err.status === 401
          ? '❌ Niste autorizovani. Prijavite se ponovo.'
          : '❌ Greška pri slanju zahtjeva.';
      }
    });
  }

}
