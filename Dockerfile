FROM golang:1.13.2-stretch

RUN mkdir -p /app/src/github.com/ranamobile
WORKDIR /app
ENV GOPATH /app

COPY pikabot /app/src/github.com/ranamobile/pikabot
COPY entry /app/entry

# RUN go get -d ./...
RUN go get golang.org/x/net/context
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/google
RUN go get google.golang.org/api/drive/v3
RUN go get github.com/nlopes/slack

VOLUME [ "/app/etc" ]
ENTRYPOINT [ "go", "run", "/app/entry/pikabot.go" ]
EXPOSE 8080
