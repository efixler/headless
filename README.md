# headless

## Execution plan

### Go is the controller

Go will launch a local headless chrome, since remote debugging seems janky. Start this as a script and then port it to a server

In either case the Go process launches and then launches the browser

## Chrome Headless Launcher

-- Set a smallish viewport, user agent and such


### yukinying

https://github.com/yukinying/chrome-headless-browser-docker?tab=readme-ov-file

(Dockerfile: https://github.com/yukinying/chrome-headless-browser-docker/blob/master/chrome-stable/Dockerfile) 

debian-stable:slim

some weird apt-gets and command line options 

### justinribeiro

This one is deboan buster (10, old, 2022)

https://hub.docker.com/r/justinribeiro/chrome-headless/

(Dockerfile: https://github.com/justinribeiro/dockerfiles/blob/master/chrome-headless/Dockerfile)

## DevTools Adaptor Libraries

- https://github.com/go-rod/rod
- https://github.com/chromedp/chromedp
