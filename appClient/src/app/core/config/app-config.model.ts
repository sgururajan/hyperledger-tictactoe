export interface AppConfig {
    productName:string
    services:{
        [apiKey:string]:string
    }
}