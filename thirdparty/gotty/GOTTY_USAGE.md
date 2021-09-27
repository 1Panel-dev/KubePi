# Usage

```
Usage: gotty [options] <command> [<arguments...>]
```

Run `gotty` with your preferred command as its arguments (e.g. `gotty top`).

By default, GoTTY starts a web server at port 8080. Open the URL on your web browser and you can see the running command as if it were running on your terminal.

## Options

```
--address value, -a value     IP address to listen (default: "0.0.0.0") [$GOTTY_ADDRESS]
--port value, -p value        Port number to liten (default: "8080") [$GOTTY_PORT]
--permit-write, -w            Permit clients to write to the TTY (BE CAREFUL) [$GOTTY_PERMIT_WRITE]
--credential value, -c value  Credential for Basic Authentication (ex: user:pass, default disabled) [$GOTTY_CREDENTIAL]
--random-url, -r              Add a random string to the URL [$GOTTY_RANDOM_URL]
--random-url-length value     Random URL length (default: 8) [$GOTTY_RANDOM_URL_LENGTH]
--tls, -t                     Enable TLS/SSL [$GOTTY_TLS]
--tls-crt value               TLS/SSL certificate file path (default: "~/.gotty.crt") [$GOTTY_TLS_CRT]
--tls-key value               TLS/SSL key file path (default: "~/.gotty.key") [$GOTTY_TLS_KEY]
--tls-ca-crt value            TLS/SSL CA certificate file for client certifications (default: "~/.gotty.ca.crt") [$GOTTY_TLS_CA_CRT]
--index value                 Custom index.html file [$GOTTY_INDEX]
--title-format value          Title format of browser window (default: "{{ .command }}@{{ .hostname }}") [$GOTTY_TITLE_FORMAT]
--reconnect                   Enable reconnection [$GOTTY_RECONNECT]
--reconnect-time value        Time to reconnect (default: 10) [$GOTTY_RECONNECT_TIME]
--max-connection value        Maximum connection to gotty (default: 0) [$GOTTY_MAX_CONNECTION]
--once                        Accept only one client and exit on disconnection [$GOTTY_ONCE]
--timeout value               Timeout seconds for waiting a client(0 to disable) (default: 0) [$GOTTY_TIMEOUT]
--permit-arguments            Permit clients to send command line arguments in URL (e.g. http://example.com:8080/?arg=AAA&arg=BBB) [$GOTTY_PERMIT_ARGUMENTS]
--width value                 Static width of the screen, 0(default) means dynamically resize (default: 0) [$GOTTY_WIDTH]
--height value                Static height of the screen, 0(default) means dynamically resize (default: 0) [$GOTTY_HEIGHT]
--ws-origin value             A regular expression that matches origin URLs to be accepted by WebSocket. No cross origin requests are acceptable by default [$GOTTY_WS_ORIGIN]
--term value                  Terminal name to use on the browser, one of xterm or hterm. (default: "xterm") [$GOTTY_TERM]
--close-signal value          Signal sent to the command process when gotty close it (default: SIGHUP) (default: 1) [$GOTTY_CLOSE_SIGNAL]
--close-timeout value         Time in seconds to force kill process after client is disconnected (default: -1) (default: -1) [$GOTTY_CLOSE_TIMEOUT]
--config value                Config file path (default: "~/.gotty") [$GOTTY_CONFIG]
--version, -v                 print the version
```

### Config File

You can customize default options and your terminal (hterm) by providing a config file to the `gotty` command. GoTTY loads a profile file at `~/.gotty` by default when it exists.

```
// Listen at port 9000 by default
port = "9000"

// Enable TSL/SSL by default
enable_tls = true

// hterm preferences
// Smaller font and a little bit bluer background color
preferences {
    font_size = 5
    background_color = "rgb(16, 16, 32)"
}
```

See the [`.gotty`](https://github.com/yudai/gotty/blob/master/.gotty) file in this repository for the list of configuration options.
