package main

import (
	pb "akstechies/go-grpc/order"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderRequestData struct {
	OrderID  int    `json:"order_id"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

func handleOrderRequest(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data OrderRequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received data: %+v\n", data)
	fmt.Println("type of data", reflect.TypeOf(data))

	// ✅ Call the function with extracted data
	order_request(int32(data.OrderID), data.Item, int32(data.Quantity))

	// ✅ Respond with success message
	response := map[string]string{
		"status":  "success",
		"message": "Order received successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func order_request(order_id int32, item string, quantity int32) {
	// Connect to the gRPC server
	// if running without docker use `localhost:50051`
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	// Prepare the request
	req := &pb.OrderRequest{
		OrderId:  order_id,
		Item:     item,
		Quantity: quantity,
	}

	// Send the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateOrder(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	// Handle the response
	fmt.Printf("Order Status: %s\n", res.Status)
	fmt.Printf("Message: %s\n", res.Message)
}

func main() {
	http.HandleFunc("POST /submit", handleOrderRequest)

	order_request(1, "Lap", 4)

	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
