package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func withdraw_from_wallet(jwtToken string, amount int, guid string) ([]bytes, error) {
	sign_in_data := sign_in_wallet_send_data(jwtToken)
	merchant_guid := "af74cf04-311f-435e-a530-0c85fbd6d154"
	URL := "https://atwallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/user/platform/stellar/account/"
	client := http.Client{}

	data, _ := json.Marshal(map[string]interface{}{
		"amount":        float32(amount),
		"asset":         "ATUSD:GBT4VVTDPCNA45MNWX5G6LUTLIEENSTUHDVXO2AQHAZ24KUZUPLPGJZH",
		"merchant_guid": merchant_guid,
	})

	request, err := http.NewRequest("POST", URL+guid+"/payout", bytes.NewBuffer(data))
	operandID := "sep0031Payout:650bd907-baf9-11ec-9da4-5405dbf726e9"
	SessionID := 123456789

	if err != nil {
		return nil, err
	}

	result, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if result.StatusCode != 200 {
		return nil, err
	}

	res := struct {
		Transaction string `json:"transaction"`
	}{}
	json.NewDecoder(result.Body).Decode(&res)
	resp_auth := authPayment(res.Transaction, guid, operandID, sign_in_data["access_token"], jwtToken, SessionID)
}

func authPayment(transaction string, guid string, operandID string, Authorization string, AuthToken string, SessionID int) error {
	URL := "https://atwallet.rock-west.net/api/v1/wallet/application/ab54ee14-15f1-4ce5-bcc3-6559451354da/user/platform/stellar/account/"

	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+Authorization)
	request.Header.Set("Session-ID", transaction)
	request.Header.Set("X-Auth-Token", AuthToken)
	request.Header.Set("X-SessionID", string(SessionID))

	data, _ := json.Marshal(map[string]interface{}{
		"transaction_status": "pending",
		"operation_status":   "pending",
	})

	req, err := http.NewRequest("PATCH", URL+guid+"/operation/"+operandID, bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	return req, nil
}
