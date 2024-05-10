import { Component, inject } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-opinion-system-page',
  templateUrl: './opinion-system-page.component.html',
  styleUrl: './opinion-system-page.component.css',
})
export class OpinionSystemPageComponent {
  username: string = 'Jan Kowalski';

  private sanitizer = inject(DomSanitizer);

  trustedURL: any = '';

  constructor(private authService: AuthService, private router: Router) {
    if (this.checkIfLoggedIn() === false) {
      this.router.navigateByUrl('/login');
    }
    this.trustedURL = this.sanitizer.bypassSecurityTrustResourceUrl(
      'https://www.bing.com/'
    );
  }
  checkIfLoggedIn(): boolean {
    return this.authService.isLoggedIn();
  }
}
