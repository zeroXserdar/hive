FROM golang:latest as builder

# Build the simulator executable.
ADD ./sanity/protocol-l2/hive.go.work /source/go.work
ADD ./common /source/common
ADD ./sanity/protocol-l2 /source/sanity/protocol-l2
WORKDIR /source/sanity/protocol-l2
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
RUN apt update
RUN apt install -y iproute2 grep
WORKDIR /sim/
COPY --from=builder /source/sanity/protocol-l2/protocol-l2 .
ENTRYPOINT ["/sim/protocol-l2"]