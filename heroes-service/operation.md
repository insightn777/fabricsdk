##########################################################
#
#
#    구 동 순 서
#
#
##########################################################

# docker image 다 지우기
docker image rm $(docker images --format "{{.ID}}")
docker stop $(docker ps --format "{{.ID}}")

# docker file 다 옮기기

#docker build -t fabricbase -f Dockerfile.base .

rm -r /root/ieetu/myfabric/*

export HOSTIP=35.196.70.249
export ORG0=35.196.70.249
export ORG1=35.243.169.182
export ORG2=35.227.56.250
export TLS=true

docker-compose up -d

cd ..
vim ca-tls/fabric-ca-tls-config.yaml
docker cp fabric-ca-tls-config.yaml ca:/root/ieetu/myfabric/ca-tls

###### CA #####

# CA 서버 가동	  // FABRIC_CA_SERVER_HOME에 MSP 들어감

docker exec -it ca /bin/bash

fabric-ca-server start -b admin:adminpw --cfg.affiliations.allowremove --cfg.identities.allowremove --ca.name ca0
fabric-ca-server start -b admin:adminpw --cfg.affiliations.allowremove --cfg.identities.allowremove --ca.name ca0 --cafiles ../ca-tls/fabric-ca-tls-config.yaml &

# CA 관리자 등록     //FABRIC_CA_CLIENT_HOME에 MSP 들어감
fabric-ca-client enroll --caname ca0 -H /root/ieetu/myfabric/ca-server-admin -u http://admin:adminpw@$HOSTIP:7054
fabric-ca-client enroll --caname ca0-tls -H /root/ieetu/myfabric/ca-tls-admin -u http://admin:adminpw@$HOSTIP:7054

# 조직 정리
fabric-ca-client affiliation remove --force org1 -H /root/ieetu/myfabric/ca-tls-admin --caname ca0-tls
fabric-ca-client affiliation remove --force org2 -H /root/ieetu/myfabric/ca-tls-admin --caname ca0-tls
fabric-ca-client affiliation add Org0 -H /root/ieetu/myfabric/ca-tls-admin --caname ca0-tls

fabric-ca-client affiliation remove --force org1 -H /root/ieetu/myfabric/ca-server-admin --caname ca0
fabric-ca-client affiliation remove --force org2 -H /root/ieetu/myfabric/ca-server-admin --caname ca0
fabric-ca-client affiliation add Org0 -H /root/ieetu/myfabric/ca-server-admin --caname ca0

# tls 등록
fabric-ca-client register --id.name ca0 --id.secret ca0pw --id.affiliation Org0 --id.type user -H /root/ieetu/myfabric/ca-tls-admin --caname ca0-tls -u http://$HOSTIP:7054
fabric-ca-client register --id.name admin0 --id.secret admin0pw --id.affiliation Org0 --id.type user -H /root/ieetu/myfabric/ca-tls-admin --caname ca0-tls -u http://$HOSTIP:7054
fabric-ca-client register --id.name peer0 --id.secret peer0pw --id.affiliation Org0 --id.type user -H /root/ieetu/myfabric/ca-tls-admin --caname ca0-tls -u http://$HOSTIP:7054
fabric-ca-client register --id.name orderer0 --id.secret orderer0pw --id.affiliation Org0 --id.type user -H /root/ieetu/myfabric/ca-tls-admin --caname ca0-tls -u http://$HOSTIP:7054

### tls enabled
fabric-ca-client enroll -u http://ca0:ca0pw@$HOSTIP:7054 -H tls --caname ca0-tls -m ca0

rm tls/f*
mv tls/msp/cacerts/*.pem tls/tlsca.crt
mv tls/msp/signcerts/*.pem tls/server.crt
mv tls/msp/keystore/*_sk tls/server.key
rm -r tls/msp

# 조직 관리자 등록    // FABRIC_CA_CLIENT_HOME에 MSP 들어감 admin=true:ecert || role=admin:ecert
fabric-ca-client register -H /root/ieetu/myfabric/ca-server-admin --caname ca0 --id.affiliation Org0 --id.name admin0 --id.secret admin0pw --id.type client --id.maxenrollments 0 --id.attrs '"hf.Registrar.Roles=client,orderer,peer,user","hf.Registrar.DelegateRoles=client,orderer,peer,user",hf.Registrar.Attributes=*,hf.GenCRL=true,hf.Revoker=true,hf.AffiliationMgr=true,hf.IntermediateCA=true,admin=true:ecert'
exit

##### ADMIN #####
docker exec -it admin /bin/bash

# 조직 관리자 MSP 생성
fabric-ca-client enroll -u http://admin0:admin0pw@$HOSTIP:7054 --caname ca0
mv msp/keystore/*_sk msp/keystore/serverkey.key

# admincerts 폴더 생성
cp -r /root/ieetu/myfabric/msp/signcerts /root/ieetu/myfabric/msp/admincerts


# 구성원 MSP 생성
fabric-ca-client register --id.name peer0 --id.secret peer0pw --id.affiliation Org0 --id.type peer --id.maxenrollments 0 --id.attrs 'role=peer:ecert' -H /root/ieetu/myfabric -u http://$HOSTIP:7054 --caname ca0
fabric-ca-client register --id.name orderer0 --id.secret orderer0pw --id.affiliation Org0 --id.type orderer --id.maxenrollments 0 --id.attrs 'role=orderer:ecert' -H /root/ieetu/myfabric -u http://$HOSTIP:7054 --caname ca0

# admin TLS enabled
mkdir tls

fabric-ca-client enroll -u http://admin0:admin0pw@$HOSTIP:7054 -H tls --caname ca0-tls -m admin0

rm tls/f*
mv tls/msp/cacerts/*.pem tls/tlsca.crt
mv tls/msp/signcerts/*.pem tls/server.crt
mv tls/msp/keystore/*_sk tls/server.key
rm -r tls/msp

# org msp 폴더 생성
mkdir -p Org0/tlscacerts
cp -r msp/admincerts Org0/
cp -r msp/cacerts Org0/
cp tls/tlsca.crt Org0/tlscacerts

# admin, peer, orderer 각각에 DNS 등
echo "35.196.70.249   orderer0 org0 peer0 admin0" >> /etc/hosts
echo "35.243.169.182  orderer1 org1 peer1 admin1" >> /etc/hosts
echo "35.227.56.250   orderer2 org2 peer2 admin2" >> /etc/hosts
exit

####### ORDERER #######
docker exec -it orderer /bin/bash

# ORDERER MSP 생성
fabric-ca-client enroll -u http://orderer0:orderer0pw@$HOSTIP:7054 --caname ca0

# TLS enabled
mkdir tls
fabric-ca-client enroll -u http://orderer0:orderer0pw@$HOSTIP:7054 -H tls --caname ca0-tls -m orderer0

rm tls/f*
mv tls/msp/cacerts/*.pem tls/tlsca.crt
mv tls/msp/signcerts/*.pem tls/server.crt
mv tls/msp/keystore/*_sk tls/server.key
rm -r tls/msp

# admin, peer, orderer 각각에 DNS 등
echo "35.196.70.249   orderer0 org0 peer0 admin0" >> /etc/hosts
echo "35.243.169.182  orderer1 org1 peer1 admin1" >> /etc/hosts
echo "35.227.56.250   orderer2 org2 peer2 admin2" >> /etc/hosts
exit

###### PEER ######
docker exec -it peer /bin/bash

# PEER MSP 생성
fabric-ca-client enroll -u http://peer0:peer0pw@$HOSTIP:7054 --caname ca0

# TLS enabled
mkdir tls
fabric-ca-client enroll -u http://peer0:peer0pw@$HOSTIP:7054 -H tls --caname ca0-tls -m peer0

rm tls/f*
mv tls/msp/cacerts/*.pem tls/tlsca.crt
mv tls/msp/signcerts/*.pem tls/server.crt
mv tls/msp/keystore/*_sk tls/server.key
rm -r tls/msp

# admin, peer, orderer 각각에 DNS 등
echo "35.196.70.249   orderer0 org0 peer0 admin0" >> /etc/hosts
echo "35.243.169.182  orderer1 org1 peer1 admin1" >> /etc/hosts
echo "35.227.56.250   orderer2 org2 peer2 admin2" >> /etc/hosts
exit

# admincerts 폴더 생성
cd /root/ieetu/myfabric
docker cp admin:/root/ieetu/myfabric/Org0 .

docker cp /root/ieetu/myfabric/Org0/admincerts orderer:/root/ieetu/myfabric/msp/
docker cp /root/ieetu/myfabric/Org0/admincerts peer:/root/ieetu/myfabric/msp/

gsutil cp -r Org0 gs://fabricsuch/

# start peer
docker exec -it peer bash

##### ORDERER 구동 #####

###### 다른 org cert 가져옴  ($FABRIC_CA_CLIENT_HOME의 각각 admin msp 안에 admincert && cacert 필요)
#?# fabric-ca-client getcacerts   // 다른 ca 서버의 cacerts 가져옴

cd /root/ieetu/myfabric/
gsutil cp -r gs://fabricsuch/Org1 .
gsutil cp -r gs://fabricsuch/Org2 .

docker cp Org1 admin:/root/ieetu/myfabric/ && docker cp Org2 admin:/root/ieetu/myfabric/
docker cp Org1/tlscacerts/tlsca.crt orderer:/root/ieetu/myfabric/tls/tlsca1.crt
docker cp Org2/tlscacerts/tlsca.crt orderer:/root/ieetu/myfabric/tls/tlsca2.crt

# configtx.yaml 작성
docker exec -it admin bash
vim configtx.yaml

# genesisblock 생성
configtxgen -profile TaskGenesis -outputBlock genesis.block
exit

docker cp admin:/root/ieetu/myfabric/genesis.block .
docker cp genesis.block orderer:/root/ieetu/myfabric
gsutil cp genesis.block gs://fabricsuch/

# orderer 실행  // 폴더에 genesis.block orderer.yaml 파일 다 있어야됨
docker exec -it orderer bash
orderer 

# 다른 orderer 다 시작시킴

##### CHANNER 생성 #####

docker exec -it admin bash

# create Channel ( admin )
configtxgen -profile TaskChannel -outputCreateChannelTx ch1.tx -channelID ch1

peer channel create -o orderer0:7050 -c ch1 -f ch1.tx --tls --cafile /root/ieetu/myfabric/tls/tlsca.crt --certfile /root/ieetu/myfabric/tls/server.crt --keyfile /root/ieetu/myfabric/tls/server.key

# peerJoin ( admin )
peer channel join -b ch1.block

# Anchorpeer Update (admin0 )
configtxgen -profile TaskChannel -outputAnchorPeersUpdate Org0Anchor.tx -channelID ch1 -asOrg Org0

peer channel create -o orderer0:7050 -c ch1 -f Org0Anchor.tx --tls --cafile /root/ieetu/myfabric/tls/tlsca.crt --certfile /root/ieetu/myfabric/tls/server.crt --keyfile /root/ieetu/myfabric/tls/server.key

# 다른 peer join & Anchorpeer update

docker cp admin:/root/ieetu/myfabric/ch1.block .
gsutil cp ch1.block gs://fabricsuch/


##### CHAINCODE 생성 #####

# chainhero 옮김
docker cp chainHero admin:/root/ieetu/gopath/src/github.com/

curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# install chaincode ( admin ) 
peer chaincode install -n heroes-service -v 1.0 -p github.com/chainHero/heroes-service/chaincode --tls --cafile /root/ieetu/myfabric/tls/tlsca.crt --certfile /root/ieetu/myfabric/tls/server.crt --keyfile /root/ieetu/myfabric/tls/server.key

# instantiate chaincode
peer chaincode instantiate -o orderer0:7050 -C ch1 -n heroes-service -v 1.0 -c '{"Args":["init",""]}' --tls --cafile /root/ieetu/myfabric/tls/tlsca.crt --certfile /root/ieetu/myfabric/tls/server.crt --keyfile /root/ieetu/myfabric/tls/server.key

# query
peer chaincode query -C ch1 -n heroes-service -c '{"Args":["invoke","query","hello"]}' --tls --cafile /root/ieetu/myfabric/tls/tlsca.crt --certfile /root/ieetu/myfabric/tls/server.crt --keyfile /root/ieetu/myfabric/tls/server.key

# invoke
peer chaincode invoke -o orderer0:7050 -C ch1 -n heroes-service -c '{"Args":["invoke","invoke","hello","asdf"]}' --tls --cafile /root/ieetu/myfabric/tls/tlsca.crt --certfile /root/ieetu/myfabric/tls/server.crt --keyfile /root/ieetu/myfabric/tls/server.key

peer logging setlogspec gossip=warning:msp=warning:info

dep ensure

go build



###### Trouble shooting ######
##############################

# channel create error
# because it doesn't contain any IP SANs 

vim /etc/ssl/openssl.cnf

# edit

[ v3_ca ]
subjectAltName = IP : 10.142.0.13

#==================================

# channel create error
# first record does not look like a TLS handshake

일단 TLS 꺼버림

#==================================

# file 수정 시
cd /root/ieetu/gopath/src/github.com/hyperledger/fabric && make peer
cd /root/ieetu/myfabric && ./createChannel.sh
vim /root/ieetu/gopath/src/github.com/hyperledger/fabric/peer/channel/create.go

# file 여러 개인 경우

mkdir $FABRIC_CA_CLIENT_HOME/tls
cp $(find $FABRIC_CA_CLIENT_HOME/msp/signcerts -name '*.pem') $FABRIC_CA_CLIENT_HOME/tls/server.crt
cp $(find $FABRIC_CA_CLIENT_HOME/msp/cacerts -name '*.pem') $FABRIC_CA_CLIENT_HOME/tls/ca.crt
cp $(find $FABRIC_CA_CLIENT_HOME/msp/keystore -name '*_sk') $FABRIC_CA_CLIENT_HOME/tls/server.key
exit

###################################


