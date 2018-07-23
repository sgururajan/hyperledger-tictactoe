import { TestBed, inject } from '@angular/core/testing';

import { NetworkAuthInterceptorService } from './network-auth-interceptor.service';

describe('NetworkAuthInterceptorService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NetworkAuthInterceptorService]
    });
  });

  it('should be created', inject([NetworkAuthInterceptorService], (service: NetworkAuthInterceptorService) => {
    expect(service).toBeTruthy();
  }));
});
