import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Chatpage } from './chatpage';

describe('Chatpage', () => {
  let component: Chatpage;
  let fixture: ComponentFixture<Chatpage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Chatpage]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Chatpage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
