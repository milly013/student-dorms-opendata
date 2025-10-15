import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface MoveInRequest {
  id: string;
  student_id: string;
  dorm_id: string;
  status: string;
  room_type: string;
  created_at: string;
  request_type: string;
  description: string;

  // dodatna polja za prikaz
  studentName?: string;
  dormName?: string;
}
export interface PopularDorm {
  dorm_id: string;
  request_count: number;
}

@Injectable({
  providedIn: 'root'
})
export class RequestService {
  private apiUrl = 'http://localhost:8000/requests';

  constructor(private http: HttpClient) { }

  getAllRequests(): Observable<MoveInRequest[]> {
    const token = localStorage.getItem('token');
    let headers = new HttpHeaders();
    if (token) {
      headers = headers.set('Authorization', `Bearer ${token}`);
    }

    return this.http.get<MoveInRequest[]>(`${this.apiUrl}/requests`, { headers });
  }
  createRequest(data: { student_id: string; dorm_id: string; room_type: string }) {
    const token = localStorage.getItem('token');
    let headers = new HttpHeaders();
    if (token) headers = headers.set('Authorization', `Bearer ${token}`);
    return this.http.post(`${this.apiUrl}/requests`, data, { headers });
  }

  approveRequest(id: string) {
    const token = localStorage.getItem('token');
    let headers = new HttpHeaders();
    if (token) headers = headers.set('Authorization', `Bearer ${token}`);
    return this.http.post(`${this.apiUrl}/requests/${id}/approve`, {}, { headers });
  }

  rejectRequest(id: string) {
    const token = localStorage.getItem('token');
    let headers = new HttpHeaders();
    if (token) headers = headers.set('Authorization', `Bearer ${token}`);
    return this.http.post(`${this.apiUrl}/requests/${id}/reject`, {}, { headers });
  }
  getPopularDorms(): Observable<PopularDorm[]> {
    const token = localStorage.getItem('token');
    let headers = new HttpHeaders();
    if (token) headers = headers.set('Authorization', `Bearer ${token}`);
    return this.http.get<PopularDorm[]>(`${this.apiUrl}/requests/popular-dorms`, { headers });
  }
  createMoveOutRequest(data: { student_id: string; dorm_id: string; room_type: string }) {
    const token = localStorage.getItem('token');
    let headers = new HttpHeaders();
    if (token) headers = headers.set('Authorization', `Bearer ${token}`);
    return this.http.post(`${this.apiUrl}/requests/move-out`, data, { headers });
  }
  createIssueRequest(data: any, token: string): Observable<any> {
    const headers = {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    };

    return this.http.post(`${this.apiUrl}/requests/issue`, data, { headers });
  }


}