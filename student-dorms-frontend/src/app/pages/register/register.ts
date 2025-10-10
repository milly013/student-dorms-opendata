import { Component } from '@angular/core';
import { AuthService } from '../../services/auth';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './register.html',
  styleUrl: './register.css'
})
export class Register {
 username = '';
  email = '';
  password = '';
  confirmPassword = '';
  message = '';
  error = '';
  loading = false;

  constructor(private authService: AuthService) {}

  register() {
    this.error = '';
    this.message = '';

    if (!this.username || !this.email || !this.password) {
      this.error = 'Sva polja su obavezna.';
      return;
    }

    if (this.password !== this.confirmPassword) {
      this.error = 'Lozinke se ne poklapaju.';
      return;
    }

    this.loading = true;

    this.authService.register(this.username, this.email, this.password).subscribe({
      next: () => {
        this.loading = false;
        this.message = 'Registracija uspešna! Sada se možete prijaviti.';
        this.username = '';
        this.email = '';
        this.password = '';
        this.confirmPassword = '';
      },
      error: (err) => {
        this.loading = false;
        console.error('❌ Greška pri registraciji:', err);
        this.error = err.error?.message || 'Došlo je do greške pri registraciji.';
      }
    });
  }
}