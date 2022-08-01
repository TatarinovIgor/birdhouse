import jwt
from flask import Flask, request, jsonify

from app.modules.create_wallet import create_wallet_send_data, sign_in_wallet_send_data, \
    get_wallet_balance, make_deposit

app = Flask(__name__)
key = "secret"


@app.route('/', methods=['GET'])
def home():
    return "<h1>Distant Reading Archive</h1><p>This site is a prototype API for distant reading of science fiction novels.</p>"


# ToDo - Test on the front
@app.route('/create_wallet_bh', methods=['POST'])
def create_wallet_request():
    # request headers
    # Should be JWT, with credentials and secret key
    jwtToken = request.args.get('auth_key')

    # inside should be parameters entered during registration
    decoded = jwt.decode(jwtToken, "secret", algorithms=["HS256"])

    return create_wallet_send_data(jwtToken)



# ToDo - Test on the front
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


# ToDo - Test on the front
@app.route('/sign_in_wallet_bh', methods=['POST'])
def sign_up_wallet_BH():
    jwtToken = request.args.get('auth_key')

    return sign_in_wallet_send_data(jwtToken)


# should save guid in bubble

# ToDo - Test on the front
@app.route('/sign_in_at', methods=['GET'])
def sign_in_wallet_at():
    jwtToken = request.headers.get('X-Auth-Token', str)
    sessionId = request.headers.get('X-Session-ID', str)

    decoded = jwt.decode(jwtToken, "secret", algorithms=["HS256"])
    return {
        "external_id": decoded["external_id"]
    }


# ToDo - Test on the front, Test generally
@app.route('/get_balance', methods=['GET'])
def get_balance():
    jwtToken = request.args.get('auth_key', type=str)
    guid = request.args.get('guid', type=str)
    return get_wallet_balance(jwtToken, guid)

# ToDo - Test on the front, Test generally, Realise functional
@app.route('/deposit', methods=['GET'])
def deposit():

    # request headers
    # Should be JWT, with credentials and secret key
    jwtToken = request.args.get('auth_key')
    amount = request.args.get('amount')
    acc_guid = request.args.get('acc_guid')

    return make_deposit(jwtToken, amount, acc_guid)


# ToDo  - Test on the front, Test generally
@app.route('/withdraw', methods=['GET'])
def payout():
    jwtToken = request.args.get('auth_key', type=str)
    guid = request.args.get('guid', type=str)
    amount = request.args.get('amount', type=str)
    return withdraw_from_wallet(jwtToken, amount, guid)


# ToDo  - Test on the front, Test generally, Realise functional
@app.route('/transfer', methods=['GET'])
def transfer():
    secret_key = request.headers.get('auth_key')
    if secret_key == key:
        return "<p>Transfer here</p>"
    return "<p>Key is invalid</p>"


# ToDo  - Test on the front, Test generally, Realise functional
@app.route('/list_of_transaction', methods=['GET'])
def list_of_transactions():
    secret_key = request.headers.get('auth_key')
    if secret_key == key:
        return "<p>Deposit here</p>"
    return "<p>Key is invalid</p>"
