package routing

import (
	"birdhouse/modules/handler"
	"birdhouse/modules/service"
	"github.com/julienschmidt/httprouter"
)

func InitRouter(router *httprouter.Router, pathName string, atWallet *service.ATWalletService) {

	routerWrap := NewRouterWrap(pathName, router)

	//GET routers
	routerWrap.POST("/create_wallet_bh", handler.MakeCreateSignUPWithWallet(atWallet))
	routerWrap.POST("/sign_in_wallet_bh", handler.MakeSignInWalletBH(atWallet))
	routerWrap.GET("/get_balance", handler.GetBalance)
	routerWrap.GET("/deposit", handler.Deposit)
	routerWrap.GET("/withdraw", handler.Withdraw)
	routerWrap.GET("/transfer", handler.Transfer)
	routerWrap.GET("/list_of_transaction", handler.ListOfTransaction)
	// call backs for AT-Wallet
	routerWrap.GET("/create_wallet_at", handler.MakeSignInAT(atWallet))
	routerWrap.GET("/sign_in_at", handler.MakeSignInAT(atWallet))
}
