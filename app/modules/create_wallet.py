import jwt
import requests
import uuid


def create_wallet_send_data(jwtToken, uid, first_name, last_name, email, phone):
    URL = "https://atwallet.rock-west.net/api/v1/wallet/application/1e4e7c78-79a0-4483-9617-e985f4733481/guest"
    # generate session and jwt from uid
    session_id = uuid.uuid1()
    encoded_jwt = jwtToken
    # create headers
    headers = {'X-Auth-Token': encoded_jwt,
               'X-Session-ID': session_id}
    # make a request
    request = requests.get(URL, headers=headers)


