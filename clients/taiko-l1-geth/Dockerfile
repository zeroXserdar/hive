FROM gcr.io/evmchain/internal-devnet-l1-node:latest

RUN apk add --update bash jq curl

ENV CHAIN_ID=31336
ENV PERIOD=12
ENV NODE_KEY="9c52ecfdf0397641379694b367a2bf4328580e4310753d4064ab0af8b132d565"

#RUN rm -rf /data/l1-node/geth/chaindata
#RUN rm -rf /data/l1-node/geth/lightchaindata
#RUN rm -rf /data/l1-node/nodesenode.sh
#RUN rm -rf /data/l1-node/geth/triecache

COPY enode.sh /hive-bin/enode.sh
RUN chmod +x /hive-bin/enode.sh

RUN sed -i 's/\(--http\.api debug,eth,net,web3,txpool,miner\)/\1,admin/' /start.sh
RUN sed -i 's/\(--ws\.api debug,eth,net,web3,txpool,miner\)/\1,admin/' /start.sh

RUN /usr/local/bin/geth console --exec 'console.log(admin.nodeInfo.name)' --maxpeers=0 --nodiscover --dev 2>/dev/null | head -1 >/version.txt
