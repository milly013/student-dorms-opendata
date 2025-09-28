import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

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

  constructor(private http: HttpClient) {}

  getDorms(): Observable<Dorm[]> {
    return this.http.get<Dorm[]>(this.apiUrl);
  }

  getDormById(id: string): Observable<Dorm> {
    return this.http.get<Dorm>(`${this.apiUrl}/dorms/${id}`);
  }

  getAllDorms(): Observable<Dorm[]> {
    return this.http.get<Dorm[]>(`${this.apiUrl}/dorms`);
  }
}
