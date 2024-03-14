# headless

[![Go Reference](https://pkg.go.dev/badge/github.com/efixler/headless.svg)](https://pkg.go.dev/github.com/efixler/headless)
[![Go Report Card](https://goreportcard.com/badge/github.com/efixler/headless)](https://goreportcard.com/report/github.com/efixler/headless)
[![License MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://github.com/efixler/headless?tab=MPL-2.0-1-ov-file)

`headless` scrapes html for a target url using a headless Chrome browser. The included `headless` and `headless-proxy` apps provide headless scraping functionality from the shell and as an HTTP proxy server.

## Table of Contents

- [Usage as a CLI Application](#usage-as-a-cli-application)
- [Usage as a Server](#usage-as-a-proxy-server)
- [Roadmap](#roadmap)
- [Acknowledgements](#acknowledgements)

## Usage As a CLI Application

### Installation

```
go install github.com/efixler/headless
```

### Usage

```
headless % ./build/headless -h
Usage: 
        headless [flags] :url
 
  -h
        Show this help message
  -H    Show browser window (don't run in headless mode)
        Environment: HEADLESS_NO_HEADLESS
  -log-level value
        Log level
        Environment: HEADLESS_LOG_LEVEL
  -user-agent value
        User agent to use (omit for browser default)
        Environment: HEADLESS_USER_AGENT
```

## Usage As a Proxy Server 

`headless-proxy` is currently experimental. It's functional as a proof-of-concept but not ready for usage
in production environments.

### Installation

```
go install github.com/efixler/headless-proxy
```

### Usage

```
headless % ./build/headless-proxy -h
Usage: 
        headless-proxy [flags] :url
 
  -h
        Show this help message
  -default-user-agent value
        Default user agent string (empty for browser default)
        Environment: HEADLESS_PROXY_DEFAULT_USER_AGENT
  -inbound-idle-timeout value
        Inbound connection keepalive idle timeout
        Environment: HEADLESS_PROXY_IDLE_TIMEOUT (default 2m0s)
  -inbound-read-timeout value
        Inbound connection read timeout
        Environment: HEADLESS_PROXY_READ_TIMEOUT (default 5s)
  -inbound-write-timeout value
        Inbound connection write timeout
        Environment: HEADLESS_PROXY_WRITE_TIMEOUT (default 30s)
  -log-level value
        Set the log level [debug|error|info|warn]
        Environment: HEADLESS_PROXY_LOG_LEVEL
  -max-concurrent value
        Maximum concurrent connections
        Environment: HEADLESS_PROXY_MAX_CONCURRENT (default 6)
  -port value
        Port to listen on
        Environment: HEADLESS_PROXY_PORT (default 8008)
```

## Roadmap

- Implemenent Proxy Authorization
- Build docker container
- Add https support
- Document https usage in perimeter-http environments
- Implement better header checking, url verification, etc.
- Proxy inbound user agent to outbound connection