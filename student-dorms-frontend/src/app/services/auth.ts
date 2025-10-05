import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';

interface LoginResponse {
  token: string;
  
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:8000/auth';

  constructor(private http: HttpClient) {}

  login(email: string, password: string): Observable<LoginResponse> {
    
    return this.http.post<LoginResponse>(`${this.apiUrl}/login`, { email, password })
      .pipe(
        tap((res: LoginResponse) => {   
          localStorage.setItem('token', res.token);
        })
      );
  }

  logout(): void {
    localStorage.removeItem('token');
  }
  getUserId(): string | null {
    const token = this.getToken();
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.user_id || null;  
    } catch (e) {
      console.error('Nevalidan token', e);
      return null;
    }
  }

  getToken(): string | null {
    return localStorage.getItem('token');
  }
   isLoggedIn(): boolean {
    return !!localStorage.getItem('token');
  }
  getUserRole(): string | null {
    const token = this.getToken();
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.role || null;
    } catch (e) {
      return null;
    }
  }
   isAdmin(): boolean {
    return this.getUserRole() === 'admin';
  }
  getUserInfo(userId: string): Observable<any> {
  return this.http.get<any>(`${this.apiUrl}/users/${userId}`);
}

isInDorm(): Observable<boolean> {
  const userId = this.getUserId();
  if (!userId) return new Observable<boolean>((observer) => {
    observer.next(false);
    observer.complete();
  });

  return this.getUserInfo(userId).pipe(
    tap(user => console.log("User info:", user)),
    // vraćamo samo true/false
    map(user => !!user.inDorm)
  );
}
}