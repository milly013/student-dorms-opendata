import { Component, OnInit } from '@angular/core';
import { MoveInRequest, RequestService } from '../../services/request.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-requests',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './requests.html',
  styleUrl: './requests.css'
})
export class Requests implements OnInit {
  
requests: MoveInRequest[] = [];
  loading = true;
  error: string | null = null;

  constructor(private requestService: RequestService) {}

  ngOnInit(): void {
    this.requestService.getAllRequests().subscribe({
      next: (data) => {
        this.requests = data;
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load requests';
        this.loading = false;
        console.error(err);
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