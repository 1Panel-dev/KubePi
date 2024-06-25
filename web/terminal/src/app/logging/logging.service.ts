import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class LoggingService {

  constructor(private http: HttpClient) {
  }

  createLoggingSession(clusterName: string,
                       namespace: string,
                       podName: string,
                       containerName: string,
                       tailLines: number,
                       follow: boolean,
                       previous: boolean ): Observable<any> {
    const url = function () {
      let baseUrl = `/kubepi/api/v1/clusters/${clusterName}/logging/session?podName=${podName}`
      if (namespace) {
        baseUrl = `${baseUrl}&&namespace=${namespace}`
      }
      if (containerName) {
        baseUrl += `&&containerName=${containerName}`
      }
      if (tailLines) {
        baseUrl += `&&tailLines=${tailLines}`
      }
      if (follow) {
        baseUrl += `&&follow=${follow}`
      }
      if (previous) {
        baseUrl += `&&previous=${previous}`
      }
      return baseUrl
    }()
    return this.http.get<any>(url)
  }
}
