import { Component } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-navbar',
  imports: [RouterLink, CommonModule],
  templateUrl: './navbar.html',
  styleUrl: './navbar.css'
})
export class Navbar {
  constructor(
    public authService: AuthService,
    private router: Router
  ) {}
  
   logout(): void {
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  goHome(): void {
    if (this.authService.isLoggedIn()) {
      this.router.navigate(['/dorms']);
    } else {
      this.router.navigate(['/login']);
    }
  }

}
