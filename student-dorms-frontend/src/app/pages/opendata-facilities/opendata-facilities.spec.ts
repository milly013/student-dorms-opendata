import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpendataFacilities } from './opendata-facilities';

describe('OpendataFacilities', () => {
  let component: OpendataFacilities;
  let fixture: ComponentFixture<OpendataFacilities>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OpendataFacilities]
    })
    .compileComponents();

    fixture = TestBed.createComponent(OpendataFacilities);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
