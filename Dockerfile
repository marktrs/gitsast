FROM golang:1.20 as builder

WORKDIR /api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/server

FROM alpine:3.17.1

WORKDIR /app

COPY --from=builder /api/Makefile  /app/Makefile
COPY --from=builder /api/bin/server /app/bin/server
COPY --from=builder /api/config/*.yaml /app/config/
COPY --from=builder /api/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/bin/server
RUN chmod +x /app/entrypoint.sh

CMD ["./entrypoint.sh"]
