import { Component } from '@angular/core';
import { AuthService } from '../../services/auth';
import { CommonModule } from '@angular/common';
import {FormsModule } from '@angular/forms';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  templateUrl: './login.html',
  styleUrls: ['./login.css'],
  imports: [CommonModule, FormsModule, HttpClientModule,RouterLink]
})
export class LoginComponent {
  email: string = '';
  password: string = '';
  errorMessage: string = '';

  constructor(
    private authService: AuthService, 
    private router: Router) { }

  onLogin() {
    this.authService.login(this.email, this.password).subscribe({
      next: (res) => {
        this.errorMessage = 'Login successful!';
        console.log('Backend response:', res);
        this.router.navigate(['/dorms']);
      },
      error: (err) => {
        this.errorMessage = 'Login failed!';
        console.error('Error:', err);
      }
    });
  }
}
