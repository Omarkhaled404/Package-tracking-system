import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { RegisterPostData } from '../interfaces/auth';
//import { LoginPostData } from '../interfaces/auth';
import { User } from '../interfaces/auth';
@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private baseURL="http://localhost:8000";

  constructor(private http: HttpClient) { }
  registerUser(postData: RegisterPostData){
    return this.http.post(`${this.baseURL}/register`, postData);
  }
  getUserDetails(postData: { email: string; password: string }): Observable<User> {
    return this.http.post<User>(`${this.baseURL}/login`, postData);
}
}
