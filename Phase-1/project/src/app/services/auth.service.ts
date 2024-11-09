import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { RegisterPostData, User } from '../interfaces/auth';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private user_id: number | null = null;
  private baseURL = "http://localhost:8000";

  constructor(private http: HttpClient) { }

  registerUser(postData: RegisterPostData) {
    return this.http.post(`${this.baseURL}/register`, postData);
  }

  getUserDetails(postData: { email: string; password: string }): Observable<User> {
    return this.http.post<User>(`${this.baseURL}/login`, postData);
  }

  setUserId(id: number) {
    this.user_id = id;
    console.log('Setting userId:', id); // Log userId being set
    localStorage.setItem('user_id', id.toString());
  }
  

  public getUserId(): number | null {
    if (this.user_id === null && typeof window !== 'undefined') {
      const storedId = localStorage.getItem('user_id');
      this.user_id = storedId ? Number(storedId) : null;
    }
    return this.user_id;
  }

  clearUserId() {
    this.user_id = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('userId');
    }
  }
}
