import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IssueRequests } from './issue-requests';

describe('IssueRequests', () => {
  let component: IssueRequests;
  let fixture: ComponentFixture<IssueRequests>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [IssueRequests]
    })
    .compileComponents();

    fixture = TestBed.createComponent(IssueRequests);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
