package routing

import (
	"birdhouse/modules/handler"
	"github.com/julienschmidt/httprouter"
)

func InitRouter(router *httprouter.Router, pathName string) {

	routerWrap := NewRouterWrap(pathName, router)

	//GET routers
	routerWrap.POST("/create_wallet_bh", handler.CreateWalletBH)
	routerWrap.GET("/create_wallet_at", handler.CreateWalletAT)
	routerWrap.POST("/sign_in_wallet_bh'", handler.SignInWalletBH)
	routerWrap.GET("/sign_in_at", handler.SignInAT)
	routerWrap.GET("/get_balance", handler.GetBalance)
	routerWrap.GET("/deposit", handler.Deposit)
	routerWrap.GET("/withdraw", handler.Withdraw)
	routerWrap.GET("/transfer", handler.Transfer)
	routerWrap.GET("/list_of_transaction", handler.ListOfTransaction)
}
