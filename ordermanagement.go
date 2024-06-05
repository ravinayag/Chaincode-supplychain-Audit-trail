package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"os"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var logger = log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

type serverConfig struct {
	CCID    string
	Address string
}

// SmartContract provides functions for managing logistics
type SmartContract struct {
	contractapi.Contract
}

// Order represents the order information
type Order struct {
	OrderNo       string `json:"orderNo"`
	Date          string `json:"date"`
	OrderDetail   string `json:"orderDetail"`
	Invoice       string `json:"invoice"`
	PackingStatus string `json:"packingStatus"`
	PaymentMethod string `json:"paymentMethod"`
	OrderTrack    string `json:"orderTrack"`
}

// InitLedger initializes the ledger with some sample data
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	orders := []Order{
		{OrderNo: "logis_ordr_1", Date: "2024-03-01", OrderDetail: "Sample order details 1", Invoice: "INV-001", PackingStatus: "Packing", PaymentMethod: "Credit Card", OrderTrack: "In Progress"},
		{OrderNo: "logis_ordr_2", Date: "2024-03-02", OrderDetail: "Sample order details 2", Invoice: "INV-002", PackingStatus: "Packing", PaymentMethod: "Cash", OrderTrack: "Shipped"},
	}

	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(order.OrderNo, orderJSON)
		if err != nil {
			return fmt.Errorf("failed to put order %s to world state: %v", order.OrderNo, err)
		}
	}

	return nil
}

// CreateOrder creates a new order in the world state with the given details
func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface, orderNo, date, orderDetail, invoice, packingStatus, paymentMethod, orderTrack string) error {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Log the start of the function
	logger.Printf("%s : Creating order: %s", timestamp, orderNo)

	// Retrieve transaction ID and caller ID
	txID := ctx.GetStub().GetTxID()
	callerID, _ := ctx.GetClientIdentity().GetID()

	// Log transaction details
	logger.Printf("%s : Transaction ID: %s, Caller ID: %s", timestamp, txID, callerID)

	// Log parameter details
	logger.Printf("%s : Parameters - OrderNo: %s, Date: %s, OrderDetail: %s, Invoice: %s, PackingStatus: %s, PaymentMethod: %s, OrderTrack: %s",
		timestamp, orderNo, date, orderDetail, invoice, packingStatus, paymentMethod, orderTrack)

	exists, err := s.OrderExists(ctx, orderNo)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s : the order %s already exists", timestamp, orderNo)
	}

	order := Order{
		OrderNo:       orderNo,
		Date:          date,
		OrderDetail:   orderDetail,
		Invoice:       invoice,
		PackingStatus: packingStatus,
		PaymentMethod: paymentMethod,
		OrderTrack:    orderTrack,
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Log the success of the operation
	logger.Printf("%s : Order created successfully: %s", timestamp, orderNo)

	return ctx.GetStub().PutState(orderNo, orderJSON)
}

// ReadOrder retrieves an order from the world state based on its order number
func (s *SmartContract) ReadOrder(ctx contractapi.TransactionContextInterface, orderNo string) (*Order, error) {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Log the start of the function
	logger.Printf("%s : Reading order: %s", timestamp, orderNo)

	// Retrieve transaction ID and caller ID
	txID := ctx.GetStub().GetTxID()
	callerID, _ := ctx.GetClientIdentity().GetID()

	// Log transaction details
	logger.Printf("%s : Transaction ID: %s, Caller ID: %s", timestamp, txID, callerID)

	// Log parameter details
	logger.Printf("%s : Parameters - OrderNo: %s", timestamp, orderNo)

	orderJSON, err := ctx.GetStub().GetState(orderNo)
	if err != nil {
		return nil, fmt.Errorf("%s : failed to read order %s from world state: %v", timestamp, orderNo, err)
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s : the order %s does not exist", timestamp, orderNo)
	}

	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return nil, err
	}

	// Log the success of the operation
	logger.Printf("%s : Queried the order successfully: %s", timestamp, orderNo)

	return &order, nil
}

// UpdateOrder updates an existing order in the world state with the provided details
func (s *SmartContract) UpdateOrder(ctx contractapi.TransactionContextInterface, orderNo, date, orderDetail, invoice, packingStatus, paymentMethod, orderTrack string) error {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Log the start of the function
	logger.Printf("%s : Updating order: %s", timestamp, orderNo)

	// Retrieve transaction ID and caller ID
	txID := ctx.GetStub().GetTxID()
	callerID, _ := ctx.GetClientIdentity().GetID()

	// Log transaction details
	logger.Printf("%s : Transaction ID: %s, Caller ID: %s", timestamp, txID, callerID)

	// Log parameter details
	logger.Printf("%s : Parameters - OrderNo: %s", timestamp, orderNo)

	exists, err := s.OrderExists(ctx, orderNo)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s : the order %s does not exist", timestamp, orderNo)
	}

	order := Order{
		OrderNo:       orderNo,
		Date:          date,
		OrderDetail:   orderDetail,
		Invoice:       invoice,
		PackingStatus: packingStatus,
		PaymentMethod: paymentMethod,
		OrderTrack:    orderTrack,
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Log the success of the operation
	logger.Printf("%s : Order updated successfully: %s", timestamp, orderNo)

	return ctx.GetStub().PutState(orderNo, orderJSON)
}

