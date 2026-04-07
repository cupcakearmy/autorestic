FROM golang:1.25-alpine as builder

WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build

FROM restic/restic:0.17.3
RUN apk add --no-cache rclone bash curl docker-cli tini
COPY --from=builder /app/autorestic /usr/bin/autorestic
ENTRYPOINT [ "/sbin/tini", "--" ]
CMD [ "autorestic" ]
