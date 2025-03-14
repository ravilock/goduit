FROM golang:1.24.1-bookworm

WORKDIR /app

RUN go install github.com/codegangsta/gin@latest
RUN git config --system --add safe.directory '*'

CMD ["gin", "-b", "goduit", "run"]
