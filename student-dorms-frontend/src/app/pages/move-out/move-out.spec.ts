import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MoveOut } from './move-out';

describe('MoveOut', () => {
  let component: MoveOut;
  let fixture: ComponentFixture<MoveOut>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MoveOut]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MoveOut);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
