import { TestBed, inject } from '@angular/core/testing';

import { ActionDispatchHelperService } from './action-dispatch-helper.service';

describe('ActionDispatchHelperService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ActionDispatchHelperService]
    });
  });

  it('should be created', inject([ActionDispatchHelperService], (service: ActionDispatchHelperService) => {
    expect(service).toBeTruthy();
  }));
});
