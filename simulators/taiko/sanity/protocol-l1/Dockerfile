FROM golang:latest as builder

# Build the simulator executable.
ADD ./sanity/protocol-l1/hive.go.work /source/go.work
ADD ./common /source/common
ADD ./sanity/protocol-l1 /source/sanity/protocol-l1
WORKDIR /source/sanity/protocol-l1
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
RUN apt update
RUN apt install -y iproute2 grep
WORKDIR /sim/
COPY --from=builder /source/sanity/protocol-l1/protocol-l1 .
ENTRYPOINT ["/sim/protocol-l1"]