from flask import Flask, request, jsonify
from kafka import KafkaProducer
import json
import os

app = Flask(__name__)

kafka_broker = os.getenv('KAFKA_BROKER', 'kafka:9092')
producer = KafkaProducer(
    bootstrap_servers=[kafka_broker],
    value_serializer=lambda v: json.dumps(v).encode('utf-8')
)


@app.route('/telemetry', methods=['POST'])
def receive_telemetry():
    data = request.json

    producer.send('telemetry-data', value={
        'device_id': data.get('device_id'),
        'metric': data.get('metric'),
        'value': data.get('value'),
        'timestamp': data.get('timestamp'),
        'location': data.get('location')
    })

    if data.get('metric') == 'temperature' and data.get('value', 0) > 35:
        producer.send('telemetry-alerts', value={
            'alert_type': 'high_temperature',
            'device_id': data.get('device_id'),
            'value': data.get('value'),
            'message': 'Температура превысила допустимый предел'
        })

    return jsonify({"status": "telemetry_received", "broker": "kafka"})


@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({"status": "healthy", "service": "telemetry"})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8082)