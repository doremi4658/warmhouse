from flask import Flask, jsonify, request
import requests
import os
import json
from threading import Thread
from kafka import KafkaConsumer

app = Flask(__name__)

devices_state = {}


def start_kafka_consumer():
    """Kafka consumer для получения алертов"""
    try:
        kafka_broker = os.getenv('KAFKA_BROKER', 'kafka:9092')

        consumer = KafkaConsumer(
            'telemetry-alerts',
            bootstrap_servers=[kafka_broker],
            value_deserializer=lambda m: json.loads(m.decode('utf-8')),
            group_id='device-service',
            auto_offset_reset='earliest'
        )

        for message in consumer:
            alert_data = message.value
            device_id = alert_data.get('device_id')

            if device_id in devices_state:
                devices_state[device_id]['last_alert'] = alert_data
                devices_state[device_id]['status'] = 'alert'

            print(f"Received alert: {alert_data}")
    except Exception as e:
        print(f"Kafka consumer error: {e}")


kafka_thread = Thread(target=start_kafka_consumer, daemon=True)
kafka_thread.start()


@app.route('/devices', methods=['GET'])
def get_devices():
    return jsonify({
        "devices": devices_state,
        "total": len(devices_state)
    })


@app.route('/devices/<device_id>', methods=['POST'])
def register_device(device_id):
    devices_state[device_id] = {
        'id': device_id,
        'status': 'online',
        'last_seen': '2024-01-15T10:00:00Z'
    }
    return jsonify({"status": "registered", "device_id": device_id})


@app.route('/devices/<device_id>/command', methods=['POST'])
def send_command(device_id):
    if device_id not in devices_state:
        return jsonify({"error": "Device not found"}), 404

    command_data = request.json
    print(f"Command received for {device_id}: {command_data}")

    return jsonify({
        "status": "command_sent",
        "device_id": device_id,
        "command": command_data.get('command'),
        "value": command_data.get('value')
    })


@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({"status": "healthy", "service": "device"})


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)