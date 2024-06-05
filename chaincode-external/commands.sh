# Go to Test Network Directory
cd worldtradex/wtx-supplychain-fabric-network

# Start the network & create channel
./network.sh up createChannel -s couchdb

# Set Chaincode PATH & Deploy the chaincode
./network.sh deployCC -ccn ordermanagement -ccp /Users/macbook/neel/Projects/worldtradex/wtx-supplychain-fabric-network/chaincode-external -ccl go

./network.sh deployCCAAS  -ccn ordermanagement -ccp /Users/macbook/neel/Projects/worldtradex/wtx-supplychain-fabric-network/chaincode-external

# ENV
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export PATH=../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config


docker build -f ./Dockerfile -t basicj_ccaas_image:latest --build-arg CC_SERVER_PORT=9999 ./chaincode-external


 docker run --rm -it -p 9229:9229 --name peer0org2_ordermanagement_ccaas --network fabric_test -e DEBUG=true -e CHAINCODE_SERVER_ADDRESS=0.0.0.0:9999 -e CHAINCODE_ID=ordermanagement_1.0:7c7dff5cdc43c77ccea028c422b3348c3c1fb5a26ace0077cf3cc627bd355ef0 -e CORE_CHAINCODE_ID_NAME=ordermanagement_1.0:7c7dff5cdc43c77ccea028c422b3348c3c1fb5a26ace0077cf3cc627bd355ef0 ordermanagement_ccaas_image:latest


 docker run -p 9999:9999 --name ordermanagement.org1.example.com --network fabric_test -e DEBUG=true -e CHAINCODE_SERVER_ADDRESS=0.0.0.0:9999 -e CHAINCODE_ID=ordermanagement_1.0.1:e907f2b717286686ffee3e89074b11dc0ed56ef06db67adefe643c2a9312d2aa -e CORE_CHAINCODE_ID_NAME=ordermanagement_1.0.1:e907f2b717286686ffee3e89074b11dc0ed56ef06db67adefe643c2a9312d2aa ordermanagement_ccaas_image:latest -d


# Invoke the chaincode 
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n ordermanagement --peerAddresses localhost:7051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitOrder","Args":[]}'


peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n ordermanagement --peerAddresses localhost:7051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitTransaction","Args":[]}'


peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n ordermanagement --peerAddresses localhost:7051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitShipEngine","Args":[]}'


# Query the chaincode
peer chaincode query -C mychannel -n ordermanagement -c '{"Args":["ReadOrder","logis_ordr_1"]}'


# Transaction Chaincode Functions
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n ordermanagement --peerAddresses localhost:7051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"Args":["CreateTransaction", "txn123", "ACH", "100.0", "1234567890", "Payment for services"]}'


# Query Transaction from Transaction ID
peer chaincode query -C mychannel -n ordermanagement -c '{"Args":["GetTransaction","txn123"]}'
