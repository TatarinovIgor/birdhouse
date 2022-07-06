import jwt
import requests
import uuid


def create_wallet_send_data(jwtToken):
    URL = "https://atwallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/sign-up"
    # generate session and jwt from uid
    session_id = str(uuid.uuid1())
    encoded_jwt = jwtToken
    # create headers
    headers = {'X-Auth-Token': encoded_jwt,
               'X-Session-ID': session_id}
    # make a request
    request = requests.post(URL, headers=headers)
    resp = request.json()
    return activate_wallet()


def activate_wallet(resp, encoded_jwt, session_id):
    URL = 'https://atwallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/user/platform/stellar/account'
    authorisation_key = resp["access_token"]
    headers = {
        'Authorization': authorisation_key,
        'X-Auth-Token': encoded_jwt,
        'X-Session-ID': session_id,
    }
    body = {
        'platform': "stellar",
        'type': "ccba7c71-27aa-40c3-9fe8-03db6934bc20",
        'name': "Private account"
    }
    request = requests.post(URL, headers, body)
    return request.json()


def sign_in_wallet_send_data(jwtToken):
    URL = "https://atwallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/sign-in"

    session_id = str(uuid.uuid1())
    encoded_jwt = jwtToken

    headers = {'X-Auth-Token': encoded_jwt,
               'X-Session-ID': session_id}

    request = requests.post(URL, headers=headers)
    return request.json()
