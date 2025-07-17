FROM golang:1.24.5-bookworm

RUN go install github.com/cespare/reflex@latest

COPY article-feed-worker-reflex.conf /usr/local/etc/reflex.conf

WORKDIR /app

VOLUME /go

CMD ["reflex", "-d", "none", "-c", "/usr/local/etc/reflex.conf"]
