import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OtherFunctionalitiesSystemPageComponent } from './other-functionalities-system-page.component';
import { NavbarComponent } from '../navbar/navbar.component';

describe('OtherFunctionalitiesSystemPageComponent', () => {
  let component: OtherFunctionalitiesSystemPageComponent;
  let fixture: ComponentFixture<OtherFunctionalitiesSystemPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [OtherFunctionalitiesSystemPageComponent, NavbarComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(OtherFunctionalitiesSystemPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
