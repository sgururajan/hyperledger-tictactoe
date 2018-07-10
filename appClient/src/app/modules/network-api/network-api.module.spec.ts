import { NetworkApiModule } from './network-api.module';

describe('NetworkApiModule', () => {
  let networkApiModule: NetworkApiModule;

  beforeEach(() => {
    networkApiModule = new NetworkApiModule();
  });

  it('should create an instance', () => {
    expect(networkApiModule).toBeTruthy();
  });
});
