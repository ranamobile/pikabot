FROM golang:1.13.2-stretch

RUN mkdir -p /app
WORKDIR /app
ENV GOPATH /app

COPY src /app/src
COPY entry /app/entry

RUN go get -d ./...
# RUN go get golang.org/x/net/context
# RUN go get golang.org/x/oauth2
# RUN go get golang.org/x/oauth2/google
# RUN go get google.golang.org/api/drive/v3
# RUN go get github.com/nlopes/slack

VOLUME [ "/app" ]
ENTRYPOINT [ "go", "run", "/app/src/pikabot.go" ]
EXPOSE 8080
