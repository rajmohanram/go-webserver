FROM alpine:latest

WORKDIR /app

COPY bin/server-linux-arm64 /app/server
COPY web/ /app/web/

EXPOSE 8080

CMD ["./server"]
