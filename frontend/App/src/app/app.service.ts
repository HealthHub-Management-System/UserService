import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, BehaviorSubject, lastValueFrom } from 'rxjs';
import { User } from './User';

@Injectable({
  providedIn: 'root',
})
export class AppService {
  private apiUrl = 'http://localhost:8080/api/v1/users';
  public users: User[] = [];
  private loggedInUser = {
    loggedInUserEmail: '',
    loggedInUserPassword: '',
    loggedInUserName: '',
    loggedInUserid: '',
    loggedInUserRole: '',
  };
  private loggedInUserSubject = new BehaviorSubject(this.loggedInUser);

  loggedInUser$ = this.loggedInUserSubject.asObservable();

  constructor(private http: HttpClient) {
    this.loadLoggedInUser();
  }

  getUserById(userId: string): void {
    const url = `${this.apiUrl}/${userId}`;

    this.http.get<any>(url, { withCredentials: true }).subscribe(
      (response: any) => {},
      (error) => {}
    );
  }

  async getUsers(
    page: number = 1,
    limit: number = 5
  ): Promise<{ users: User[]; total: number }> {
    try {
      const response = await lastValueFrom(
        this.http.get<{ users: User[]; total: number }>(
          `${this.apiUrl}?page=${page}&limit=${limit}`,
          {
            withCredentials: true,
          }
        )
      );
      return response;
    } catch (error) {
      console.error('Error retrieving users:', error);
      throw error;
    }
  }

  deleteUser(userId: string): Observable<any> {
    const url = `${this.apiUrl}/${userId}`;
    return this.http.delete<any>(url, { withCredentials: true });
  }

  addUser(user: any): void {
    this.http.post<any>(this.apiUrl, user, { withCredentials: true }).subscribe(
      (response: any) => {},
      (error) => {}
    );
  }

  async addUserRegistration(user: any): Promise<void> {
    this.http.post<any>(this.apiUrl, user).subscribe(
      (response: any) => {},
      (error) => {}
    );
  }

  async updateUser(email: string, name: string): Promise<void> {
    const url = `${this.apiUrl}/${this.getLoggedInUserId()}`;
    let user = {
      name: name,
      email: email,
      password: this.getLoggedInUserPassword(),
      role: this.getLoggedInUserRole(),
    };
    this.http.put<any>(url, user, { withCredentials: true }).subscribe(
      async (response: any) => {
        await this.setLoggedInUserEmail(email);
        await this.setLoggedInUserName(name);
        this.saveLoggedInUser();
      },
      (error) => {}
    );
  }

  async getAllUsers(): Promise<User[]> {
    let page = 1;
    const limit = 100;
    let allUsers: User[] = [];
    let totalUsers = 0;

    do {
      const { users, total } = await this.getUsers(page, limit);
      allUsers = allUsers.concat(users);
      totalUsers = total;
      page++;
    } while (allUsers.length < totalUsers);

    return allUsers;
  }

  saveLoggedInUser(): void {
    localStorage.setItem('loggedInUser', JSON.stringify(this.loggedInUser));
    this.loggedInUserSubject.next(this.loggedInUser);
  }

  loadLoggedInUser(): void {
    const userData = localStorage.getItem('loggedInUser');
    if (userData) {
      this.loggedInUser = JSON.parse(userData);
      this.loggedInUserSubject.next(this.loggedInUser);
    }
  }

  async setLoggedInUserEmail(email: string): Promise<void> {
    this.loggedInUser.loggedInUserEmail = email;
    this.saveLoggedInUser();
  }

  getLoggedInUserEmail(): string {
    return this.loggedInUser.loggedInUserEmail;
  }

  async setLoggedInUserPassword(password: string): Promise<void> {
    this.loggedInUser.loggedInUserPassword = password;
    this.saveLoggedInUser();
  }

  getLoggedInUserPassword(): string {
    return this.loggedInUser.loggedInUserPassword;
  }

  async setLoggedInUserName(name: string): Promise<void> {
    try {
      this.loggedInUser.loggedInUserName = name;
    } catch (error) {
      console.error('Error retrieving users:', error);
    }
  }

  getLoggedInUserName(): string {
    return this.loggedInUser.loggedInUserName;
  }

  async setLoggedInUserId(email: string, password: string): Promise<void> {
    try {
      const users = await this.getAllUsers();
      const loggedInUser = users.find((user) => user.email === email);
      if (loggedInUser) {
        this.loggedInUser.loggedInUserid = loggedInUser.id;
        this.saveLoggedInUser();
      }
    } catch (error) {
      console.error('Error retrieving users:', error);
    }
  }

  getLoggedInUserId(): string {
    return this.loggedInUser.loggedInUserid;
  }

  async setLoggedInUserRole(email: string): Promise<void> {
    try {
      const users = await this.getAllUsers();
      const loggedInUser = users.find((user) => user.email === email);
      if (loggedInUser) {
        this.loggedInUser.loggedInUserRole = loggedInUser.role;
        this.saveLoggedInUser();
      }
    } catch (error) {
      console.error('Error retrieving users:', error);
    }
  }

  getLoggedInUserRole(): string {
    return this.loggedInUser.loggedInUserRole;
  }

  async clearLoggedInUser(): Promise<void> {
    this.loggedInUser = {
      loggedInUserEmail: '',
      loggedInUserPassword: '',
      loggedInUserName: '',
      loggedInUserid: '',
      loggedInUserRole: '',
    };
    localStorage.removeItem('loggedInUser');
    this.loggedInUserSubject.next(this.loggedInUser);
  }

  async getCurrentUser(): Promise<User | null> {
    try {
      const user = await lastValueFrom(
        this.http.get<User>(`${this.apiUrl}/current`, { withCredentials: true })
      );
      this.loggedInUser.loggedInUserName = user.name;
      this.loggedInUser.loggedInUserEmail = user.email;
      this.loggedInUser.loggedInUserid = user.id;
      this.loggedInUser.loggedInUserRole = user.role;
      this.saveLoggedInUser();
      return user;
    } catch (error) {
      console.error('Error retrieving current user:', error);
      return null;
    }
  }
}
