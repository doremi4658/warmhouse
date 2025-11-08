from flask import Flask, request, jsonify
import random
import os

app = Flask(__name__)


@app.route('/temperature', methods=['GET'])
def get_temperature():
    location = request.args.get('location', 'default')

    # Генерируем случайную температуру от -30 до +40
    temperature = round(random.uniform(-30, 40), 1)

    return jsonify({
        'location': location,
        'temperature': temperature,
        'unit': 'celsius'
    })


@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({'status': 'healthy'})


if __name__ == '__main__':
    port = int(os.environ.get('PORT', 8081))
    app.run(host='0.0.0.0', port=port, debug=True)