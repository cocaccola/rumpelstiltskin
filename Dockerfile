FROM golang:1.19 AS builder

WORKDIR /build
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN make linux


FROM alpine:3.17

WORKDIR /
RUN apk add --no-cache curl
COPY --from=builder /build/rumpelstiltskin ./
ENTRYPOINT [ "/rumpelstiltskin" ]
