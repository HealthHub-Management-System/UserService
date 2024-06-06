import { Component, inject } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-other-functionalities-system-page',
  templateUrl: './other-functionalities-system-page.component.html',
  styleUrl: './other-functionalities-system-page.component.css',
})
export class OtherFunctionalitiesSystemPageComponent {
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
