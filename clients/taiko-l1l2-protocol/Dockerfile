FROM gcr.io/evmchain/internal-devnet-protocol:latest

ENV PRIVATE_KEY='0xf214f2b2cd398c806f84e317254e0f0b801d0643303237d97a22a48e01628897'
#ENV MAINNET_URL='http://172.17.0.2:8545'
#ENV PROTOCOL_LAYER=L1/L2
ENV PROTOCOL_LAYER=L1
ENV LOG_LEVEL=debug
#ENV MAINNET_URL=programmatically generated
#ENV NETWORK_HOST=programmatically generated
#ENV PRIVATE_KEY=(L2 deployer private key)(secret)
#ENV L2_GENESIS_HASH=(L2 genesis hash)
#ENV TAIKO_L2_ADDRESS(L2 taiko address)
#ENV L2_SIGNAL_SERVICE=(L2 signal service l2 address)
#ENV GUARDIAN_PROVERS=(L2 guardian prover addresses)(join ",")
#ENV CONTRACT_OWNERS=(L2 contract owner)
#ENV TAIKO_TOKEN_PREMINT_RECIPIENT=(L2 token recipient address)
#ENV SHARED_ADDRESS_MANAGER=(L2 value: "0x0000000000000000000000000000000000000000")
#ENV TIMELOCK_CONTROLLER="0x0000000000000000000000000000000000000000"

ENV PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
ENV PROPOSER=0x0000000000000000000000000000000000000000
ENV PROPOSER_ONE=0x0000000000000000000000000000000000000000
ENV GUARDIAN_PROVERS="0x1000777700000000000000000000000000000001,0x1000777700000000000000000000000000000002,0x1000777700000000000000000000000000000003,0x1000777700000000000000000000000000000004,0x1000777700000000000000000000000000000005"
ENV MIN_GUARDIANS=3
ENV TAIKO_L2_ADDRESS=0x1000777700000000000000000000000000000001
ENV L2_SIGNAL_SERVICE=0x1000777700000000000000000000000000000007
ENV SECURITY_COUNCIL=0x60997970C51812dc3A010C7d01b50e0d17dc79C8
ENV TAIKO_TOKEN_PREMINT_RECIPIENT=0xa0Ee7A142d267C1f36714E4a8F75612F20a79720
ENV TAIKO_TOKEN_NAME="Taiko Token Katla"
ENV TAIKO_TOKEN_SYMBOL=TTKOk
ENV SHARED_ADDRESS_MANAGER=0x0000000000000000000000000000000000000000
ENV L2_GENESIS_HASH=0xee1950562d42f0da28bd4550d88886bc90894c77c9c9eaefef775d4c8223f259

#RUN git clone --depth=1 git@github.com:taikoxyz/taiko-mono.git
#WORKDIR /taiko-mono/packages/protocol
#RUN pnpm compile && pnpm generate:genesis ../../../k8s-configs/internal-devnet/config/genesis/internal_l2a.js > /config.output

COPY ./internal_l2a.js /internal_l2a.js
WORKDIR /taiko-mono/packages/protocol
ENV PATH="/root/.foundry/bin:${PATH}"

#RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | /bin/bash
#RUN /bin/bash -c "source ~/.nvm/nvm.sh"
#RUN /bin/bash -c "source ~/.nvm/nvm.sh && nvm --version && nvm install stable && nvm use stable && node --version"

#RUN /bin/bash -c "source ~/.nvm/nvm.sh && pnpm install"
#RUN /bin/bash -c "source ~/.nvm/nvm.sh && pnpm install husky tshy"
#RUN /bin/bash -c 'source ~/.nvm/nvm.sh && pnpm install "@nomicfoundation/hardhat-verify"'
#RUN /bin/bash -c 'source ~/.nvm/nvm.sh && pnpm install "@nomicfoundation/hardhat-ethers"'
#RUN /bin/bash -c "source ~/.nvm/nvm.sh && pnpm compile && pnpm generate:genesis /internal_l2a.js > /config.output"

RUN pnpm install
RUN pnpm install husky tshy
RUN pnpm install "@nomicfoundation/hardhat-verify"
RUN pnpm install "@nomicfoundation/hardhat-ethers"
RUN pnpm compile && pnpm generate:genesis /internal_l2a.js > /config.output
RUN #pnpm compile && pnpm generate:genesis ../../../k8s-configs/internal-devnet/config/genesis/internal_l2a.js > /config.output

RUN apt-get update
RUN apt-get install -y netcat-traditional
COPY ./start.sh /start.sh

RUN chmod u+x /start.sh

ENTRYPOINT ["/start.sh"]