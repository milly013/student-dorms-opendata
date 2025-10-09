import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateDorm } from './create-dorm';

describe('CreateDorm', () => {
  let component: CreateDorm;
  let fixture: ComponentFixture<CreateDorm>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateDorm]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateDorm);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
