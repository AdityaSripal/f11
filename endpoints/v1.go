package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/greg-szabo/f11/defaults"
	"net/http"
	//clientrpc "github.com/tendermint/tendermint/client/rpc"
	//"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	fromPrivateKey = "abc"
)

//var rpcClient clientrpc.Client

func InitializeV1() {
	//	rpcClient = rpcClient.NewHTTP(defaults.Node, "/websocket")
}

/*
func createContext() context.CoreContext {
	return context.CoreContext{
		ChainID:         defaults.TestnetName,
		Height:          0,
		Gas:             200000,
		TrustNode:       false,
		NodeURI:         defaults.Node,
		FromAddressName: "faucetAccount",
		AccountNumber:   0,
		Sequence:        0,
		Client:          nil, //rpcClient
		Decoder:         nil,
		AccountStore:    "acc",
	}
}
*/
type ClaimMessageV1 struct {
	Message string `json:"message"`
	Block   int64  `json:"block"`
	Hash    string `json:"hash"`
}

func ClaimHandlerV1(w http.ResponseWriter, r *http.Request) {
	defaults.Headers(w)
	/*
		// Get receiving address
		vars := mux.Vars(r)
		toBech32, ok := vars["to"]
		if !ok {
			defaults.UserError(w, errors.New("receiving address required"))
			return
		}

		// Get Hex addresses
		from, err := sdk.GetAccAddressBech32(defaults.FromKey)
		if err != nil {
			defaults.InternalError(w, err)
			return
		}

		to, err := sdk.GetAccAddressBech32(toBech32)
		if err != nil {
			defaults.UserError(w, err)
			return
		}

		// Parse coins
		coins, err := sdk.ParseCoins(defaults.Amount)
		if err != nil {
			defaults.InternalError(w, err)
			return
		}

		ctx := createContext()

		//Todo: Implement account check for enough coins from develop branch (x/bank/client/cli/sendtx.go)

		// build and sign the transaction
		msg := client.BuildMsg(from, to, coins)

		// Broadcast to Tendermint
		cdc := app.MakeCodec()

		res, err := ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, msg, cdc)
		if err != nil {
			defaults.InternalError(w, err)
			return
		}

		json.NewEncoder(w).Encode(ClaimMessageV1{"transaction committed", res.Height, res.Hash.String()})
	*/
	json.NewEncoder(w).Encode(ClaimMessageV1{"transaction committed", 0, "hash"})
}

func AddRoutesV1(r *mux.Router) {

	r.HandleFunc("/v1/claim/{to}", ClaimHandlerV1)

}
