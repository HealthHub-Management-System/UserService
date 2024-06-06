import { Component, OnInit } from '@angular/core';
import { AppService } from '../app.service';
import { User } from '../User';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-users-management-system-page',
  templateUrl: './users-management-system-page.component.html',
  styleUrls: ['./users-management-system-page.component.css'],
})
export class UsersManagementSystemPageComponent implements OnInit {
  users: User[] = [];
  filteredUsers: User[] = [];
  paginatedUsers: User[] = [];
  allUsers: User[] = [];
  username: string = '';
  searchTerm: string = '';

  currentPage: number = 1;
  itemsPerPage: number = 5;
  totalPages: number = 1;
  totalUsers: number = 0;

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
      const response = await this.appService.getUsers(
        this.currentPage,
        this.itemsPerPage
      );
      this.users = response.users;
      this.totalUsers = response.total;
      this.applySearchFilter();
      this.allUsers = await this.appService.getAllUsers();
    } catch (error) {
      console.error('Error retrieving users:', error);
    }
  }

  applySearchFilter(): void {
    if (this.searchTerm.trim() === '') {
      this.filteredUsers = this.users;
      this.updatePaginatedUsers();
    } else {
      this.filteredUsers = this.allUsers.filter((user) =>
        user.name.toLowerCase().includes(this.searchTerm.trim().toLowerCase())
      );
      this.updatePaginatedUsers2();
    }
  }

  async updatePaginatedUsers(): Promise<void> {
    this.totalPages = Math.ceil(this.totalUsers / this.itemsPerPage);
    this.paginatedUsers = this.filteredUsers.slice(0, this.itemsPerPage);
  }
  async updatePaginatedUsers2(): Promise<void> {
    this.totalPages = Math.ceil(this.filteredUsers.length / this.itemsPerPage);
    this.paginatedUsers = this.filteredUsers.slice(0, this.itemsPerPage);
  }

  nextPage(): void {
    if (this.currentPage < this.totalPages) {
      this.currentPage++;
      this.loadUsers();
    }
  }

  previousPage(): void {
    if (this.currentPage > 1) {
      this.currentPage--;
      this.loadUsers();
    }
  }

  goToPage(event: Event): void {
    const target = event.target as HTMLInputElement;
    const pageNumber = Number(target.value);

    if (pageNumber >= 1 && pageNumber <= this.totalPages) {
      this.currentPage = pageNumber;
      this.loadUsers();
    }
  }

  deleteUser(userId: string): void {
    this.appService.deleteUser(userId).subscribe(() => {
      this.loadUsers();
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

  mapRole(role: string): string {
    switch (role) {
      case 'patient':
        return 'pacjent';
      case 'doctor':
        return 'lekarz';
      case 'admin':
        return 'administrator';
      default:
        return role;
    }
  }
}
