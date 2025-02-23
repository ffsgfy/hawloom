FROM golang:1.23-alpine AS build
WORKDIR /build

RUN apk add git
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833

COPY cmd ./cmd
COPY db ./db
COPY internal ./internal
COPY go.mod go.sum sqlc.yaml ./

RUN sqlc generate
RUN templ generate -path ./internal/ui
RUN go build ./cmd/hawloom

FROM alpine:3.20
WORKDIR /hawloom

USER root
RUN apk add curl

USER guest
COPY --from=build /build/hawloom /usr/bin/

ENTRYPOINT ["/usr/bin/hawloom"]
