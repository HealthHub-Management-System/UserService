import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';
import { AppService } from '../app.service';

@Component({
  selector: 'app-home-page',
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.css',
})
export class HomePageComponent {
  user = {
    name: '',
    role: '',
  };

  constructor(
    private router: Router,
    private authService: AuthService,
    private appService: AppService
  ) {
    this.user.name = this.appService.getLoggedInUserName();
    this.user.role = this.appService.getLoggedInUserRole();
    if (this.checkIfLoggedIn() === false) {
      this.router.navigateByUrl('/login');
    }
  }

  checkIfLoggedIn(): boolean {
    console.log('App-service');
    console.log('email: ', this.appService.getLoggedInUserEmail());
    console.log('password: ', this.appService.getLoggedInUserPassword());
    console.log('id: ', this.appService.getLoggedInUserId());
    console.log('name: ', this.appService.getLoggedInUserName());
    console.log('role: ', this.appService.getLoggedInUserRole());
    return this.authService.isLoggedIn();
  }

  navigateTo(url: string): void {
    this.router.navigateByUrl(url);
  }
}
