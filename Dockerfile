FROM golang:1.19-alpine as builder

WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build

FROM restic/restic:0.15.0
RUN apk add --no-cache rclone bash
COPY --from=builder /app/autorestic /usr/bin/autorestic
ENTRYPOINT []
CMD [ "autorestic" ]
