import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PaymentSystemPageComponent } from './payment-system-page.component';
import { NavbarComponent } from '../navbar/navbar.component';

describe('PaymentSystemPageComponent', () => {
  let component: PaymentSystemPageComponent;
  let fixture: ComponentFixture<PaymentSystemPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [PaymentSystemPageComponent, NavbarComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(PaymentSystemPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
