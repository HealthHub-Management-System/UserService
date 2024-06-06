import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginPageComponent } from './login-page/login-page.component';
import { RegisterPageComponent } from './register-page/register-page.component';
import { HomePageComponent } from './home-page/home-page.component';
import { NavbarComponent } from './navbar/navbar.component';
import { PaymentSystemPageComponent } from './payment-system-page/payment-system-page.component';
import { MedicalDocumentatioSystemPageComponent } from './medical-documentation-system-page/medical-documentatio-system-page.component';
import { AppointmentSystemPageComponent } from './appointment-system-page/appointment-system-page.component';
import { OpinionSystemPageComponent } from './opinion-system-page/opinion-system-page.component';
import { OtherFunctionalitiesSystemPageComponent } from './other-functionalities-system-page/other-functionalities-system-page.component';
import { ProfilePageComponent } from './profile-page/profile-page.component';
import { UsersManagementSystemPageComponent } from './users-management-system-page/users-management-system-page.component';
import { AddUserPageComponent } from './add-user-page/add-user-page.component';
import { HttpClientModule } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';
import { AppService } from './app.service';
import { AuthService } from './auth.service';

@NgModule({
  declarations: [
    AppComponent,
    LoginPageComponent,
    RegisterPageComponent,
    HomePageComponent,
    NavbarComponent,
    PaymentSystemPageComponent,
    MedicalDocumentatioSystemPageComponent,
    AppointmentSystemPageComponent,
    OpinionSystemPageComponent,
    OtherFunctionalitiesSystemPageComponent,
    ProfilePageComponent,
    UsersManagementSystemPageComponent,
    AddUserPageComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    CommonModule,
    AppRoutingModule,
    HttpClientModule,
  ],
  providers: [CookieService, AppService, AuthService],
  bootstrap: [AppComponent],
})
export class AppModule {}
