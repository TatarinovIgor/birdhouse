package routing

import (
	"birdhouse/modules/handler"
	"birdhouse/modules/middleware"
	"birdhouse/modules/service"
	"github.com/julienschmidt/httprouter"
)

func InitRouter(router *httprouter.Router, pathName string, atWallet *service.ATWalletService) {

	routerWrap := NewRouterWrap(pathName, router)

	//GET routers
	routerWrap.POST("/create_wallet_bh", middleware.AuthMiddleware(atWallet, handler.MakeCreateSignUPWithWalletBH(atWallet)))
	routerWrap.POST("/sign_in_wallet_bh", middleware.AuthMiddleware(atWallet, handler.MakeSignInWalletBH(atWallet)))
	routerWrap.GET("/deposit_wallet_link", middleware.AuthMiddleware(atWallet, handler.MakeFPFLinkForWallet(atWallet, true)))
	routerWrap.GET("/withdraw_wallet_link", middleware.AuthMiddleware(atWallet, handler.MakeFPFLinkForWallet(atWallet, false)))
	routerWrap.GET("/get_balance", middleware.AuthMiddleware(atWallet, handler.MakeGetBalance(atWallet)))
	routerWrap.GET("/deposit", middleware.AuthMiddleware(atWallet, handler.MakeFPFLinkForWallet(atWallet, true)))
	routerWrap.GET("/withdraw", middleware.AuthMiddleware(atWallet, handler.MakeFPFLinkForWallet(atWallet, false)))
	routerWrap.GET("/transfer/deposit", handler.TransferDeposit(atWallet))
	routerWrap.GET("/list_of_transaction", handler.ListOfTransaction)
	// call backs for AT-Wallet
	routerWrap.GET("/create_wallet_at", handler.MakeSignInAT(atWallet))
	routerWrap.GET("/sign_in_at", handler.MakeSignInAT(atWallet))
}
