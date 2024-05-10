import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, lastValueFrom, map, switchMap, throwError } from 'rxjs';

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

  constructor(private http: HttpClient) {}

  async updateUsersList(): Promise<void> {}

  getUserById(userId: string): void {
    const url = `${this.apiUrl}/${userId}`;

    this.http.get<any>(url, { withCredentials: true }).subscribe(
      (response: any) => {
        console.log(response);
      },
      (error) => {
        console.log(error);
      }
    );
  }

  async getUsers(): Promise<User[]> {
    try {
      const response = await lastValueFrom(
        this.http.get<{ users: User[] }>(`${this.apiUrl}?limit=100`, {
          withCredentials: true,
        })
      );
      const users = response.users;
      return users;
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
      (response: any) => {
        console.log(response);
      },
      (error) => {
        console.log(error);
      }
    );
  }
  addUserRegistration(user: any): void {
    this.http.post<any>(this.apiUrl, user).subscribe(
      (response: any) => {
        console.log(response);
        console.log('Registry User-service:', user);
      },
      (error) => {
        console.log(error);
      }
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
      (response: any) => {
        console.log(response);
      },
      (error) => {
        console.log(error);
      }
    );
    await this.setLoggedInUserEmail(email);
    await this.setLoggedInUserName(email, this.getLoggedInUserPassword());
    await this.updateUsersList();
  }
  async setLoggedInUserEmail(email: string): Promise<void> {
    this.loggedInUser.loggedInUserEmail = email;
  }
  getLoggedInUserEmail(): string {
    return this.loggedInUser.loggedInUserEmail;
  }
  async setLoggedInUserPassword(password: string): Promise<void> {
    this.loggedInUser.loggedInUserPassword = password;
  }
  getLoggedInUserPassword(): string {
    return this.loggedInUser.loggedInUserPassword;
  }
  async setLoggedInUserName(email: string, password: string): Promise<void> {
    try {
      const users = await this.getUsers();
      const loggedInUser = users.find((user) => user.email === email);
      if (loggedInUser) {
        this.loggedInUser.loggedInUserName = loggedInUser.name;
      }
    } catch (error) {
      console.error('Error retrieving users:', error);
    }
  }
  getLoggedInUserName(): string {
    return this.loggedInUser.loggedInUserName;
  }

  async setLoggedInUserId(email: string, password: string): Promise<void> {
    try {
      const users = await this.getUsers();
      const loggedInUser = users.find((user) => user.email === email);
      if (loggedInUser) {
        this.loggedInUser.loggedInUserid = loggedInUser.id;
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
      const users = await this.getUsers();
      const loggedInUser = users.find((user) => user.email === email);
      if (loggedInUser) {
        this.loggedInUser.loggedInUserRole = loggedInUser.role;
      }
    } catch (error) {
      console.error('Error retrieving users:', error);
    }
  }
  getLoggedInUserRole(): string {
    return this.loggedInUser.loggedInUserRole;
  }
}
