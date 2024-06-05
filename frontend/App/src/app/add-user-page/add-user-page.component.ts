import { Component, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import { AppService } from '../app.service';
import { User } from '../User';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';

@Component({
  selector: 'app-add-user-page',
  templateUrl: './add-user-page.component.html',
  styleUrls: ['./add-user-page.component.css'],
})
export class AddUserPageComponent {
  user = {
    email: '',
    name: '',
    password: '',
    role: '',
  };

  userName: string = '';
  @ViewChild('userForm') userForm: NgForm | undefined;
  repeatedpassword: string = '';

  constructor(
    private appService: AppService,
    private router: Router,
    private authService: AuthService
  ) {
    this.userName = this.appService.getLoggedInUserName();
    if (this.checkIfLoggedIn() === false) {
      this.router.navigateByUrl('/login');
    }
  }

  submitForm(userForm: NgForm) {
    if (userForm.valid && this.user.password === this.repeatedpassword) {
      this.appService.addUser(this.user);
      window.alert('Pomyślne dodanie użytkownika!');
      userForm.resetForm();
      this.user = { name: '', email: '', role: '', password: '' };
      this.router.navigate(['/usersmanagement']);
    } else {
      window.alert('Dodanie użytkownika nie powiodło się!');
    }
  }
  checkIfLoggedIn(): boolean {
    return this.authService.isLoggedIn();
  }
}
