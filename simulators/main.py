import argparse
import hashlib
import json
import uuid
import requests
import time
import threading
import math
from datetime import datetime, timezone
from parking import generate_parking_data
from iot_utils import generate_rsa_key_pair, sign_data, encrypt_data, save_keys

DEVICE_GENERATORS = {
    "parking": generate_parking_data
}

def send_data_loop(device_index, device_type, device_id, frequency, backend_url, sin_mode=False):
    # Choose the generator function based on device type
    generator = DEVICE_GENERATORS[device_type]
    if generator is None:
        print(f"Unknown device type: {device_type}")
        return

    # Generate keys for the device
    private_key, public_key = generate_rsa_key_pair()
    save_keys(device_id, private_key, public_key)

    # Register the device with the backend
    public_pem = public_key.public_bytes(
        encoding="PEM",
        format="SubjectPublicKeyInfo"
    ).decode()

    response = requests.post(f"{backend_url}/api/v1/send", json={
        "device_type": device_type,
        "region": "us-east-1",
        "public_key": public_pem
    })

    if response.status_code != 200:
        print(f"[Device {device_index}] Registration failed: {response.status_code} - {response.text}")
        return
    
    # Get the device ID and backend public key from the response
    device_id = response.json().get("device_id")
    backend_public_key = response.json().get("backend_public_key")
    print(f"[Device {device_index}] Registered with Id: {device_id} - backend public key: {backend_public_key}")

    # Start sending data
    while True:
        # Generate data
        device_id = uuid.uuid4().hex
        data_id = uuid.uuid4().hex
        payload = generator(device_id)
        if payload is None:
            print(f"[Device {device_index}] Data generation failed")
            break

        # Prepare the data package
        json_data = json.dumps(payload)
        # signature = sign_data(json_data.encode(), private_key)
        # encrypted_data = encrypt_data(json_data.encode(), backend_public_key)
        data_hash = hashlib.sha256(json_data.encode()).hexdigest()

        package = {
            "device_id": device_id,
            "data": payload,
            "hash": data_hash,
            "data_id": data_id
            # "data": encrypted_data.decode(),
            # "signature": signature.decode(),
        }

        # Send the data to the backend
        try:
            res = requests.post(backend_url, json=package)
            if res.status_code != 200:
                print(f"[Device {device_index}] Failed to send data: {res.status_code} - {res.text}")
        except Exception as e:
            print(f"[Device {device_index}] Failed to send: {e}")

        # Sleep for the specified frequency
        if sin_mode:
            t = time.time()
            freq = frequency + math.sin(t) * (frequency / 2)
            sleep_time = max(1, 1 / freq)
        else:
            sleep_time = 1 / frequency

        time.sleep(sleep_time)


if __name__ == '__main__':
    # Parse command line arguments
    parser = argparse.ArgumentParser()
    parser.add_argument("--type", required=True, choices=["parking", "weather1", "weather2"], help="Device type")
    parser.add_argument("--count", type=int, default=1, help="Number of device copies")
    parser.add_argument("--freq", type=float, default=1.0, help="Frequency of sending data")
    parser.add_argument("--sin", action='store_true', help="Use sinusoidal frequency")
    parser.add_argument("--url", default="http://localhost:8080", help="Backend URL")
    args = parser.parse_args()

    # Create threads for each device
    for i in range(args.count):
        device_id = f"{args.type}_{i}"
        t = threading.Thread(target=send_data_loop, args=(i, args.type, device_id, args.freq, args.url, args.sin))
        t.start()