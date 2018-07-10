import { AppConfigService } from "./app-config.service";

export function AppConfigFactory(configService: AppConfigService) {
    return ()=> configService.load();
}