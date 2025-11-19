import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Segments } from './segments';

describe('Segments', () => {
  let component: Segments;
  let fixture: ComponentFixture<Segments>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Segments]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Segments);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
