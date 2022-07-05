# importing the requests library
import requests


# api-endpoint
def create_wallet(uid, first_name, last_name, email, phone):
    URL = "https://atwallet.rock-west.net/api/v1/wallet/application/1e4e7c78-79a0-4483-9617-e985f4733481/guest"
    # generate session and jwt from uid
    headers = {'content-type': 'application/json'}
    # request = requests.get(URL)
    return first_name, last_name, email, phone
