import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AppComponent } from './app.component';
import { LoginPageComponent } from './login-page/login-page.component';
import { RegisterPageComponent } from './register-page/register-page.component';
import { HomePageComponent } from './home-page/home-page.component';
import { PaymentSystemPageComponent } from './payment-system-page/payment-system-page.component';
import { MedicalDocumentatioSystemPageComponent } from './medical-documentation-system-page/medical-documentatio-system-page.component';
import { AppointmentSystemPageComponent } from './appointment-system-page/appointment-system-page.component';
import { OpinionSystemPageComponent } from './opinion-system-page/opinion-system-page.component';
import { OtherFunctionalitiesSystemPageComponent } from './other-functionalities-system-page/other-functionalities-system-page.component';
import { ProfilePageComponent } from './profile-page/profile-page.component';
import { UsersManagementSystemPageComponent } from './users-management-system-page/users-management-system-page.component';
import { AddUserPageComponent } from './add-user-page/add-user-page.component';

const routes: Routes = [
  { path: '', component: AppComponent },
  { path: 'login', component: LoginPageComponent },
  { path: 'register', component: RegisterPageComponent },
  { path: 'home', component: HomePageComponent },
  { path: 'payment', component: PaymentSystemPageComponent },
  {
    path: 'medicaldocumentation',
    component: MedicalDocumentatioSystemPageComponent,
  },
  { path: 'appointment', component: AppointmentSystemPageComponent },
  { path: 'opinion', component: OpinionSystemPageComponent },
  {
    path: 'otherfunctionalities',
    component: OtherFunctionalitiesSystemPageComponent,
  },
  { path: 'usersmanagement', component: UsersManagementSystemPageComponent },
  { path: 'addUser', component: AddUserPageComponent },
  { path: 'profile', component: ProfilePageComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
