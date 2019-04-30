package main

import (
	"fmt"
	sdk "github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/keys"	
	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/tx"
	jwt "github.com/dgrijalva/jwt-go"
    "net/http"
    "io/ioutil"
    "encoding/json"
)


type BasicMessage struct {
	Wallet string
}

type SignedMessage struct {
	BasicMessage
	ChainId       string
	AccountNumber int64
	Sequence      int64
}

type CreateOrder struct {
	SignedMessage
	BaseAssetSymbol  string
	QuoteAssetSymbol string
	Op               int8
	Price            int64
	Quantity         int64
}

type Response struct {
	Response string
}


func broadcastTx(tx []byte) ([]tx.TxCommitResult, error) {
	keyManager, _ := keys.NewKeyManager()
	client, err := sdk.NewDexClient("testnet-dex.binance.org", types.TestNetwork, keyManager)
	if err != nil {
		fmt.Println("Error: ")
		fmt.Println(err)
		panic(err)
	}

	param := map[string]string{}
	param["sync"] = "true"
	commits, err := client.PostTx([]byte(tx), param)
	return commits, err
}

func main() {

	secret := "YOURSECRET"
	url := "http://127.0.0.1:1234/v1/order/create"

	// Create Order
	co := CreateOrder{
		SignedMessage: SignedMessage {
			BasicMessage: BasicMessage {
				Wallet: "TESTWALLET",
			},
			ChainId: "Binance-Chain-Nile",
			AccountNumber: 667929,
			Sequence: 18,
		},
		BaseAssetSymbol: "ANN-457",
		QuoteAssetSymbol: "BNB",
		Op: 1,
		Price: 100000000,
		Quantity: 100000000,
	}
	_ = co

	// Marshal our payload
	payload, err := json.Marshal(co)
	if err != nil {
		panic(err)
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"payload": string(payload),
		})
	tokenString, err := token.SignedString([]byte(secret))

	// Send request to DexVault
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "BEARER " + tokenString)
	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}

	// Read response
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Raw response: " + string(body))
	// Unmarshal response
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	// Broadcast transaction
	commits, err := broadcastTx([]byte(response.Response))
	if err != nil {
		panic(err)
	}

	fmt.Println("Order submitted!")
	fmt.Println(commits)
}