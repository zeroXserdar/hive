FROM gcr.io/evmchain/taiko-client:main

RUN apk add --update bash jq curl
RUN which bash
RUN echo $SHELL

COPY saved_env.txt /saved_env.txt

COPY enode.sh /hive-bin/enode.sh
RUN chmod +x /hive-bin/enode.sh
#ENV CHAIN_ID=31336
#ENV PERIOD=12
#ENV NODE_KEY="9c52ecfdf0397641379694b367a2bf4328580e4310753d4064ab0af8b132d565"
#
#RUN rm -rf /data/l1-node/geth/chaindata
#RUN rm -rf /data/l1-node/geth/lightchaindata
#RUN rm -rf /data/l1-node/nodes
#RUN rm -rf /data/l1-node/geth/triecache
#
#RUN /usr/local/bin/geth console --exec 'console.log(admin.nodeInfo.name)' --maxpeers=0 --nodiscover --dev 2>/dev/null | head -1 >/version.txt
RUN taiko-client -v >/version.txt

COPY ./start.sh /start.sh

RUN chmod u+x /start.sh

ENTRYPOINT ["/start.sh"]