import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Dorm, DormService } from '../../services/dorm';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RequestService } from '../../services/request.service';
import { AuthService } from '../../services/auth';
import { forkJoin } from 'rxjs';


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
  issueDescription: string = '';
  issueMessage: string | null = null;
  usersInDorm: any[] = []; 
  loadingUsers: boolean = false;
  userDormId: string | null = null;


  constructor(
    private route: ActivatedRoute,
    private dormsService: DormService,
    private requestService: RequestService,
    public authService: AuthService,
    private cdr: ChangeDetectorRef   

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
        this.userDormId = this.dorm.id;
        this.userBelongsToDormMap[this.dorm.id] = user.dorm_id === this.dorm.id;
      }
    });
  }

  loadDorm(id: string): void {
    this.dormsService.getDormById(id).subscribe((data) => {
      this.dorm = data;
      this.loadUsersInDorm();
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
  submitIssue(): void {
    if (!this.dorm) {
      this.issueMessage = 'Greška: dom nije učitan.';
      return;
    }

    const studentId = this.authService.getUserId();
    if (!studentId) {
      this.issueMessage = 'Molimo prijavite se prije prijave kvara.';
      return;
    }

    if (!this.issueDescription.trim()) {
      this.issueMessage = 'Molimo unesite opis kvara.';
      return;
    }

    const token = this.authService.getToken() || '';
    const issueData = {
      student_id: studentId,
      dorm_id: this.dorm.id,
      description: this.issueDescription,
      
    };

    this.requestService.createIssueRequest(issueData, token).subscribe({
      next: (res) => {
        this.issueMessage = '✅ Prijava kvara je uspješno poslana!';
        this.issueDescription = '';
      },
      error: (err) => {
        console.error('❌ Greška pri prijavi kvara:', err);
        this.issueMessage = err.status === 401
          ? '❌ Niste autorizovani. Prijavite se ponovo.'
          : '❌ Greška pri slanju prijave.';
      }
    });
  }
  loadUsersInDorm(): void {
  if (!this.dorm || !this.dorm.users || this.dorm.users.length === 0) {
    this.usersInDorm = [];
    return;
  }

  this.loadingUsers = true;

  const userObservables = this.dorm.users.map(userId =>
    this.authService.getPublicUserInfo(userId)
  );

  forkJoin(userObservables).subscribe({
    next: (users) => {
      this.usersInDorm = users;
      this.loadingUsers = false;
      this.cdr.detectChanges();

    },
    error: (err) => {
      console.error('❌ Greška pri učitavanju korisnika:', err);
      this.loadingUsers = false;
    }
  });
}


}
