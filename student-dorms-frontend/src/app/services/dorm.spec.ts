import { TestBed } from '@angular/core/testing';

import { Dorm } from './dorm';

describe('Dorm', () => {
  let service: Dorm;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Dorm);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
