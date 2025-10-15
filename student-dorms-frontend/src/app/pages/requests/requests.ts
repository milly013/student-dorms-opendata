import { Component, OnInit, NgZone } from '@angular/core';
import { MoveInRequest, RequestService } from '../../services/request.service';
import { DormService } from '../../services/dorm';
import { AuthService } from '../../services/auth';
import { CommonModule } from '@angular/common';
import { forkJoin, map, of, switchMap } from 'rxjs';

@Component({
  selector: 'app-requests',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './requests.html',
  styleUrls: ['./requests.css']
})
export class Requests implements OnInit {
  requests: MoveInRequest[] = [];
  loading = true;
  error: string | null = null;
  private loadingRequests = false;

  constructor(
    private requestService: RequestService,
    private dormService: DormService,
    private authService: AuthService,
    private ngZone: NgZone
  ) {}

  ngOnInit(): void {
    this.loadRequests();
  }

  loadRequests() {
    if (this.loadingRequests) return;
    this.loadingRequests = true;
    this.loading = true;
    this.error = '';
    this.requests = [];

    this.requestService.getAllRequests()
      .pipe(
        // Filtriramo samo move_in i move_out zahtjeve
        map(requests => requests.filter(r => r.request_type === 'move_in' || r.request_type === 'move_out')),
        switchMap(filteredRequests => {
          if (!filteredRequests || filteredRequests.length === 0) return of([]);

          const requestsWithDetails$ = filteredRequests.map(req => {
            const student$ = req.student_id
              ? this.authService.getUserInfo(req.student_id).pipe(
                  map(user => user?.username || 'Nepoznat student')
                )
              : of('Nepoznat student');

            const dorm$ = req.dorm_id
              ? this.dormService.getDormById(req.dorm_id).pipe(
                  map(dorm => dorm?.name || 'Nepoznat dom')
                )
              : of('Nepoznat dom');

            return forkJoin({ studentName: student$, dormName: dorm$ }).pipe(
              map(({ studentName, dormName }) => ({
                ...req,
                studentName,
                dormName
              }))
            );
          });

          return forkJoin(requestsWithDetails$);
        })
      )
      .subscribe({
        next: requests => {
          this.ngZone.run(() => {
            this.requests = requests;
            this.loading = false;
            this.loadingRequests = false;
            console.log('✅ Učitani zahtjevi:', this.requests);
          });
        },
        error: err => {
          this.ngZone.run(() => {
            console.error('❌ Greška pri učitavanju zahtjeva:', err);
            this.error = 'Greška pri učitavanju zahtjeva';
            this.loading = false;
            this.loadingRequests = false;
          });
        }
      });
  }

  // Odmah mijenjamo status na "approved"
  onApprove(req: MoveInRequest) {
    this.requestService.approveRequest(req.id).subscribe({
      next: () => {
        req.status = 'approved';
        console.log('✅ Zahtjev odobren:', req.id);
      },
      error: err => {
        console.error('❌ Greška pri odobravanju:', err);
        alert('Greška pri odobravanju zahtjeva.');
      }
    });
  }

  // Za move-in i move-out: mora se unijeti razlog pri odbijanju
  onReject(req: MoveInRequest) {
    const reason = prompt('Unesite razlog odbijanja zahtjeva:');
    if (!reason || reason.trim() === '') {
      alert('Morate unijeti razlog odbijanja.');
      return;
    }

    this.requestService.rejectRequestWithReason(req.id, reason).subscribe({
      next: () => {
        req.status = 'rejected';
        req.rejectionReason = reason; // čuvamo razlog za prikaz
        console.log('❌ Zahtjev odbijen:', req.id, 'Razlog:', reason);
      },
      error: err => {
        console.error('❌ Greška pri odbijanju zahtjeva:', err);
        alert('Greška pri odbijanju zahtjeva.');
      }
    });
  }
}
