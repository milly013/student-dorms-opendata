import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DormDetail } from './dorm-detail';

describe('DormDetail', () => {
  let component: DormDetail;
  let fixture: ComponentFixture<DormDetail>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DormDetail]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DormDetail);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
