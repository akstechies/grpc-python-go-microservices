import grpc
from concurrent import futures
import time

import order_pb2
import order_pb2_grpc

# Implement the gRPC service
class OrderService(order_pb2_grpc.OrderServiceServicer):
    def CreateOrder(self, request, context):
        print(f"Received order: {request.order_id}, Item: {request.item}, Quantity: {request.quantity}")
        return order_pb2.OrderResponse(
            status="SUCCESS",
            message=f"Order {request.order_id} for {request.item} created successfully!"
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    order_pb2_grpc.add_OrderServiceServicer_to_server(OrderService(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Order Service running on port 50051...")
    try:
        while True:
            time.sleep(86400)  # Keep running
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()
