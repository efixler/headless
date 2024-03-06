FROM golang:latest AS builder


RUN apt -y update && apt -y upgrade

WORKDIR /headless
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go install -v ./cmd/headless
RUN go install -v ./cmd/headless-proxy
WORKDIR /go/bin


FROM debian:12-slim

RUN apt -y update && apt -y upgrade
RUN apt-get -y install \
    ca-certificates \
    chromium \
    gnupg wget apt-transport-https

RUN mkdir -p /headless/bin
COPY --from=builder /go/bin/* /headless/bin/
EXPOSE 8080/tcp
ENV HEADLESS_PROXY_PORT="8080"
ENV HEADLESS_PROXY_DEFAULT_USER_AGENT=":firefox:"
CMD ["cd", "/"]
ENTRYPOINT ["/headless/bin/headless-proxy"]