import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AppointmentSystemPageComponent } from './appointment-system-page.component';
import { NavbarComponent } from '../navbar/navbar.component';

describe('AppointmentSystemPageComponent', () => {
  let component: AppointmentSystemPageComponent;
  let fixture: ComponentFixture<AppointmentSystemPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [AppointmentSystemPageComponent, NavbarComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(AppointmentSystemPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
