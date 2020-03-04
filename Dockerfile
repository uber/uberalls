FROM golang:1.14.0-alpine3.11 as go

WORKDIR /go/src/uberalls

COPY . .

RUN apk add --no-cache git gcc g++ &&\
  go get github.com/Masterminds/glide &&\
  glide install &&\
  go get -v &&\
  go build

###########################################
FROM alpine:3.11

ARG UID

ENV UBERALLS_CONFIG=/home/docker.json

WORKDIR /home

RUN adduser -D -u $UID -h /home user &&\
  mkdir volumes &&\
  chown user:user volumes

COPY --chown=user config/docker.json $UBERALLS_CONFIG

COPY --from=go --chown=user  /go/src/uberalls/uberalls uberalls

USER user 

EXPOSE 3000

VOLUME [ "/home/volumes" ]

CMD ["/home/uberalls"]
