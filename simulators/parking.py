from datetime import datetime
import random

def generate_parking_data():
    return {
        "parking_lot_id": random.randint(1, 100),
        "name": "TestParkinglot",
        "address": "TestStreet",
        "city": "TestCity",
        "country": "Test",
        "latitude": round(random.uniform(50.0, 51.0), 6),
        "longitude": round(random.uniform(11.0, 12.0), 6),
        "general_occupied": random.randint(0, 10),
        "general_total": 13,
        "family_occupied": random.randint(0, 1),
        "family_total": 1,
        "disabled_occupied": random.randint(0, 12),
        "disabled_total": 12,
        "electrocharger_occupied": random.randint(0, 5),
        "electrocharger_total": 5,
        "time": datetime.now().isoformat(),
    }