// DeleteOrder deletes an order from the world state based on its order number
func (s *SmartContract) DeleteOrder(ctx contractapi.TransactionContextInterface, orderNo string) error {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Log the start of the function
	logger.Printf("%s : Deleting order: %s", timestamp, orderNo)

	// Retrieve transaction ID and caller ID
	txID := ctx.GetStub().GetTxID()
	callerID, _ := ctx.GetClientIdentity().GetID()

	// Log transaction details
	logger.Printf("%s : Transaction ID: %s, Caller ID: %s", timestamp, txID, callerID)

	// Log parameter details
	logger.Printf("%s : Parameters - OrderNo: %s", timestamp, orderNo)

	exists, err := s.OrderExists(ctx, orderNo)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s : the order %s does not exist", timestamp, orderNo)
	}

	// Log the success of the operation
	logger.Printf("%s : Order deleted successfully: %s", timestamp, orderNo)

	return ctx.GetStub().DelState(orderNo)
}

// OrderExists checks if an order exists in the world state based on its order number
func (s *SmartContract) OrderExists(ctx contractapi.TransactionContextInterface, orderNo string) (bool, error) {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	orderJSON, err := ctx.GetStub().GetState(orderNo)
	if err != nil {
		return false, fmt.Errorf("%s : failed to read order %s from world state: %v", timestamp, orderNo, err)
	}

	return orderJSON != nil, nil
}

// GetAllOrders returns all orders stored in the world state
func (s *SmartContract) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var order Order
		err = json.Unmarshal(queryResponse.Value, &order)
		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: &order}
		results = append(results, queryResult)
	}

	// Log the success of the operation
	logger.Printf("%s : All Assets Queried successfully", timestamp)
	return results, nil
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Order
}

// GetHistoryForKey returns the history of changes for a given order number
func (s *SmartContract) GetHistoryForKey(ctx contractapi.TransactionContextInterface, orderNo string) ([]HistoryQueryResult, error) {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Log the start of the function
	logger.Printf("%s : Retrieving history for order: %s", timestamp, orderNo)

	// Retrieve transaction ID and caller ID
	txID := ctx.GetStub().GetTxID()
	callerID, _ := ctx.GetClientIdentity().GetID()

	// Log transaction details
	logger.Printf("%s : Transaction ID: %s, Caller ID: %s", timestamp, txID, callerID)

	// Log parameter details
	logger.Printf("%s : Parameters - OrderNo: %s", timestamp, orderNo)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(orderNo)
	if err != nil {
		return nil, fmt.Errorf("%s : failed to get history for order %s: %v", timestamp, orderNo, err)
	}
	defer resultsIterator.Close()

	var history []HistoryQueryResult
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var order Order
		if queryResponse.IsDelete {
			order = Order{}
		} else {
			err = json.Unmarshal(queryResponse.Value, &order)
			if err != nil {
				return nil, err
			}
		}

		historyQueryResult := HistoryQueryResult{
			TxId:      queryResponse.TxId,
			Timestamp: time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String(),
			IsDelete:  queryResponse.IsDelete,
			Order:     order,
		}
		history = append(history, historyQueryResult)
	}

	// Log the success of the operation
	logger.Printf("%s : History retrieved successfully for order: %s", timestamp, orderNo)

	return history, nil
}

// HistoryQueryResult structure used for handling result of history query
type HistoryQueryResult struct {
	TxId      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"isDelete"`
	Order     Order  `json:"order"`
}

func main() {
	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("error creating logistics chaincode: %s", err)
	}

	server := &shim.ChaincodeServer{
		CCID:     config.CCID,
		Address:  config.Address,
		CC:       chaincode,
		TLSProps: getTLSProperties(),
	}

	if err := server.Start(); err != nil {
		log.Panicf("error starting logistics chaincode: %s", err)
	}
}

func getTLSProperties() shim.TLSProperties {
	tlsDisabledStr := getEnvOrDefault("CHAINCODE_TLS_DISABLED", "true")
	key := getEnvOrDefault("CHAINCODE_TLS_KEY", "")
	cert := getEnvOrDefault("CHAINCODE_TLS_CERT", "")
	clientCACert := getEnvOrDefault("CHAINCODE_CLIENT_CA_CERT", "")

	tlsDisabled := getBoolOrDefault(tlsDisabledStr, false)
	var keyBytes, certBytes, clientCACertBytes []byte
	var err error

	if !tlsDisabled {
		keyBytes, err = os.ReadFile(key)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
		certBytes, err = os.ReadFile(cert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}
	if clientCACert != "" {
		clientCACertBytes, err = os.ReadFile(clientCACert)
		if err != nil {
			log.Panicf("error while reading the crypto file: %s", err)
		}
	}

	return shim.TLSProperties{
		Disabled:      tlsDisabled,
		Key:           keyBytes,
		Cert:          certBytes,
		ClientCACerts: clientCACertBytes,
	}
}

func getEnvOrDefault(env, defaultVal string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		value = defaultVal
	}
	return value
}

func getBoolOrDefault(value string, defaultVal bool) bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultVal
	}
	return parsed
}

