import { ApiServiceModule } from './api-service.module';

describe('ApiServiceModule', () => {
  let apiServiceModule: ApiServiceModule;

  beforeEach(() => {
    apiServiceModule = new ApiServiceModule();
  });

  it('should create an instance', () => {
    expect(apiServiceModule).toBeTruthy();
  });
});
