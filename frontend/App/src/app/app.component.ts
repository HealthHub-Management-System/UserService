import { Component } from '@angular/core';
import { AuthService } from './auth.service';
import { AppService } from './app.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  routerOutletActivated = false;

  constructor(
    private authService: AuthService,
    private appService: AppService
  ) {}

  async ngOnInit(): Promise<void> {
    const isLoggedIn = this.authService.isLoggedIn();
    if (!isLoggedIn) {
      await this.appService.clearLoggedInUser();
    }
  }
}
