import grpc
import order_pb2
import order_pb2_grpc

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = order_pb2_grpc.OrderServiceStub(channel)
    
    request = order_pb2.OrderRequest(order_id=1, item="Laptop", quantity=2)
    response = stub.CreateOrder(request)
    
    print(f"Order Status: {response.status}")
    print(f"Message: {response.message}")

if __name__ == "__main__":
    run()
