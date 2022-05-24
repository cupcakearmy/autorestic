FROM golang:1.18-alpine as builder

WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build

FROM alpine
RUN apk add --no-cache restic rclone bash openssh
COPY --from=builder /app/autorestic /usr/bin/autorestic
CMD [ "autorestic" ]
