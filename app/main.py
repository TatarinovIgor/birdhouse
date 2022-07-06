import jwt
from flask import Flask, request, jsonify

from app.modules.create_wallet import create_wallet_send_data

app = Flask(__name__)

key = "secret"

@app.route('/', methods=['GET'])
def home():
    return "<h1>Distant Reading Archive</h1><p>This site is a prototype API for distant reading of science fiction novels.</p>"


@app.route('/create_wallet_BH', methods=['POST'])
def create_wallet_request():
    # request headers
    # Should be JWT, with credentials and secret key
    jwtToken = request.args.get('auth_key')

    # inside should be parameters entered during registration
    decoded = jwt.decode(jwtToken, "secret", algorithms=["HS256"])


    # request body
    uid = request.args.get('uid')
    first_name = request.args.get("first_name")
    last_name = request.args.get("last_name")
    email = request.args.get("email")
    phone = request.args.get("phone_number")
    #create_wallet_send_data(jwtToken, uid, first_name, last_name, email, phone)
    if decoded["uid"] == uid and decoded["first_name"] == first_name:
        return jsonify(
            message="success"
        )
    return jsonify(
        message="invalid key or parameters"
    )


@app.route('/create_wallet_AT', methods=['GET'])
def create_wallet_requestAT():
    # request headers
    X_Auth_Token = request.headers.get('X-Auth-Token')
    X_Session_ID = request.headers.get('X-Session-ID')
    # request body

    # authentication (key validation)
    #first_name_response, last_name_response, email_response, phone_response = create_wallet(secret_key, uid, first_name,
    #                                                                                       last_name, email, phone)

    #uid = request.args.get('uid')
    #first_name = request.args.get("first_name")
    #last_name = request.args.get("last_name")
    #email = request.args.get("email")
    #phone = request.args.get("phone")

    # authentication (key validation)
    #first_name_response, last_name_response, email_response, phone_response = create_wallet(secret_key, uid, first_name,
    #                                                                                        last_name, email, phone)

    #if secret_key == key:  # on production key should be defined in environment
    #    return jsonify(
    #        first_name=first_name_response,
    #        last_name=last_name_response,
    #        email=email_response,
    #        phone=phone_response,
    #    )
    return "msg"

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
