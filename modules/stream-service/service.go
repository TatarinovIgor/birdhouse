package stream_service

import (
	"context"
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/operations"
	"time"
)

func ListenAndServe() {
	client := horizonclient.DefaultTestNetClient
	opRequest := horizonclient.OperationRequest{ForAccount: "GCAESSCCYR6NYLLNBUMJ24N42QGMNFTD3KWGB6SXCHZUUQIEMDYA6456", Cursor: "now"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(10 * time.Second)
		cancel()
	}()

	printHandler := func(op operations.Operation) {
		fmt.Println(op)
	}
	err := client.StreamPayments(ctx, opRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
}
