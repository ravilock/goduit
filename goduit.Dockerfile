FROM golang:1.25.1-bookworm

WORKDIR /app

RUN go install github.com/codegangsta/gin@latest
RUN git config --system --add safe.directory '*'

CMD ["gin", "-b", "goduit", "run"]
