import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpendataHub } from './opendata-hub';

describe('OpendataHub', () => {
  let component: OpendataHub;
  let fixture: ComponentFixture<OpendataHub>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OpendataHub]
    })
    .compileComponents();

    fixture = TestBed.createComponent(OpendataHub);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
