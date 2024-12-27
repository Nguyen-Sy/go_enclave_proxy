import socket

# Define the target enclave's vSock CID and port
VSOCK_CID = 17  # This is typically the CID of the Nitro Enclave
VSOCK_PORT = 5006

def connect_to_vsock_server():
    # Create a vsock socket
    client_socket = socket.socket(socket.AF_VSOCK, socket.SOCK_STREAM)

    # Connect to the Nitro Enclave server (using the vSock CID and port)
    client_socket.connect((VSOCK_CID, VSOCK_PORT))
    print(f"Connected to vSock Ping-Pong server on CID {VSOCK_CID} and port {VSOCK_PORT}")

    try:
        # Send a "ping" message to the server
        message = "ping"
        client_socket.sendall(message.encode())
        print(f"Sent: {message}")

        # Receive the response ("pong") from the server
        response = client_socket.recv(1024)
        print(f"Received: {response.decode()}")
    finally:
        client_socket.close()
        print("Connection closed.")

if __name__ == '__main__':
    connect_to_vsock_server()