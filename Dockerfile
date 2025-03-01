# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER=1.17
ARG ALPINE_VER=3.16

FROM golang:${GO_VER}-alpine${ALPINE_VER}

WORKDIR /Users/macbook/neel/Projects/worldtradex/wtx-supplychain-fabric-network/chaincode-external
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 9999
CMD ["chaincode-external"]
