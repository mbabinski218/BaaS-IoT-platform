import json
import time
import threading
import random
import requests
from datetime import datetime, timezone
from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
import sys

interval = 5
random_interval_enabled = False
running = True
api_url = "fake_api_url"

def generate_keys():
    private_key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
    public_key = private_key.public_key()
    private_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.TraditionalOpenSSL,
        encryption_algorithm=serialization.NoEncryption()
    )
    public_pem = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )
    return private_key, public_pem

def sign_data(private_key, data):
    return private_key.sign(
        data.encode('utf-8'),
        padding.PKCS1v15(),
        hashes.SHA256()
    )

def send_to_backend(data, private_key, public_key):
    signature = sign_data(private_key, json.dumps(data))
    headers = {
        'Content-Type': 'application/json',
        'X-Public-Key': public_key.decode('utf-8'),
        'X-Signature': signature.hex()
    }
    try:        
        response = requests.post(api_url, headers=headers, data=json.dumps(data))
        print("Sent:", data)
        print("Response:", response.status_code, response.text)
    except Exception as e:
        print("Error sending data:", e)

def update_config():
    global interval, running, api_url, random_interval_enabled
    while running:
        try:
            cmd = input().strip()
            if cmd.startswith("interval"):
                _, val = cmd.split()
                if val == "random":
                    random_interval_enabled = True
                    print("Random interval enabled")
                else:
                    random_interval_enabled = False
                    interval = int(val)
                    print(f"Updated interval to {interval} seconds")
            elif cmd.startswith("api"):
                _, api_url = cmd.split()
                print(f"API url updated to {api_url}")
            elif cmd == "stop":
                running = False
                print("Stopping...")
        except Exception as e:
            print("Command error:", e)

def simulate_data():
    return {
        "device_type": "air_quality_sensor",
        "region": "central",
        "co2_level": round(random.uniform(300, 600), 2),
        "pm2_5": round(random.uniform(10, 50), 2),
        "pm10": round(random.uniform(20, 80), 2),
        "humidity": round(random.uniform(30, 70), 2),
        "temperature": round(random.uniform(10, 35), 2),
        "timestamp": datetime.now(timezone.utc).isoformat()
    }

def loop_send():
    while running:
        data = simulate_data()
        send_to_backend(data, private_key, public_key)
        if random_interval_enabled:
            time.sleep(random.randint(1, 20))            
        else:        
            time.sleep(interval)

private_key, public_key = generate_keys()
threading.Thread(target=update_config, daemon=True).start()
loop_send()
