FROM golang:1.13-alpine AS builder

RUN go version
RUN mkdir /src

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0

WORKDIR /src

COPY . .

RUN go build -o ./bin/globber  ./cmd/globber 

FROM alpine:3.9 
RUN apk --no-cache add ca-certificates

COPY --from=builder /src/bin/globber /bin/globber
COPY --from=builder /src/static/* /static/
COPY --from=builder /src/templates/* /templates/

EXPOSE 3000
 