import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MultiStepRegisterComponent } from './signup';

describe('Signup', () => {
  let component: MultiStepRegisterComponent;
  let fixture: ComponentFixture<MultiStepRegisterComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MultiStepRegisterComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MultiStepRegisterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
