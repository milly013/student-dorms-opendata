import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpendataFreeSpots } from './opendata-free-spots';

describe('OpendataFreeSpots', () => {
  let component: OpendataFreeSpots;
  let fixture: ComponentFixture<OpendataFreeSpots>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OpendataFreeSpots]
    })
    .compileComponents();

    fixture = TestBed.createComponent(OpendataFreeSpots);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
