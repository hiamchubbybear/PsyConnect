import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Sessions } from './sessions';

describe('Sessions', () => {
  let component: Sessions;
  let fixture: ComponentFixture<Sessions>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Sessions]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Sessions);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
