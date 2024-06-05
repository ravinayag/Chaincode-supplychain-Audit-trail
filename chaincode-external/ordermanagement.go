package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var logger = log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

type serverConfig struct {
	CCID    string
	Address string
}

// SmartContract provides functions for managing supply chain, shipments, and payments
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

// TransactionType represents the type of transaction
type TransactionType string

const (
	ACHTransaction        TransactionType = "ACH"
	CreditCardTransaction TransactionType = "CreditCard"
	// Add more transaction types as needed
)

// TransactionData represents the structure of transactional data to be stored
type TransactionData struct {
	ID                 string          `json:"id"`
	Type               TransactionType `json:"type"`
	Amount             float64         `json:"amount"`
	Account            string          `json:"account"`
	TransactionDetails string          `json:"transactionDetails"`
	// Add more fields as needed
}

// ShipEngineData represents the structure of ShipEngine data to be stored
type ShipEngineData struct {
	ID          string `json:"id"`
	ShipmentID  string `json:"shipmentId"`
	TrackingURL string `json:"trackingUrl"`
	// Add more fields as needed
}

// Init initializes the chaincode
func (s *SmartContract) InitOrder(ctx contractapi.TransactionContextInterface) error {
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

// CreateOrder creates a new order in the supply chain
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

	// Check if order already exists
	exists, err := s.OrderExists(ctx, orderNo)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s : the order %s already exists", timestamp, orderNo)
	}

	// Create new order object
	order := Order{
		OrderNo:       orderNo,
		Date:          date,
		OrderDetail:   orderDetail,
		Invoice:       invoice,
		PackingStatus: packingStatus,
		PaymentMethod: paymentMethod,
		OrderTrack:    orderTrack,
	}

	// Marshal order object to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Save order to ledger
	err = ctx.GetStub().PutState(orderNo, orderJSON)
	if err != nil {
		return fmt.Errorf("failed to put order %s to world state: %v", orderNo, err)
	}

	// Log the success of the operation
	logger.Printf("%s : Order created successfully: %s", timestamp, orderNo)

	return nil
}

// ReadOrder retrieves an order from the ledger based on its order number
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

	// Retrieve order from ledger
	orderJSON, err := ctx.GetStub().GetState(orderNo)
	if err != nil {
		return nil, fmt.Errorf("%s : failed to read order %s from world state: %v", timestamp, orderNo, err)
	}
	if orderJSON == nil {
		return nil, fmt.Errorf("%s : the order %s does not exist", timestamp, orderNo)
	}

	// Unmarshal order JSON into Order struct
	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return nil, err
	}

	// Log the success of the operation
	logger.Printf("%s : Queried the order successfully: %s", timestamp, orderNo)

	return &order, nil
}

// UpdateOrder updates an existing order in the supply chain
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

	// Check if order exists
	exists, err := s.OrderExists(ctx, orderNo)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s : the order %s does not exist", timestamp, orderNo)
	}

	// Update existing order
	order := Order{
		OrderNo:       orderNo,
		Date:          date,
		OrderDetail:   orderDetail,
		Invoice:       invoice,
		PackingStatus: packingStatus,
		PaymentMethod: paymentMethod,
		OrderTrack:    orderTrack,
	}

	// Marshal updated order object to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Update order in ledger
	err = ctx.GetStub().PutState(orderNo, orderJSON)
	if err != nil {
		return fmt.Errorf("failed to update order %s in world state: %v", orderNo, err)
	}

	// Log the success of the operation
	logger.Printf("%s : Order updated successfully: %s", timestamp, orderNo)

	return nil
}

// DeleteOrder deletes an order from the supply chain
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

	// Check if order exists
	exists, err := s.OrderExists(ctx, orderNo)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s : the order %s does not exist", timestamp, orderNo)
	}

	// Delete order from ledger
	err = ctx.GetStub().DelState(orderNo)
	if err != nil {
		return fmt.Errorf("failed to delete order %s from world state: %v", orderNo, err)
	}

	// Log the success of the operation
	logger.Printf("%s : Order deleted successfully: %s", timestamp, orderNo)

	return nil
}

