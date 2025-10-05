import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PopularDorms } from './popular-dorms';

describe('PopularDorms', () => {
  let component: PopularDorms;
  let fixture: ComponentFixture<PopularDorms>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PopularDorms]
    })
    .compileComponents();

    fixture = TestBed.createComponent(PopularDorms);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
