FROM alpine:latest AS timezone_build
RUN apk --no-cache add tzdata ca-certificates  


FROM golang:1.20.6-alpine3.18 AS builder

RUN apk --no-cache add tzdata ca-certificates

ADD . /go/api

WORKDIR /go/api

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.7

RUN /go/bin/swag init -d adapter/cli --parseDependency --parseInternal --parseDepth 3 -o adapter/http/rest/docs

RUN mkdir deploy
RUN go clean --modcache
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o go_app adapter/cli/main.go 
RUN mv go_app ./deploy/go_app
RUN mv config.json ./deploy/config.json
RUN mv adapter/http/rest/docs/ ./deploy/docs
RUN mv database ./deploy/database

FROM scratch AS production

COPY --from=timezone_build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=timezone_build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/api/deploy /api/

WORKDIR /api

ENTRYPOINT  ["./go_app", "serve"]