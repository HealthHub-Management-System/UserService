// import { TestBed } from '@angular/core/testing';

// import { AppService } from './app.service';

// describe('AppService', () => {
//   let service: AppService;

//   beforeEach(() => {
//     TestBed.configureTestingModule({});
//     service = TestBed.inject(AppService);
//   });

//   it('should be created', () => {
//     expect(service).toBeTruthy();
//   });
// });
import { TestBed } from '@angular/core/testing';
import {
  HttpClientTestingModule,
  HttpTestingController,
} from '@angular/common/http/testing';
import { AppService } from './app.service';
import { User } from './User';

describe('AppService', () => {
  let service: AppService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [AppService],
    });
    service = TestBed.inject(AppService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  // it('should get users', () => {
  //   const users: User[] = [
  //     {
  //       id: '1',
  //       name: 'John Doe',
  //       email: 'johndoe@gmail.com',
  //       password: 'JohnDoe123.KK',
  //       role: 'patient',
  //     },
  //     {
  //       id: '2',
  //       name: 'Jane Doe',
  //       email: 'johndoe1233@gmail.com',
  //       password: 'JohnDoe12&3.KK',
  //       role: 'doctor',
  //     },
  //   ];

  //   service.getUsers().subscribe((response) => {
  //     // expect(response).toEqual(users);
  //     console.log(response);
  //   });

  //   const req = httpMock.expectOne('http://localhost:8080/api/v1/users');
  //   expect(req.request.method).toBe('GET');
  //   req.flush(users);
  // });

  it('should add a user', () => {
    const newUser = {
      name: 'John Doe',
      email: 'johndoe@gmail.com',
      password: 'JohnDoe123.KK',
      role: 'patient',
    };

    service.addUser(newUser);

    const req = httpMock.expectOne('http://localhost:8080/api/v1/users');
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(newUser);
    req.flush({});
  });

  it('should delete a user', () => {
    const userId = '1';

    service.deleteUser(userId).subscribe();

    const req = httpMock.expectOne(
      `http://localhost:8080/api/v1/users/${userId}`
    );
    expect(req.request.method).toBe('DELETE');
  });

  // it('should get a user by ID', () => {
  //   const userId = '1';
  //   const user: User = {
  //     id: userId,
  //     name: 'John Doe',
  //     email: 'johndoe@gmail.com',
  //     password: 'JohnDoe123.KK',
  //     role: 'patient',
  //   };

  //   service.getUserById(userId);

  //   const req = httpMock.expectOne(
  //     `http://localhost:8080/api/v1/users/${userId}`
  //   );
  //   expect(req.request.method).toBe('GET');
  //   req.flush(user);
  // });

  // it('should update a user', () => {
  //   const updatedUser: User = {
  //     id: '1',
  //     name: 'Updated User',
  //     email: 'updated@example.com',
  //     password: 'Updated123!KK',
  //     role: 'patient',
  //   };

  //   service.updateUser(updatedUser);

  //   const req = httpMock.expectOne(
  //     `http://localhost:8080/api/v1/users/${updatedUser.id}`
  //   );
  //   expect(req.request.method).toBe('PUT');
  //   expect(req.request.body).toEqual(updatedUser);
  //   req.flush({});
  // });
});
