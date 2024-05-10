import { Injectable } from '@angular/core';
import { AppService } from './app.service';
import { User } from './User';
import { CookieService } from 'ngx-cookie-service';
import { Router } from '@angular/router';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  // private users: User[] = [];

  constructor(
    private userService: AppService,
    private http: HttpClient,
    private cookiesService: CookieService
  ) {}

  private apiUrl = 'http://localhost:8080/api/v1/users';

  login(email: string, password: string): Observable<HttpResponse<any>> {
    return this.http.post<any>(
      `${this.apiUrl}/login`,
      { email: email, password: password },
      { observe: 'response', withCredentials: true }
    );
  }
  logout(): void {
    this.http
      .post<any>(`${this.apiUrl}/logout`, {}, { withCredentials: true })
      .subscribe(
        (response) => {
          console.log('Logout:', response);
        },
        (error) => {
          console.error('Logout error:', error);
        }
      );
  }
  isLoggedIn(): boolean {
    return this.cookiesService.check('session');
  }
}
