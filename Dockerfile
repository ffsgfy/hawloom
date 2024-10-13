FROM golang:1.23-alpine AS build
WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd/
COPY internal ./internal/
RUN go build ./cmd/hawloom

FROM alpine:3.20
WORKDIR /hawloom

USER root
RUN apk add curl

USER guest
COPY --from=build /build/hawloom ./

CMD ["./hawloom"]
