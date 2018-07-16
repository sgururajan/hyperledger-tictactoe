import { AppConfigService } from "../core/config/app-config.service";


export function ApiFactory(serviceName: string) {
    return (configService: AppConfigService)=>{
        return configService.config.services[serviceName];
    }
}