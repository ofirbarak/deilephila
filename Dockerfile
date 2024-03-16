FROM golang:1.22

WORKDIR /usr/src/deilephila

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build

CMD make build-plugin && ./bin/deilephila ./configTemplate.yml ./bin/plugin.so