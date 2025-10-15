
import { Component, NgZone, OnInit } from '@angular/core';
import { MoveInRequest, RequestService } from '../../services/request.service';
import { CommonModule } from '@angular/common';
import { DormService } from '../../services/dorm';
import { AuthService } from '../../services/auth';
import { forkJoin, map, switchMap, of } from 'rxjs';

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
  ) { }

  ngOnInit(): void {
    this.loadRequests();
  }

  loadRequests() {
    if (this.loadingRequests) return;
    this.loadingRequests = true;

    this.loading = true;
    this.error = '';
    this.requests = [];

    this.requestService.getAllRequests().pipe(
      switchMap((requests) => {
        if (!requests || requests.length === 0) return of([]);

        // ✅ Filtriramo samo move_in i move_out zahtjeve
        const filteredRequests = requests.filter(req =>
          req.request_type === 'move_in' || req.request_type === 'move_out'
        );

        if (filteredRequests.length === 0) return of([]);

        const requestsWithDetails$ = filteredRequests.map(req => {
          const student$ = req.student_id
            ? this.authService.getUserInfo(req.student_id).pipe(
              map(user => user?.username || user?.name || 'Nepoznat student')
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
    ).subscribe({
      next: (requests) => {
        this.ngZone.run(() => { // 👈 forsira update UI-a
          this.requests = requests;
          this.loading = false;
          this.loadingRequests = false;
          console.log('✅ Učitani zahtjevi:', this.requests);
        });
      },
      error: (err) => {
        this.ngZone.run(() => { // 👈 i ovdje isto
          console.error('❌ Greška:', err);
          this.error = 'Greška pri učitavanju zahtjeva';
          this.loading = false;
          this.loadingRequests = false;
        });
      }
    });
  }

  approve(reqId: string) {
    this.requestService.approveRequest(reqId).subscribe({
      next: () => {
        this.requests = this.requests.map(r =>
          r.id === reqId ? { ...r, status: 'approved' } : r
        );
      },
      error: (err) => {
        console.error('Failed to approve request:', err);
        alert(`Approve error: ${err.status} - ${err.message}\nDetails: ${JSON.stringify(err.error)}`);
      }
    });
  }

  reject(reqId: string) {
    this.requestService.rejectRequest(reqId).subscribe({
      next: () => {
        this.requests = this.requests.map(r =>
          r.id === reqId ? { ...r, status: 'rejected' } : r
        );
      },
      error: (err) => {
        console.error('Failed to reject request:', err);
        alert(`Reject error: ${err.status} - ${err.message}\nDetails: ${JSON.stringify(err.error)}`);
      }
    });
  }
}