// OrderExists checks if an order exists in the supply chain
func (s *SmartContract) OrderExists(ctx contractapi.TransactionContextInterface, orderNo string) (bool, error) {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Retrieve order from ledger
	orderJSON, err := ctx.GetStub().GetState(orderNo)
	if err != nil {
		return false, fmt.Errorf("%s : failed to read order %s from world state: %v", timestamp, orderNo, err)
	}

	return orderJSON != nil, nil
}

// GetAllOrders returns all orders stored in the supply chain
func (s *SmartContract) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Retrieve all orders from ledger
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResult

	// Iterate over all orders
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Unmarshal order JSON into Order struct
		var order Order
		err = json.Unmarshal(queryResponse.Value, &order)
		if err != nil {
			return nil, err
		}

		// Create QueryResult object and append to results
		queryResult := QueryResult{Key: queryResponse.Key, Record: &order}
		results = append(results, queryResult)
	}

	// Log the success of the operation
	logger.Printf("%s : Queried all orders successfully", timestamp)

	return results, nil
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Order
}

// TransactionExists checks if a transaction with given ID exists in the ledger
func (s *SmartContract) TransactionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	dataBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return dataBytes != nil, nil
}

// Init ShiptEngine
func (s *SmartContract) InitShipEngine(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("ShipEngine chaincode initialized")
	return nil
}

// ShipEngineDataExists checks if ShipEngineData with given ID exists in the ledger
func (s *SmartContract) ShipEngineDataExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	dataBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return dataBytes != nil, nil
}

// CreateShipEngineData adds a new ShipEngineData to the ledger
func (s *SmartContract) CreateShipEngineData(ctx contractapi.TransactionContextInterface, id string, shipmentID string, trackingURL string) error {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Log the start of the function
	logger.Printf("%s : Creating ShipEngine data: %s", timestamp, id)

	// Check if ShipEngineData already exists
	exists, err := s.ShipEngineDataExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s : the ShipEngineData with ID %s already exists", timestamp, id)
	}

	// Create new ShipEngineData object
	data := ShipEngineData{
		ID:          id,
		ShipmentID:  shipmentID,
		TrackingURL: trackingURL,
	}

	// Marshal ShipEngineData object to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Save ShipEngineData to ledger
	err = ctx.GetStub().PutState(id, dataBytes)
	if err != nil {
		return err
	}

	// Log the success of the operation
	logger.Printf("%s : ShipEngine data created successfully: %s", timestamp, id)

	return nil
}

// CreateTransaction adds a new transaction to the ledger
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, id string, transactionTypeStr string, amount float64, account string, transactionDetails string) error {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Log the start of the function
	logger.Printf("%s : Creating transaction: %s", timestamp, id)

	transactionType := TransactionType(transactionTypeStr)

	// Check if transaction already exists
	exists, err := s.TransactionExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%s : the transaction with ID %s already exists", timestamp, id)
	}

	// Create new TransactionData object
	data := TransactionData{
		ID:                 id,
		Type:               transactionType,
		Amount:             amount,
		Account:            account,
		TransactionDetails: transactionDetails,
	}

	// Marshal TransactionData object to JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Save TransactionData to ledger
	err = ctx.GetStub().PutState(id, dataBytes)
	if err != nil {
		return err
	}

	// Log the success of the operation
	logger.Printf("%s : Transaction created successfully: %s", timestamp, id)

	return nil
}

// GetTransaction retrieves transaction from the ledger based on ID
func (s *SmartContract) GetTransaction(ctx contractapi.TransactionContextInterface, id string) (*TransactionData, error) {
	// Record the timestamp
	timestamp := time.Now().Format(time.RFC3339)
	logger.Printf("Timestamp: %s", timestamp)

	// Log the start of the function
	logger.Printf("%s : Retrieving transaction: %s", timestamp, id)

	// Retrieve transaction from ledger
	dataBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if dataBytes == nil {
		return nil, fmt.Errorf("the transaction with ID %s does not exist", id)
	}

	// Unmarshal TransactionData JSON into struct
	var data TransactionData
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}

	// Log the success of the operation
	logger.Printf("%s : Retrieved transaction successfully: %s", timestamp, id)

	return &data, nil
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
