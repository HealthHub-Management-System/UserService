import { Component, Input, OnInit } from '@angular/core';
import { AuthService } from '../auth.service';
import { AppService } from '../app.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css'],
})
export class NavbarComponent implements OnInit {
  @Input() username: string | undefined;

  constructor(
    private authService: AuthService,
    private appService: AppService
  ) {}

  isMenuOpen: boolean = false;

  ngOnInit(): void {
    this.appService.loggedInUser$.subscribe((user) => {
      this.username = user.loggedInUserName;
    });
  }

  toggleMenu(): void {
    this.isMenuOpen = !this.isMenuOpen;
  }

  closeMenu(): void {
    this.isMenuOpen = false;
  }

  logout(): void {
    this.appService.setLoggedInUserEmail('');
    this.appService.setLoggedInUserPassword('');
    this.authService.logout();
  }
}
