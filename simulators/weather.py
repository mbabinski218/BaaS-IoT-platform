from datetime import datetime
import random

def generate_weather_data():
    return {
        "temperature": round(random.uniform(-10.0, 35.0), 2),
        "humidity": float(random.randint(0, 100)),
        "status": random.choice(["sunny", "cloudy", "rainy", "snowy"]),
        "timestamp": datetime.now().isoformat()
    }
