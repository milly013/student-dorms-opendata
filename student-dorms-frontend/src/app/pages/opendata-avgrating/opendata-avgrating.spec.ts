import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpendataAvgrating } from './opendata-avgrating';

describe('OpendataAvgrating', () => {
  let component: OpendataAvgrating;
  let fixture: ComponentFixture<OpendataAvgrating>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OpendataAvgrating]
    })
    .compileComponents();

    fixture = TestBed.createComponent(OpendataAvgrating);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
