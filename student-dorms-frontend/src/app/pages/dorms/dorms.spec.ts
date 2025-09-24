import { ComponentFixture, TestBed } from '@angular/core/testing';

import { Dorms } from './dorms';

describe('Dorms', () => {
  let component: Dorms;
  let fixture: ComponentFixture<Dorms>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [Dorms]
    })
    .compileComponents();

    fixture = TestBed.createComponent(Dorms);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
