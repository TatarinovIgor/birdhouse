from flask import Flask, request, jsonify

from app.modules.create_wallet import create_wallet

app = Flask(__name__)

key = "cRfUjXn2r5u8x!A%D*G-KaPdSgVkYp3s6v9y$B?E(H+MbQeThWmZq4t7w!z%C*F)J@NcRfUjXn2r5u8x/A?D(G+KaPdSgVkYp3s6v9y$B&E)H@McQfThWmZq4t7w!z%C*F-JaNdRgUkXn2r5u8x/A?D(G+KbPeShVmYq3s6v9y$B&E)H@McQfTjWnZr4u7w!z%C*F-JaNdRgUkXp2s5v8y/A?D(G+KbPeShVmYq3t6w9z$C&E)H@McQfTjWnZr4u7x!A%D*G-JaNdRgUkXp2s5v8y/B?E(H+MbPeShVmYq3t6w9z$C&F)J@NcRfTjWnZr4u7x!A%D*G-KaPdSgVkXp2s5v8y/B?E(H+MbQeThWmZq3t6w9z$C&F)J@NcRfUjXn2r5u7x!A%D*G-KaPdSgVkYp3s6v9y/B?E(H+MbQeThWmZq4t7w!z%C&F)J@NcRfUjXn2r5u8x/A?D(G-KaPdSgVkYp3s6v9y$B&E)H@MbQeThWmZq4t7w!z%C*F-Ja"


@app.route('/', methods=['GET'])
def home():
    return "<h1>Distant Reading Archive</h1><p>This site is a prototype API for distant reading of science fiction novels.</p>"


@app.route('/create_wallet', methods=['GET'])
def create_wallet_request():
    # request headers
    secret_key = request.headers.get('auth_key')

    # request body
    uid = request.args.get('uid')
    first_name = request.args.get("first_name")
    last_name = request.args.get("last_name")
    email = request.args.get("email")
    phone = request.args.get("phone")

    # authentication (key validation)
    first_name_response, last_name_response, email_response, phone_response = create_wallet(uid, first_name, last_name, email, phone)

    if secret_key == key:  # on production key should be defined in environment
        return jsonify(
            first_name=first_name_response,
            last_name=last_name_response,
            email=email_response,
            phone=phone_response,
        )
    return "<p>Key is invalid</p>"


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
