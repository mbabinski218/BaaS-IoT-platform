import threading
import requests
import argparse
import hashlib
import cbor2
import math
import time
import uuid
from parking import generate_parking_data
from weather import generate_weather_data

DEVICE_GENERATORS = {
    "parking": generate_parking_data,
    "weather": generate_weather_data
}

def send_data_loop(device_index, device_type, device_id, frequency, backend_url, sin_mode=False):
    # Choose the generator function based on device type
    generator = DEVICE_GENERATORS[device_type]
    if generator is None:
        print(f"Unknown device type: {device_type}")
        return

    # Start sending data
    while True:
        # Generate data
        device_id = uuid.uuid4().__str__()
        data_id = uuid.uuid4().__str__()
        payload = generator()
        if payload is None:
            print(f"[Device {device_index}] Data generation failed")
            break

        # Prepare the data package
        cbor_bytes = cbor2.dumps(payload, canonical=True)
        hash = hashlib.sha256(cbor_bytes).hexdigest()

        package = {
            "device_id": device_id,
            "data": payload,
            "hash": hash,
            "data_id": data_id
        }

        print(package)

        # Send the data to the backend
        try:
            res = requests.post(f"{backend_url}/api/v1/send", json=package)
            if res.status_code != 201:
                print(f"[Device {device_index}] Failed to send data: {res.status_code} - {res.text}")
            if res.status_code == 201:
                print(f"[Device {device_index}] Data sent successfully: {res.json()}")
        except Exception as e:
            print(f"[Device {device_index}] Failed to send: {e}")

        # Sleep for the specified frequency
        if sin_mode:
            t = time.time()
            freq = frequency + math.sin(t) * (frequency / 2)
            sleep_time = max(1, 1 / freq)
        else:
            sleep_time = frequency

        time.sleep(sleep_time)


if __name__ == '__main__':
    # Parse command line arguments
    parser = argparse.ArgumentParser()
    parser.add_argument("--type", required=True, choices=["parking", "weather", "weather2"], help="Device type")
    parser.add_argument("--count", type=int, default=1, help="Number of device copies")
    parser.add_argument("--freq", type=float, default=30.0, help="Frequency of sending data")
    parser.add_argument("--sin", action='store_true', help="Use sinusoidal frequency")
    parser.add_argument("--url", default="http://localhost:8080", help="Backend URL")
    args = parser.parse_args()

    # Create threads for each device
    for i in range(args.count):
        device_id = f"{args.type}_{i}"
        t = threading.Thread(target=send_data_loop, args=(i, args.type, device_id, args.freq, args.url, args.sin))
        t.start()
        time.sleep(args.freq / args.count)