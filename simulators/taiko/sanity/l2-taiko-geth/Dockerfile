FROM golang:latest as builder
# Build the simulator executable with local docker
ADD . /source/l2-taiko-geth
WORKDIR /source/l2-taiko-geth
RUN go build -v .

# or

# Build the simulator executable with Hive
#ADD ./sanity/l2-taiko-geth/hive.go.work /source/go.work
#ADD ./common /source/common
#ADD ./sanity/l2-taiko-geth /source/sanity/l2-taiko-geth
#WORKDIR /source/sanity/l2-taiko-geth
#RUN go build -v .

# Build the simulator executable with local Hive
#ADD ./ /source/
#WORKDIR /source/simulators/taiko/sanity/l2-taiko-geth
RUN #go build -v .


#Build the simulator run container.
FROM golang:latest
RUN apt update
RUN apt install -y iproute2 grep
WORKDIR /sim/
# local docker
COPY --from=builder /source/l2-taiko-geth/l2-taiko-geth .

# or

# Hive
#COPY --from=builder /source/sanity/l2-taiko-geth/l2-taiko-geth .

# or

# local Hive
#COPY --from=builder /source/simulators/taiko/sanity/l2-taiko-geth/l2-taiko-geth .

ENTRYPOINT ["/sim/l2-taiko-geth"]