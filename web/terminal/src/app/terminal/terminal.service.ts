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

  createTerminalSession(clusterName: string, namespace: string, podName: string, containerName: string, shell: string): Observable<any> {
    const url = function () {
      let baseUrl = `/kubepi/api/v1/clusters/${clusterName}/terminal/session?podName=${podName}&&containerName=${containerName}&&shell=${shell}`
      if (namespace) {
        baseUrl = `${baseUrl}&&namespace=${namespace}`
      }
      return baseUrl
    }()
    return this.http.get<any>(url)
  }
  createNodeShellTerminalSession(clusterName: string, nodeName: string): Observable<any> {
    const url = function () {
      let baseUrl = `/kubepi/api/v1/clusters/${clusterName}/node_terminal/session?nodeName=${nodeName}`
      return baseUrl
    }()
    return this.http.get<any>(url)
  }
}
