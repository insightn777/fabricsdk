##########################################
#  Dockerfile
#  Hyperledger/fabric
#  ieetu/such
##########################################
# 1. Base Image    // docker가 bionic(Ubuntu 18.04.1 LTS) 까지만 지원
FROM ubuntu:18.04 AS fabricbase
# dependencies install
RUN apt-get update && apt-get install -y apt-utils \
    && apt-get install -y curl python-pip git libltdl-dev tree openssh-server net-tools vim apt-transport-https software-properties-common
# go install & mkdir
RUN curl -sL https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz | tar -C /usr/local -zx && mkdir -p /root/ieetu/gopath /root/ieetu/myfabric
# go env    // 이상하게 PATH는 따로 ENV 안주면 적용안됨
ENV GOPATH=/root/ieetu/gopath \
    GOROOT=/usr/local/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin:$GOPATH/src/github.com/hyperledger/fabric-ca/bin:$GOPATH/src/github.com/hyperledger/fabric/.build/bin \
    CORE_LOGGING_GRPC=DEBUG \
    FABRIC_LOGGING_SPEC=debug \
    FABRIC_CFG_PATH=/root/ieetu/myfabric \
    FABRIC_CA_CLIENT_HOME=/root/ieetu/myfabric \
    FABRIC_HOME=$GOPATH/src/github.com/hyperledger/fabric
# docker install
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
    && apt-get update && apt-get install -y docker-ce && pip install docker-compose
# fabric core/ca git clone
RUN mkdir -p $GOPATH/src/github.com/hyperledger && cd $GOPATH/src/github.com/hyperledger \ 
    && git clone -b release-1.4 https://github.com/hyperledger/fabric \
    && git clone -b release-1.4 https://github.com/hyperledger/fabric-ca \
## fabric-ca-client install
    && cd $GOPATH/src/github.com/hyperledger/fabric-ca && make fabric-ca-client
WORKDIR /root/ieetu/myfabric
# ----------------------------------------
# 2. CA Image
FROM fabricdocker_base AS ca
# ca-client / server install
RUN mkdir /root/ieetu/myfabric/ca-server /root/ieetu/myfabric/ca-server-admin /root/ieetu/myfabric/ca-tls-admin /root/ieetu/myfabric/ca-tls \
    && cd $GOPATH/src/github.com/hyperledger/fabric-ca && make fabric-ca-server
ENV FABRIC_CA_SERVER_HOME=/root/ieetu/myfabric/ca-server
# ----------------------------------------
# 3. Orderer Image
FROM fabricdocker_base AS orderer
# orderer install     // FABRIC_CFG_PATH 에는 orderer.yaml 이 있어야됨
RUN cd $GOPATH/src/github.com/hyperledger/fabric && make orderer \
    && cd $FABRIC_HOME/sampleconfig && cp orderer.yaml $FABRIC_CFG_PATH
# ----------------------------------------
# 4. Peer Image
FROM fabricdocker_base AS peer
# peer install     // FABRIC_CFG_PATH 에는 core.yaml 이 있어야됨
RUN cd $GOPATH/src/github.com/hyperledger/fabric && make configtxgen && make peer \
    && cd $FABRIC_HOME/sampleconfig && cp core.yaml $FABRIC_CFG_PATH
# ----------------------------------------
# 5. Kafka Image
FROM hyperledger/fabric-kafka AS kafka
RUN echo "35.196.70.249   org0" >> /etc/hosts \
&& echo "35.243.169.182  org1" >> /etc/hosts \
&& echo "35.227.56.250   org2" >> /etc/hosts
# ----------------------------------------
# 6. Zookeeper Image
FROM hyperledger/fabric-zookeeper AS zookeeper
RUN echo "35.196.70.249   org0" >> /etc/hosts \
&& echo "35.243.169.182  org1" >> /etc/hosts \
&& echo "35.227.56.250   org2" >> /etc/hosts
# ----------------------------------------
# 7. Admin Image
# FROM fabricdocker_peer AS admin
# 
# RUN mkdir $GOPATH/bin \ 
#     && go get -u github.com/hyperledger/fabric-sdk-go 
# ----------------------------------------
###########################################################################################


