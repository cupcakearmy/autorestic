FROM golang:1.18-alpine as builder

WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build

FROM alpine
RUN apk add --no-cache restic rclone bash openssh
COPY --from=builder /app/autorestic /usr/bin/autorestic
COPY entrypoint.sh /entrypoint.sh
COPY crond.sh /crond.sh
RUN chmod +x /entrypoint.sh /crond.sh
# show autorestic cron logs in docker
RUN ln -sf /proc/1/fd/1 /var/log/autorestic-cron.log
# run autorestic-cron every minute
RUN echo -e "*/1 * * * * bash /crond.sh" >> /etc/crontabs/root

CMD [ "/entrypoint.sh" ]
