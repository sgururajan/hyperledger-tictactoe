import { Observable } from "rxjs";
import {throwError} from "rxjs/internal/observable/throwError";

export abstract class BaseApiService {
    protected extract(res:Response) {
        const body = res.json()||{};
        if (res.status<200 || res.status>300) {
            const genericBody =<any>res.json();
            const genericMessage="server error: " + res.status;
            const errMsg = genericBody.message || genericMessage;
            throw new Error(errMsg);
        }
        return body;
    }

    protected handleError(error:any) {
        const errMsg=error.message || "server error";
        return throwError(errMsg);
        //return Observable.throw(errMsg);
    }
}
