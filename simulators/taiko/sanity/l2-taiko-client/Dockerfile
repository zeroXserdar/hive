FROM golang:latest as builder

# Build the simulator executable.
ADD ./sanity/l2-taiko-client/hive.go.work /source/go.work
ADD ./common /source/common
ADD ./sanity/l2-taiko-client /source/sanity/l2-taiko-client
WORKDIR /source/sanity/l2-taiko-client
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
RUN apt update
RUN apt install -y iproute2 grep
WORKDIR /sim/
COPY --from=builder /source/sanity/l2-taiko-client/l2-taiko-client .
ENTRYPOINT ["/sim/l2-taiko-client"]