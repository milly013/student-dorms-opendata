import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpendataPrices } from './opendata-prices';

describe('OpendataPrices', () => {
  let component: OpendataPrices;
  let fixture: ComponentFixture<OpendataPrices>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OpendataPrices]
    })
    .compileComponents();

    fixture = TestBed.createComponent(OpendataPrices);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
