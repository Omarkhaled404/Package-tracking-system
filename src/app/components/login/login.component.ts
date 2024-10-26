import { Component } from '@angular/core';
import { CardModule } from 'primeng/card';
import { InputTextModule } from 'primeng/inputtext';
import { FormsModule } from '@angular/forms'; 
import { PasswordModule } from 'primeng/password';
import { ButtonModule } from 'primeng/button';
import { RouterLink, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service'; // Import AuthService for login
import { MessageService } from 'primeng/api'; // Assuming you use PrimeNG's MessageService

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CardModule, InputTextModule, FormsModule, PasswordModule, ButtonModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  login = {
    email: '',
    password: '',
  }

  constructor(
    private authService: AuthService,
    private messageService: MessageService,
    private router: Router
  ) {}

  onLogin() {
    console.log('Logging in with:', this.login);
  
    this.authService.getUserDetails(this.login).subscribe({
      next: (response) => {
        this.messageService.add({ severity: 'success', summary: 'Success', detail: 'Login successful' });
        this.router.navigate(['home']);
        console.log('Login response:', response);
      },
      error: (err) => {
        console.error('Login error:', err); // Log full error object for better debugging
        if (err.error) {
          console.error('Error message from backend:', err.error);
        }
        this.messageService.add({ severity: 'error', summary: 'Error', detail: 'Invalid email or password' });
      }
    });
  }
  
}
