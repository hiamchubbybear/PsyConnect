import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Folders } from './folders';

describe('Folders', () => {
  let component: Folders;
  let fixture: ComponentFixture<Folders>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Folders]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Folders);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
