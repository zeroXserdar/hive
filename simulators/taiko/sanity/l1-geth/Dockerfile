FROM golang:latest as builder

# Build the simulator executable.
ADD ./sanity/l1-geth/hive.go.work /source/go.work
ADD ./common /source/common
ADD ./sanity/l1-geth /source/sanity/l1-geth
WORKDIR /source/sanity/l1-geth
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
RUN apt update
RUN apt install -y iproute2 grep
WORKDIR /sim/
COPY --from=builder /source/sanity/l1-geth/l1-geth .
ENTRYPOINT ["/sim/l1-geth"]