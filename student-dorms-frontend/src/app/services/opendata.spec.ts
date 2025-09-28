import { TestBed } from '@angular/core/testing';

import { Opendata } from './opendata';

describe('Opendata', () => {
  let service: Opendata;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Opendata);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
