package main

import (
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/x/bank/client"
	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/greg-szabo/f11/config"
	f11context "github.com/greg-szabo/f11/context"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bech32"
	"github.com/tomasen/realip"
	"log"
	"net/http"
)

func V1ClaimHandler(ctx *f11context.Context, w http.ResponseWriter, r *http.Request) (status int, err error) {
	status = http.StatusInternalServerError

	var claim struct {
		Address  string `json:"address"`
		Response string `json:"response"`
	}

	// decode JSON response from body
	err = json.NewDecoder(r.Body).Decode(&claim)
	if err != nil {
		return
	}

	// make sure address is bech32 encoded
	hrp, decodedAddress, err := bech32.DecodeAndConvert(claim.Address)
	if err != nil {
		return
	}

	// encode the address in bech32
	encodedAddress, err := bech32.ConvertAndEncode(hrp, decodedAddress)
	if err != nil {
		return
	}

	// make sure captcha is valid
	clientIP := realip.FromRequest(r)

	if !ctx.DisableRecaptcha {
		var captchaPassed bool
		captchaPassed, err = recaptcha.Confirm(clientIP, claim.Response)
		if err != nil {
			return
		}
		if !captchaPassed {
			return status, errors.New("shoo robot, recaptcha failed")
		}
	} else {
		log.Print("Recaptcha disabled")
	}

	message := "transaction committed"
	var height int64 = 0
	hash := "SendDisabled"
	if !ctx.DisableSend {
		height, hash, status, err = V1SendTx(ctx, encodedAddress)
		if err != nil {
			return
		}

		if ctx.DbSession != nil {
			ctx.Cfg, err = config.AddSequenceInDB(ctx.DbSession, ctx.Cfg)
			if err != nil {
				log.Printf("critical error: could not update sequence number in database. Current value: %d, error: %v", ctx.Cfg.Sequence, err)
			}
		} else {
			log.Println("DynamoDB disabled, sequence number only stored in memory.")
		}
	} else {
		status = http.StatusOK
	}

	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
		Hash    string `json:"hash"`
		Height  int64  `json:"height"`
	}{
		Message: message,
		Height:  height,
		Hash:    hash,
	})
	return
}

func V1SendTx(ctx *f11context.Context, toBech32 string) (height int64, hash string, status int, err error) {
	status = http.StatusInternalServerError
	// Get Hex addresses
	from, err := sdk.AccAddressFromBech32(ctx.Cfg.AccountAddress)
	if err != nil {
		return
	}

	to, err := sdk.AccAddressFromBech32(toBech32)
	if err != nil {
		status = http.StatusBadRequest
		return
	}

	publicKey, err := sdk.GetAccPubKeyBech32(ctx.Cfg.PublicKey)
	if err != nil {
		return
	}

	// Parse coins
	coins, err := sdk.ParseCoins(ctx.Cfg.Amount)
	if err != nil {
		return
	}

	coreCtx := context.CoreContext{
		ChainID:         ctx.Cfg.TestnetName,
		Height:          0,
		Gas:             200000,
		TrustNode:       false,
		NodeURI:         ctx.Cfg.Node,
		FromAddressName: "faucetAccount",
		AccountNumber:   0,
		Sequence:        0,
		Client:          ctx.RpcClient,
		Decoder:         nil, //authcmd.GetAccountDecoder(cdc),
		AccountStore:    "acc",
	}

	//Todo: Implement account check for enough coins
	//Derive coin number from sequence number (c - s = remaining coins)

	// build the transaction
	msg := client.BuildMsg(from, to, coins)

	// No fee
	fee := sdk.Coin{}

	// There's nothing to see here, move along.
	memo := "faucet drop"

	// Message
	signMsg := auth.StdSignMsg{
		ChainID:       ctx.Cfg.TestnetName,
		AccountNumber: ctx.Cfg.AccountNumber,
		Sequence:      ctx.Cfg.Sequence,
		Msgs:          []sdk.Msg{msg},
		Memo:          memo,
		Fee:           auth.NewStdFee(coreCtx.Gas, fee),
	}
	bz := signMsg.Bytes()

	// Get private key
	privateKeyBytes, err := config.GetPrivkeyBytesFromString(ctx.Cfg.PrivateKey)
	if err != nil {
		return
	}
	privateKey, err := crypto.PrivKeyFromBytes(privateKeyBytes)

	// Sign message
	sig, err := privateKey.Sign(bz)

	sigs := []auth.StdSignature{{
		PubKey:        publicKey,
		Signature:     sig,
		AccountNumber: ctx.Cfg.AccountNumber,
		Sequence:      ctx.Cfg.Sequence,
	}}

	// marshal bytes
	tx := auth.NewStdTx(signMsg.Msgs, signMsg.Fee, sigs, memo)

	// Broadcast to Tendermint
	txBytes, err := ctx.Cdc.MarshalBinary(tx)
	if err != nil {
		return
	}
	res, err := coreCtx.BroadcastTx(txBytes)
	if err != nil {
		return
	}

	status = http.StatusOK
	height = res.Height
	hash = res.Hash.String()
	return

}