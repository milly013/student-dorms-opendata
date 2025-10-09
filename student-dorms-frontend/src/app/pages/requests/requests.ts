
import { Component, OnInit } from '@angular/core';
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
  private loadingRequests = false; // ✅ zaštita od duplog poziva

  constructor(
    private requestService: RequestService,
    private dormService: DormService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    this.loadRequests();
  }

  loadRequests() {
    if (this.loadingRequests) return; // sprečava dupli poziv
    this.loadingRequests = true;

    this.loading = true;
    this.error = '';
    this.requests = []; // ✅ resetuj listu pre svakog učitavanja

    this.requestService.getAllRequests().pipe(
      switchMap((requests) => {
        if (!requests || requests.length === 0) return of([]); // ako nema zahtjeva

        const requestsWithDetails$ = requests.map(req => {
          const student$ = req.student_id
            ? this.authService.getUserInfo(req.student_id).pipe(
                map(user => {
                  console.log('👤 Dobavljen user:', user);
                  return user?.username || user?.name || 'Nepoznat student';
                })
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
        this.requests = requests;
        this.loading = false;
        this.loadingRequests = false;
        console.log('✅ Učitani zahtjevi:', this.requests);
      },
      error: (err) => {
        console.error('❌ Greška:', err);
        this.error = 'Greška pri učitavanju zahtjeva';
        this.loading = false;
        this.loadingRequests = false;
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
