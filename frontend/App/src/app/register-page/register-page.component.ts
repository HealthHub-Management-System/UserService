import { Component } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router } from '@angular/router';
import { AppService } from '../app.service';

@Component({
  selector: 'app-register-page',
  templateUrl: './register-page.component.html',
  styleUrls: ['./register-page.component.css'],
})
export class RegisterPageComponent {
  user = {
    name: '',
    email: '',
    password: '',
    role: 'patient',
  };
  repeatedpassword: string = '';
  formSubmitted = false;

  constructor(private router: Router, private appService: AppService) {}

  submitForm(userForm: NgForm) {
    if (userForm.valid) {
      this.appService.addUser(this.user);
      userForm.resetForm();
      this.user = { name: '', email: '', role: 'patient', password: '' };
      this.router.navigate(['/login']);
    }
  }
}
