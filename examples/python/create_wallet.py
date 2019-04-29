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
	"wallet": NAME,
}

headers = {
	"Authorization": "BEARER " + enc(d),
	"Content-Type": "application/json"
}

r = requests.post(HOST + "/v1/wallet/create", headers=headers)

print r.text