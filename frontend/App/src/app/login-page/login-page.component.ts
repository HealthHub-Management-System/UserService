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
      await this.appService.setLoggedInUserEmail(this.email);
      await this.appService.setLoggedInUserPassword(this.password);
      await this.appService.setLoggedInUserId(this.email, this.password);
      await this.appService.setLoggedInUserName(this.email, this.password);
      await this.appService.setLoggedInUserRole(this.email);
      // const user = {
      //   email: this.email,
      //   password: this.password,
      //   id: this.appService.getLoggedInUserId(),
      //   name: this.appService.getLoggedInUserName(),
      //   role: this.appService.getLoggedInUserRole(),
      // };

      // console.log('email: ', this.appService.getLoggedInUserEmail());
      // console.log('password: ', this.appService.getLoggedInUserPassword());
      // console.log('id: ', this.appService.getLoggedInUserId());
      // console.log('name: ', this.appService.getLoggedInUserName());
      // console.log('role: ', this.appService.getLoggedInUserRole());

      this.router.navigate(['/home']);
    } catch (error) {
      console.error('Login error:', error);
      this.loginError = true;
    }
  }
}
