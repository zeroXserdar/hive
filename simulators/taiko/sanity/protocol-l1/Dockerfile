FROM golang:latest as builder

# Build the simulator executable.
ADD . /source
WORKDIR /source/sanity/protocol-l1
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
RUN apt update
RUN apt install -y iproute2 grep
ADD . /source
WORKDIR /source/sanity/protocol-l1
COPY --from=builder /source/sanity/protocol-l1/protocol-l1 .
ENTRYPOINT ["./protocol-l1"]