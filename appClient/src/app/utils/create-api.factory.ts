import { AppConfigService } from "../core/config/app-config.service";

export function CreateApiFactory(serviceName: string) {
    console.log("createapifactory");
    return (configService: AppConfigService)=> {
        console.log(`app config service: ${configService}`);
        console.log(configService.config.services[serviceName]);
        return configService.config.services[serviceName];
    };
}