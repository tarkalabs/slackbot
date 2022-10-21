FROM golang:1.19-buster

WORKDIR /go/src/slackbot
COPY go.* ./
RUN go mod download

RUN go install github.com/cosmtrek/air@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install golang.org/x/tools/gopls@latest \
  && go install gotest.tools/gotestsum@latest
