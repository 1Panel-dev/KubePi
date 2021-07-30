import {AfterViewInit, ChangeDetectorRef, Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {Terminal} from "xterm";
import {FitAddon} from "xterm-addon-fit";
import {debounce} from "lodash";
import {mergeScan} from "rxjs/operators";

declare let SockJS: any;


@Component({
  selector: 'app-logging',
  templateUrl: './logging.component.html',
  styleUrls: ['./logging.component.css']
})
export class LoggingComponent implements AfterViewInit {

  constructor(private readonly cdr_: ChangeDetectorRef) {
  }

  @ViewChild('anchor', {static: true}) anchorRef: ElementRef;
  term: Terminal;
  namespace: string;
  podName: string;
  clusterName: string;
  container: string;
  tailLines: number;
  follow: boolean;

  private debouncedFit_: Function


  ngAfterViewInit(): void {
    this.initTerm()
    this.setupConnection()
    // if (this.namespace && this.podName && this.container) {
    //   // this.setupConnection()
    //   // this.initTerm()
    // } else {
    //   alert("please set param: namespace,pod  and container name ")
    // }
  }

  initTerm() {
    this.term = new Terminal({
      fontSize: 12,
      fontFamily: 'Consolas, "Courier New", monospace',
      bellStyle: 'sound',
      cursorBlink: true,
    })
    const fitAddon = new FitAddon();
    this.term.loadAddon(fitAddon);
    this.term.open(this.anchorRef.nativeElement)
    this.debouncedFit_ = debounce(() => {
      fitAddon.fit()
      this.cdr_.markForCheck();
    }, 100)
    this.debouncedFit_();
  }

  setupConnection() {
    let options = {
      debug: true,
      devel: true,
      protocols_whitelist: ['websocket', 'xdr-streaming', 'xhr-streaming', 'iframe-eventsource', 'iframe-htmlfile', 'xdr-polling', 'xhr-polling', 'iframe-xhr-polling', 'jsonp-polling']
    };
    let sock = new SockJS('/echo', undefined, options);
    sock.onopen = () => {
      console.log("open")
      sock.send("aaaaaa")
    }
    sock.onmessage = (msg: any) => {
      console.log(msg)
      this.term.write(msg.data)
    }
    sock.onclose = () => {
      console.log("connection close")
    }
  }
}
