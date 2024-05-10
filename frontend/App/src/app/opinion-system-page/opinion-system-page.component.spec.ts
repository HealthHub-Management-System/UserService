import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpinionSystemPageComponent } from './opinion-system-page.component';
import { NavbarComponent } from '../navbar/navbar.component';

describe('OpinionSystemPageComponent', () => {
  let component: OpinionSystemPageComponent;
  let fixture: ComponentFixture<OpinionSystemPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [OpinionSystemPageComponent, NavbarComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(OpinionSystemPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
