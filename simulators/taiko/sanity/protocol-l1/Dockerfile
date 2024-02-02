FROM golang:latest as builder

# Build the simulator executable.
ADD . /source
WORKDIR /source
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
RUN apt update
RUN apt install -y iproute2 grep
ADD . /source
WORKDIR /source/
COPY --from=builder /source/protocol-l1 .
ENTRYPOINT ["./protocol-l1"]