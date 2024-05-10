import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddUserPageComponent } from './add-user-page.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { NavbarComponent } from '../navbar/navbar.component';
import { FormsModule } from '@angular/forms';

describe('AddUserPageComponent', () => {
  let component: AddUserPageComponent;
  let fixture: ComponentFixture<AddUserPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AddUserPageComponent, NavbarComponent],
      imports: [HttpClientTestingModule, FormsModule],
    }).compileComponents();

    fixture = TestBed.createComponent(AddUserPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
