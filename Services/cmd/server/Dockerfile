# syntax=docker/dockerfile:1
FROM golang AS builder
ARG timezone=Asia/Bangkok
ENV TZ $timezone
WORKDIR /app

RUN apt-get update && \
    apt-get -y install tini && \
    apt-get -y clean

COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN go test -v ./... && \
    CGO_ENABLED=0 go build -v \
    -installsuffix 'static' \
    -ldflags="-X 'main.version=$(git rev-parse --short HEAD)' -X 'main.build=$(date --iso-8601=seconds)'" \
    -o dist/server ./cmd/server

FROM gcr.io/distroless/base:nonroot
LABEL maintainer="Chatchanan Panyaprasirtkit <changnoi2547@gmail.com>"

ARG timezone=Asia/Bangkok

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV TZ $timezone

WORKDIR /app
COPY --from=builder /usr/bin/tini-static /usr/bin/tini
COPY --from=builder /app/dist/server /app/server
COPY ./configs /app/configs
COPY ./internal/assets /app/internal/assets

EXPOSE 8080 8443

ENTRYPOINT ["tini", "--", "/app/server"]
CMD ["--env=prd"]
