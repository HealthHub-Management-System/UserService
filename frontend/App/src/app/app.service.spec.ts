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
});
