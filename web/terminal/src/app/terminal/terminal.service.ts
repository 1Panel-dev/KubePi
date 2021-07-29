import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Pod} from "./terminal";

@Injectable({
  providedIn: 'root'
})
export class TerminalService {

  constructor(private http: HttpClient) {
  }

  createTerminalSession(clusterName: string, namespace: string, podName: string, containerName: string): Observable<any> {
    const url = function () {
      let baseUrl = `/api/v1/clusters/${clusterName}/terminal/session?podName=${podName}&&containerName=${containerName}`
      if (namespace) {
        baseUrl = `${baseUrl}&&namespace=${namespace}`
      }
      return baseUrl
    }()
    return this.http.get<any>(url)
  }
}
