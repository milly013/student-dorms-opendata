import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from './auth';

export interface Dorm {
  id: string;
  name: string;
  address: string;
  city: string;
  capacity: number;
  price: number;
  occupied: number;
  type: string;
  amenities: string[];
  description?: string;
  ratings?: { user_id: string; score: number }[];
  comments?: { user_id: string; message: string; date: string }[];
  preferences?: string[];
}

@Injectable({
  providedIn: 'root'
})
export class DormService {
  private apiUrl = 'http://localhost:8000/dorms'; 

  constructor(private http: HttpClient, private authService: AuthService) {
    
  }

  createDorm(dorm: Dorm): Observable<any> {
    const token = this.authService.getToken(); // ili localStorage.getItem('token')
    const headers = token ? new HttpHeaders({ Authorization: `Bearer ${token}` }) : undefined;

    return this.http.post(`${this.apiUrl}/dorms`, dorm, { headers });
  }

  getDorms(): Observable<Dorm[]> {
    return this.http.get<Dorm[]>(this.apiUrl);
  }

  getDormById(id: string): Observable<Dorm> {
    return this.http.get<Dorm>(`${this.apiUrl}/dorms/${id}`);
  }

  getAllDorms(): Observable<Dorm[]> {
    return this.http.get<Dorm[]>(`${this.apiUrl}/dorms`);
  }

  addRating(dormId: string, score: number): Observable<any> {
    const token = localStorage.getItem('token'); 

    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`
    });

    return this.http.post(
      `${this.apiUrl}/dorms/rating`,
      { dorm_id: dormId, score: score },
      { headers }
    );
  }

}
