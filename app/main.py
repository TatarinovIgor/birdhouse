import flask

app = flask.Flask(__name__)


@app.route('/', methods=['GET'])
def home():
    return "<h1>Distant Reading Archive</h1><p>This site is a prototype API for distant reading of science fiction novels.</p>"


@app.route('/create_wallet', methods=['GET'])
def create_wallet():
    return "<p>AT wallet creation here</p>"


@app.route('/top_up', methods=['GET'])
def top_up():
    return "<p>Top up account here</p>"


@app.route('/payout', methods=['GET'])
def payout():
    return "<p>Payout here</p>"


@app.route('/transfer', methods=['GET'])
def transfer():
    return "<p>Transfer here</p>"


@app.route('/get_balance', methods=['GET'])
def get_balance():
    return "<p>Get Balance here</p>"


