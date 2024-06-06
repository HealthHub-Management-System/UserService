import { Component } from '@angular/core';
import { NgForm } from '@angular/forms';
import { AppService } from '../app.service';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-profile-page',
  templateUrl: './profile-page.component.html',
  styleUrl: './profile-page.component.css',
})
export class ProfilePageComponent {
  isMenuOpen: boolean = false;

  user = {
    email: '',
    name: '',
    password: '',
    role: '',
  };

  constructor(
    private appService: AppService,
    private router: Router,
    private authService: AuthService
  ) {
    this.user.email = this.appService.getLoggedInUserEmail();
    this.user.name = this.appService.getLoggedInUserName();
    if (this.checkIfLoggedIn() === false) {
      this.router.navigateByUrl('/login');
    }
  }

  submitForm(userForm: NgForm) {
    if (userForm.valid) {
      this.confirmUpdateUser();
    }
  }

  toggleMenu() {
    this.isMenuOpen = !this.isMenuOpen;
  }

  confirmUpdateUser(): void {
    const confirmation = confirm('Czy na pewno chcesz zmienić dane?');
    this.appService.updateUser(this.user.email, this.user.name);
    this.router.navigateByUrl('/home');
  }

  checkIfLoggedIn(): boolean {
    return this.authService.isLoggedIn();
  }
}
