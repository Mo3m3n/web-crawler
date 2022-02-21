FROM golang:1.17-alpine AS builder

COPY . /src/

RUN cd /src && go build -o webcrawler cmd/server/webcrawler.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /src/webcrawler .

ENTRYPOINT ["./webcrawler"]
CMD []
