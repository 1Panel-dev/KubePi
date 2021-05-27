export interface PodContainerList {
  containers: string[];
}

export interface ShellFrame {
  Op: string;
  Data?: string;
  SessionID?: string;
  Rows?: number;
  Cols?: number;
}

export interface TerminalResponse {
  id: string
}

export interface Pod {
  kind: string
  apiVersion: string
  metadata: PodMetadata
  containers: Container[]
}

export interface PodMetadata {
  name: string;
}

export interface Container {
  name: string;
  image: string;
}

export interface SockJSSimpleEvent {
  type: string;

  toString(): string;
}

export interface SJSCloseEvent extends SockJSSimpleEvent {
  code: number;
  reason: string;
  wasClean: boolean;
}

export interface SJSMessageEvent extends SockJSSimpleEvent {
  data: string;
}

