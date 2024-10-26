import { Component, inject } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { InputTextModule } from 'primeng/inputtext';
import { Password, PasswordModule } from 'primeng/password';
import { passwordMismatchValidator } from '../../Shared/password-mismatch.directive';
import { AuthService } from '../../services/auth.service';
import { RegisterPostData } from '../../interfaces/auth';
//import { response, Router } from 'express';
import { TreeSelectModule } from 'primeng/treeselect';
import { DropdownModule } from 'primeng/dropdown';
import { MessageService } from 'primeng/api';
@Component({
  selector: 'app-register',
  standalone: true,
  imports: [ReactiveFormsModule, CardModule, InputTextModule, PasswordModule, ButtonModule, RouterLink, DropdownModule, TreeSelectModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {
 
  private router = inject(Router);
  private messageService = inject(MessageService);
  private registerService = inject(AuthService);
  registerForm =new FormGroup({
    name: new FormControl('', [Validators.required]),
    email: new FormControl('', [Validators.required,Validators.pattern(/[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/)]),
    //phone: new FormControl('', [Validators.required]),
    phone: new FormControl('', [Validators.required,  Validators.pattern(/^[0-9]{11}$/)]),
    password: new FormControl('', [Validators.required]),
    confirmPassword: new FormControl('', [Validators.required]),
    role: new FormControl('', [Validators.required])
  },
  {
    validators:  passwordMismatchValidator
  }
)

onRegister(){
  const postData = {...this.registerForm.value};
  this.registerService.registerUser(postData as RegisterPostData).subscribe({
    next: (response) => {
      this.messageService.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Registered successfully'
      });
      this.router.navigate(['login']);
      console.log('Response:', response); // Log to verify response structure
    },
    error: (err) => {
      this.messageService.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Something went wrong'
      });
      console.error('Error:', err);
    }
  });
}

  get fullName(){
     return this.registerForm.controls['name'];
  }
  get email() {
    return this.registerForm.controls['email'];
  }
  get phone() {
    return this.registerForm.get('phone');
  }


  get password() {
    return this.registerForm.controls['password'];
  }

  get confirmPassword() {
    return this.registerForm.controls['confirmPassword'];
  }
  get role() {
    return this.registerForm.get('role');
  }
}
