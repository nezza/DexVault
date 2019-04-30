#!/usr/bin/env python

import jwt
import requests
import json
import sys

try:
	SECRET = sys.argv[1]
	NAME = sys.argv[2]
except:
	print "Usage: create_wallet.py SECRET NAME"


HOST = "http://127.0.0.1:1234"


def enc(data):
	return jwt.encode({"payload": json.dumps(data)}, SECRET, algorithm='HS256')


d = {
	"Wallet": NAME,
	"ChainId": "Binance-Chain-Nile",
	"AccountNumber": 667929,
	"Sequence": 13,
	"BaseAssetSymbol": "BNB",
	"QuoteAssetSymbol": "ANN-457",
	"Op": 1,
	"Price": 100000000,
	"Quantity": 100000000,
}



headers = {
	"Authorization": "BEARER " + enc(d),
	"Content-Type": "application/json"
}

r = requests.post(HOST + "/v1/order/create", headers=headers)

print r.text