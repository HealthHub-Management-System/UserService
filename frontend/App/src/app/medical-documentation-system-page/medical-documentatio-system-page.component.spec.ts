import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MedicalDocumentatioSystemPageComponent } from './medical-documentatio-system-page.component';
import { NavbarComponent } from '../navbar/navbar.component';

describe('MedicalDocumentatioSystemPageComponent', () => {
  let component: MedicalDocumentatioSystemPageComponent;
  let fixture: ComponentFixture<MedicalDocumentatioSystemPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [MedicalDocumentatioSystemPageComponent, NavbarComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MedicalDocumentatioSystemPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
