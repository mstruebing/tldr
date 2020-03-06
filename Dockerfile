FROM golang:1.12.4-alpine as build

RUN apk add --no-cache git
WORKDIR /tldr
COPY . /tldr
RUN GO111MODULE=on CGO_ENABLED=0 go build -o bin/tldr cmd/tldr/main.go

#

FROM alpine:latest

RUN apk add --no-cache git
WORKDIR /tldr/
COPY --from=build /tldr/bin/tldr /bin/

CMD ["tldr"]
