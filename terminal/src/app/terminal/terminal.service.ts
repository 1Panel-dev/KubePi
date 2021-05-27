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

  createTerminalSession(namespace: string, podName: string, containerName: string): Observable<any> {
    const url = function () {
      let baseUrl = `/api/terminal/session?podName=${podName}&&containerName=${containerName}`
      if (namespace) {
        baseUrl = `${baseUrl}&&namespace=${namespace}`
      }
      return baseUrl
    }()
    return this.http.get<any>(url)
  }

  readPod(podName: string, namespace?: string): Observable<Pod> {
    const url = function () {
      let baseUrl = '/api/proxy/api/v1'
      if (namespace) {
        baseUrl = `${baseUrl}/namespace/${namespace}`
      }
      baseUrl = `${baseUrl}/pods/${podName}`
      return baseUrl
    }()
    return this.http.get<Pod>(url)
  }
}
