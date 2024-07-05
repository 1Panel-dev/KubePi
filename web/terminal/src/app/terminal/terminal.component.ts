import {AfterViewInit, ChangeDetectorRef, Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {Terminal} from 'xterm';
import {FitAddon} from 'xterm-addon-fit';
import {ReplaySubject, Subject} from "rxjs";
import {ShellFrame, SJSCloseEvent, SJSMessageEvent} from "./terminal";
import {ActivatedRoute, Router} from "@angular/router";
import {debounce} from 'lodash';
import {takeUntil} from "rxjs/operators";
import {TerminalService} from "./terminal.service";

declare let SockJS: any;


@Component({
  selector: 'app-terminal',
  templateUrl: './terminal.component.html',
  styleUrls: ['./terminal.component.css']
})

export class TerminalComponent implements AfterViewInit {
  @ViewChild('anchor', {static: true}) anchorRef: ElementRef;
  term: Terminal;
  podName: string;
  container: string;
  shell: string;

  private readonly namespace_: string
  clusterName: string;
  private connecting_: boolean
  private connectionClosed_: boolean
  private conn_: WebSocket
  private connected_ = false;
  private debouncedFit_: Function
  private connSubject_ = new ReplaySubject<ShellFrame>(100)
  private incommingMessage$_ = new Subject<ShellFrame>()
  // private readonly endpoint_=
  private readonly unsubscribe_ = new Subject<void>()
  private readonly keyEvent$_ = new ReplaySubject<KeyboardEvent>(2)


  constructor(private activatedRoute_: ActivatedRoute,
              private readonly cdr_: ChangeDetectorRef,
              private _router: Router,
              private terminalService: TerminalService
  ) {
    this.clusterName = this.activatedRoute_.snapshot.queryParams["cluster"]
    this.namespace_ = this.activatedRoute_.snapshot.queryParams["namespace"]
    this.podName = this.activatedRoute_.snapshot.queryParams["pod"]
    this.container = this.activatedRoute_.snapshot.queryParams["container"]
    this.shell = this.activatedRoute_.snapshot.queryParams["shell"]
  }


  ngAfterViewInit(): void {
    if (this.namespace_ && this.podName && this.container) {
      this.setupConnection()
    } else {
      alert("please set param: namespace,pod  and container name ")
    }
  }

  ngOnDestroy(): void {
    this.unsubscribe_.next();
    this.unsubscribe_.complete();

    if (this.conn_) {
      this.conn_.close();
    }

    if (this.connSubject_) {
      this.connSubject_.complete();
    }

    if (this.term) {
      this.term.dispose();
    }
    this.incommingMessage$_.complete();
  }

  onPodContainerChange(podContainer: string): void {
    this._router.navigate([`/terminal/${this.namespace_}/${this.podName}/${podContainer}`], {
      queryParamsHandling: 'preserve',
    });
  }

  disconnect(): void {
    if (this.conn_) {
      this.conn_.close();
    }

    if (this.connSubject_) {
      this.connSubject_.complete();
      this.connSubject_ = new ReplaySubject<ShellFrame>(100);
    }

    if (this.term) {
      this.term.dispose();
    }

    this.incommingMessage$_.complete();
    this.incommingMessage$_ = new Subject<ShellFrame>();
  }


  initTerm(): void {
    if (this.connSubject_) {
      this.connSubject_.complete()
      this.connSubject_ = new ReplaySubject<ShellFrame>(100);
    }

    if (this.term) {
      this.term.dispose()
    }

    this.term = new Terminal({
      fontSize: 14,
      fontFamily: 'Consolas, "Courier New", monospace',
      bellStyle: 'sound',
      cursorBlink: true,
    });
    const fitAddon = new FitAddon();
    this.term.loadAddon(fitAddon);
    this.term.open(this.anchorRef.nativeElement);
    this.debouncedFit_ = debounce(() => {
      fitAddon.fit()
      this.cdr_.markForCheck();
    }, 100)
    this.debouncedFit_();
    window.addEventListener('resize', () => this.debouncedFit_())
    this.connSubject_.pipe(takeUntil(this.unsubscribe_)).subscribe(frame => {
      this.handleConnectionMessage(frame);
    })
    this.term.onData(this.onTerminalSendingString.bind(this))
    this.term.onResize(this.onTerminalResize.bind(this))
    this.term.onKey(event => {
      this.keyEvent$_.next(event.domEvent)
    })
    this.cdr_.markForCheck()
  }

  private onConnectionOpen(sessionId: string): void {
    const startData = {Op: 'bind', SessionID: sessionId};

    this.conn_.send(JSON.stringify(startData));
    this.connSubject_.next(startData);
    this.connected_ = true;
    this.connecting_ = false;
    this.connectionClosed_ = false;

    // Make sure the terminal is with correct display size.
    this.onTerminalResize();

    // Focus on connection
    this.term.focus();
    this.cdr_.markForCheck();
  }


  private async setupConnection(): Promise<void> {
    if (!(this.container && this.podName && this.namespace_ && !this.connecting_)) {
      return;
    }

    this.connecting_ = true;
    this.connectionClosed_ = false;

    try {
      const {data} = await this.terminalService.createTerminalSession(this.clusterName, this.namespace_, this.podName, this.container, this.shell).toPromise()
      const id = data.id
      this.conn_ = new SockJS(`/kubepi/api/v1/ws/terminal/sockjs?${id}`);
      this.conn_.onopen = this.onConnectionOpen.bind(this, id);
      this.conn_.onmessage = this.onConnectionMessage.bind(this);
      this.conn_.onclose = this.onConnectionClose.bind(this);
      this.initTerm()
      this.cdr_.markForCheck();
    } catch (e:any) {
      this.term.write(e.error.message)
    }
  }

  private handleConnectionMessage(frame: ShellFrame): void {
    if (frame.Op === 'stdout') {
      if (frame.Data)
        this.term.write(frame.Data);
    }

    if (frame.Op === 'toast') {
      alert(frame.Data)
    }

    this.incommingMessage$_.next(frame);
    this.cdr_.markForCheck();
  }

  private onConnectionMessage(evt: SJSMessageEvent): void {
    const msg = JSON.parse(evt.data);
    this.connSubject_.next(msg);
  }

  private onConnectionClose(_evt?: SJSCloseEvent): void {
    this.term.write("\n\n******** The connection failed. Unsupported, interrupted or timed out. ********")
    if (!this.connected_) {
      return;
    }
    this.conn_.close();
    this.connected_ = false;
    this.connecting_ = false;
    this.connectionClosed_ = true;
    // alert(_evt?.reason)
    this.cdr_.markForCheck();
  }

  private onTerminalSendingString(str: string): void {
    if (this.connected_) {
      this.conn_.send(JSON.stringify({
        Op: 'stdin',
        Data: str,
        Cols: this.term.cols,
        Rows: this.term.rows,
      }))
    }
  }

  private onTerminalResize(): void {
    if (this.connected_) {
      this.conn_.send(
        JSON.stringify({
          Op: 'resize',
          Cols: this.term.cols,
          Rows: this.term.rows,
        })
      );
    }
  }
}
