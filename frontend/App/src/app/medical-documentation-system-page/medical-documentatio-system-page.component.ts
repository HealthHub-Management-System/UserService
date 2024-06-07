import { Component, inject } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-medical-documentatio-system-page',
  templateUrl: './medical-documentatio-system-page.component.html',
  styleUrl: './medical-documentatio-system-page.component.css',
})
export class MedicalDocumentatioSystemPageComponent {
  private sanitizer = inject(DomSanitizer);

  trustedURL: any = '';

  constructor(private authService: AuthService, private router: Router) {
    if (this.checkIfLoggedIn() === false) {
      this.router.navigateByUrl('/login');
    }
    this.trustedURL = this.sanitizer.bypassSecurityTrustResourceUrl(
      'http://localhost:2137'
    );
  }
  checkIfLoggedIn(): boolean {
    return this.authService.isLoggedIn();
  }
}
