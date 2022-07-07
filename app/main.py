import jwt
from flask import Flask, request, jsonify

from app.modules.create_wallet import create_wallet_send_data, sign_in_wallet_send_data

app = Flask(__name__)
key = "secret"


@app.route('/', methods=['GET'])
def home():
    return "<h1>Distant Reading Archive</h1><p>This site is a prototype API for distant reading of science fiction novels.</p>"


@app.route('/create_wallet_bh', methods=['POST'])
def create_wallet_request():
    # request headers
    # Should be JWT, with credentials and secret key
    jwtToken = request.args.get('auth_key')

    # inside should be parameters entered during registration
    decoded = jwt.decode(jwtToken, "secret", algorithms=["HS256"])

    return jsonify(
        create_wallet_send_data(jwtToken)
    )


@app.route('/create_wallet_at', methods=['GET'])
def create_wallet_requestAT():
    jwtToken = request.headers.get('X-Auth-Token', type=str)
    sessionID = request.headers.get('X-Session-ID', type=str)

    decoded = jwt.decode(jwtToken, "secret", algorithms=["HS256"])

    return jsonify(
        external_id=decoded["external_id"],
        first_name=decoded["first_name"],
        last_name=decoded["last_name"],
        email=decoded["email"],
        phone=decoded["phone"],
    )


@app.route('/sign_in_wallet_bh', methods=['POST'])
def sign_up_wallet_BH():
    jwtToken = request.args.get('auth_key')

    return jsonify(
        sign_in_wallet_send_data(jwtToken)
    )


@app.route('/sign_in_at', methods=['GET'])
def sign_in_wallet_at():
    jwtToken = request.headers.get('X-Auth-Token', str)
    sessionId = request.headers.get('X-Session-ID', str)

    decoded = jwt.decode(jwtToken, "secret", algorithms=["HS256"])
    return {
        "external_id": decoded["external_id"]
    }


#@app.route('/get_wallet', methods=[])
@app.route('/top_up', methods=['GET'])
def top_up():
    secret_key = request.headers.get('auth_key')
    if secret_key == key:
        return "<p>Top up account here</p>"
    return "<p>Key is invalid</p>"


@app.route('/payout', methods=['GET'])
def payout():
    secret_key = request.headers.get('auth_key')
    if secret_key == key:
        return "<p>Payout here</p>"
    return "<p>Key is invalid</p>"


@app.route('/transfer', methods=['GET'])
def transfer():
    secret_key = request.headers.get('auth_key')
    if secret_key == key:
        return "<p>Transfer here</p>"
    return "<p>Key is invalid</p>"


@app.route('/get_balance', methods=['GET'])
def get_balance():
    secret_key = request.headers.get('auth_key')
    if secret_key == key:
        return "<p>Get Balance here</p>"
    return "<p>Key is invalid</p>"
