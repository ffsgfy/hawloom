FROM golang:1.23-alpine AS build
WORKDIR /build

RUN apk add git curl libstdc++
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833
RUN curl -L https://github.com/tailwindlabs/tailwindcss/releases/download/v4.0.8/tailwindcss-linux-x64-musl -o ./tailwindcss && chmod a+x ./tailwindcss

COPY cmd ./cmd
COPY db ./db
COPY internal ./internal
COPY go.mod go.sum sqlc.yaml ./

RUN sqlc generate
RUN templ generate -path ./internal/ui
RUN ./tailwindcss -i ./internal/ui/css/style.css -o ./static/style.css --minify
RUN go build ./cmd/hawloom

FROM alpine:3.20
WORKDIR /hawloom

USER root
RUN apk add curl

USER guest
COPY --from=build /build/hawloom /usr/bin/
COPY --from=build /build/static ./static

ENTRYPOINT ["/usr/bin/hawloom"]
