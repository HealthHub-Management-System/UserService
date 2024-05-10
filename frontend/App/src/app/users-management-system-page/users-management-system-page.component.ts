import { Component, OnInit } from '@angular/core';
import { AppService } from '../app.service';
import { User } from '../User';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';
import { lastValueFrom } from 'rxjs';

@Component({
  selector: 'app-users-management-system-page',
  templateUrl: './users-management-system-page.component.html',
  styleUrls: ['./users-management-system-page.component.css'],
})
export class UsersManagementSystemPageComponent implements OnInit {
  users: User[] = [];
  filteredUsers: User[] = [];
  username: string = '';
  searchTerm: string = '';

  constructor(
    private appService: AppService,
    private authService: AuthService,
    private router: Router
  ) {
    this.username = this.appService.getLoggedInUserName();
    if (this.checkIfLoggedIn() === false) {
      this.router.navigateByUrl('/login');
    }
  }

  ngOnInit(): void {
    this.loadUsers();
  }
  async loadUsers(): Promise<void> {
    try {
      const users = await this.appService.getUsers();
      this.users = users;
      this.applySearchFilter();
    } catch (error) {
      console.error('Error retrieving users:', error);
    }
  }

  applySearchFilter(): void {
    if (this.searchTerm.trim() === '') {
      this.filteredUsers = this.users;
    } else {
      this.filteredUsers = this.users.filter((user) =>
        user.name.toLowerCase().includes(this.searchTerm.trim().toLowerCase())
      );
    }
  }

  deleteUser(userId: string): void {
    this.appService.deleteUser(userId).subscribe(() => {
      this.users = this.users.filter((user) => user.id !== userId);
      this.applySearchFilter();
    });
  }
  confirmDeleteUser(userId: string): void {
    const confirmation = confirm('Czy na pewno chcesz usunąć użytkownika?');
    if (confirmation) {
      this.deleteUser(userId);
    }
  }
  checkIfLoggedIn(): boolean {
    return this.authService.isLoggedIn();
  }
}
