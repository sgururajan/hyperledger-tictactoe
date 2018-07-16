import { TestBed, inject } from '@angular/core/testing';

import { NetworkAuthGuardService } from './network-auth-guard.service';

describe('NetworkAuthGuardService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NetworkAuthGuardService]
    });
  });

  it('should be created', inject([NetworkAuthGuardService], (service: NetworkAuthGuardService) => {
    expect(service).toBeTruthy();
  }));
});
