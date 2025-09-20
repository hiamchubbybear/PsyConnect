import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Refunds } from './refunds';

describe('Refunds', () => {
  let component: Refunds;
  let fixture: ComponentFixture<Refunds>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Refunds]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Refunds);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
