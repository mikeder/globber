FROM --platform=linux/amd64 golang AS builder

RUN go version
RUN mkdir /src

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0

WORKDIR /src

COPY . .

RUN go build -o ./bin/admin ./cmd/admin
RUN go build -o ./bin/globber  ./cmd/globber

FROM alpine as globber
RUN apk --no-cache add ca-certificates curl

COPY --from=builder /src/bin/admin /bin/admin
COPY --from=builder /src/bin/globber /bin/globber
COPY --from=builder /src/static ./static/
COPY --from=builder /src/templates ./templates/

EXPOSE 3000
ENTRYPOINT /bin/globber
