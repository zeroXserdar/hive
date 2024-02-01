FROM golang:latest as builder

# Build the simulator executable.
ADD . /source
#WORKDIR /source/sanity/l1-geth/
WORKDIR /source
RUN go build -v .

# Build the simulator run container.
FROM golang:latest
ADD . /source
WORKDIR /source/sanity/l1-geth
COPY --from=builder /source/sanity/l1-geth/l1-geth .
ENTRYPOINT ["./l1-geth"]