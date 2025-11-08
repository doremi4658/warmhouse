from flask import Flask, request, jsonify
import random

app = Flask(__name__)


@app.route('/temperature', methods=['GET'])
def get_temperature():
    location = request.args.get('location', 'unknown')
    temperature = round(random.uniform(-30, 40), 1)

    return jsonify({
        'location': location,
        'temperature': temperature,
        'unit': 'celsius'
    })


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081)