package telegram

import (
	"birdhouse/modules/service"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

func (receiver *TelegramService) CreateJWT(config *tgbotapi.Message, email, phone string) (*service.UserAccount, error) {
	user := config.From
	userID := user.ID

	if !user.IsBot {
		userData := &service.UserData{
			ExternalID: strconv.FormatInt(userID, 10),
			FirsName:   user.FirstName,
			LastName:   user.LastName,
			Phone:      "nil",
			Email:      "nil",
		}

		payload, _ := NewPayload(*userData)

		payloadToken := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
		tokenString, err := payloadToken.SignedString(receiver.privateKey)
		if err != nil {
			return &service.UserAccount{}, err
		}
		token, err := receiver.atWallet.SignUp(tokenString)
		if err != nil {
			return &service.UserAccount{}, err
		}

		// this is only test, pub key should be removed in prod
		result, err := receiver.atWallet.CreateStellarWallet(payloadToken.Raw, token.AccessToken,
			"ccba7c71-27aa-40c3-9fe8-03db6934bc20", "BirdHouseClientAccount")

		fmt.Println(result.GUID)
		return result, nil
	}
	fmt.Println("no bots allowed here")

	return &service.UserAccount{}, nil
}

func NewPayload(userData service.UserData) (*service.TokenClaims, error) {
	payload := &service.TokenClaims{
		Payload: userData,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			IssuedAt:  time.Now().Unix(),
			Subject:   "telegram request",
		},
	}
	return payload, nil
}
