import { Component, OnInit, NgZone  } from '@angular/core';
import { MoveInRequest, RequestService } from '../../services/request.service';
import { DormService } from '../../services/dorm';
import { AuthService } from '../../services/auth';
import { CommonModule } from '@angular/common';
import { forkJoin, map, of, switchMap } from 'rxjs';

@Component({
  selector: 'app-issue-requests',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './issue-requests.html',
  styleUrl: './issue-requests.css'
})
export class IssueRequests implements OnInit {
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
    this.loadIssueRequests();
  }

  loadIssueRequests() {
    if (this.loadingRequests) return;
    this.loadingRequests = true;
    this.loading = true;
    this.error = '';
    this.requests = [];

    this.requestService.getAllRequests().pipe(
      // 🔹 Filtriramo samo zahtjeve tipa "issue"
      map(requests => requests.filter(r => r.request_type === 'issue_report')),

      switchMap((filteredRequests) => {
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
    ).subscribe({
      next: (requests) => {
        this.ngZone.run(() => {
          this.requests = requests;
          this.loading = false;
          this.loadingRequests = false;
          console.log('✅ Učitani zahtjevi za kvarove:', this.requests);
        });
      },
      error: (err) => {
        this.ngZone.run(() => {
          console.error('❌ Greška pri učitavanju zahtjeva:', err);
          this.error = 'Greška pri učitavanju zahtjeva za kvarove';
          this.loading = false;
          this.loadingRequests = false;
        });
      }
    });
  }
}