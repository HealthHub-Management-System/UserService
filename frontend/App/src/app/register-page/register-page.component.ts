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
  showMessage: boolean = false;
  isSuccess: boolean = false;

  constructor(private router: Router, private appService: AppService) {}

  submitForm(userForm: NgForm) {
    if (userForm.valid) {
      this.appService.addUserRegistration(this.user);
      console.log('Registry User:', this.user);
      window.alert('Pomyślna rejestracja!');
      userForm.resetForm();
      this.user = { name: '', email: '', role: 'patient', password: '' };
      this.router.navigate(['/login']);
    } else {
      window.alert('Rejestracja nie powiodła się!');
    }
  }
}
