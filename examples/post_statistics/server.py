from flask import Flask, request, render_template
from flask_socketio import SocketIO

APP = Flask(__name__)
SOCKET = SocketIO(APP, async_mode='eventlet')


@APP.route('/')
def index():
    return render_template('index.html')


@APP.route('/update', methods=['POST'])
def update():
    payload = request.get_json()
    print(payload)
    SOCKET.emit('update', payload)
    return '', 200


if __name__ == '__main__':
    SOCKET.run(APP, debug=True)
