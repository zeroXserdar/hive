FROM golang:latest as builder

# Build the simulator executable.
ADD . /source
WORKDIR /source
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
ADD . /source
WORKDIR /source
COPY --from=builder /source/sanity-l1-geth .
ENTRYPOINT ["./sanity-l1-geth"]