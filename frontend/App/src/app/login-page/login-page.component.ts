import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';
import { CookieService } from 'ngx-cookie-service';
import { AppService } from '../app.service';

@Component({
  selector: 'app-login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.css'],
})
export class LoginPageComponent {
  email: string = '';
  password: string = '';
  loginError: boolean = false;

  constructor(
    private authService: AuthService,
    private cookieService: CookieService,
    private router: Router,
    private appService: AppService
  ) {}

  async login(): Promise<void> {
    try {
      const response = await this.authService
        .login(this.email, this.password)
        .toPromise();

      await Promise.all([
        this.appService.setLoggedInUserEmail(this.email),
        this.appService.setLoggedInUserPassword(this.password),
        this.appService.setLoggedInUserId(this.email, this.password),
        this.appService.setLoggedInUserName(this.email, this.password),
        this.appService.setLoggedInUserRole(this.email),
      ]);

      this.router.navigate(['/home']);
    } catch (error) {
      console.error('Login error:', error);
      this.loginError = true;
    }
  }
}
