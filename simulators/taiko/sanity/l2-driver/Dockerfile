FROM golang:latest as builder

# Build the simulator executable.
ADD ./sanity/l2-driver/hive.go.work /source/go.work
ADD ./common /source/common
ADD ./sanity/l2-driver /source/sanity/l2-driver
WORKDIR /source/sanity/l2-driver
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
RUN apt update
RUN apt install -y iproute2 grep
WORKDIR /sim/
COPY --from=builder /source/sanity/l2-driver/l2-driver .
ENTRYPOINT ["/sim/l2-driver"]