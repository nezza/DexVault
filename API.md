# API

## The endpoints

### POST

- [/v1/address](###-/v1/address)
- [/v1/wallet/create](###-/v1/wallet/create)
- [/v1/order/create](###-/v1/order/create)
- [/v1/order/cancel](###-/v1/order/cancel)
- [/v1/token/burn](###-/v1/token/burn)
- [/v1/token/freeze](###-/v1/token/freeze)
- [/v1/token/unfreeze](###-/v1/token/unfreeze)
- [/v1/token/issue](###-/v1/token/issue)
- [/v1/token/mint](###-/v1/token/mint)
- [/v1/token/send](###-/v1/token/send)
- [/v1/listPair](###-/v1/listPair)
- [/v1/proposal/submit](###-/v1/proposal/submit)
- [/v1/proposal/vote](###-/v1/proposal/vote)
- [/v1/deposit/](###-/v1/deposit/)


### /v1/address

Method: `POST`

Payload:
```
{
	"Wallet": "walletname"
}
```

Response:
```
{
	"Response": "tbnb1mrk0c5q485px083l2vakjhq8pfur8pzh2n8hce"
}
```

### /v1/wallet/create

Method: `POST`

Payload:
```
{
	"Wallet": "walletname"
}
```

Response: Address of newly created wallet
```
{
	"Response": "tbnb1mrk0c5q485px083l2vakjhq8pfur8pzh2n8hce"
}
```

### /v1/order/create

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"BaseAssetSymbol": "BNB",
	"QuoteAssetSymbol": "BTC",
	"Op": 1,
	"Price": 1000,
	"Quantity": 1000
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```


### /v1/order/cancel

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"BaseAssetSymbol": "BNB",
	"QuoteAssetSymbol": "BTC",
	"RefId": "ORDER ID"
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```


### /v1/token/burn

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"Symbol": "BNB",
	"Amount": 1234
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```

### /v1/token/freeze

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"Symbol": "BNB",
	"Amount": 1234
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```

### /v1/token/unfreeze

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"Symbol": "BNB",
	"Amount": 1234
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```

### /v1/token/issue

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"Name": "Tokenname",
	"Supply": 1234,
	"Mintable": true
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```

### /v1/token/mint

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"Symbol": "BNB",
	"Amount": 1234
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```

### /v1/token/send

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"Transfers": [
		{
			"ToAddr": "tbnb1mrk0c5q485px083l2vakjhq8pfur8pzh2n8hce",
			"Coins": [
				"symbol": "BNB",
				"free": "0.000",
				"locked": "0.000",
				"frozen": "0.0000"
			]
		}
	]
}
```

type Transfer struct {
	ToAddr types.AccAddress
	Coins  types.Coins
}
Response:
```
{
	"Response": "HEX SIGNATURE"
}
```

### /v1/listPair

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"ProposalID": 123456,
	"BaseAssetSymbol": "BNB",
	"QuoteAssetSymbol": "BTC",
	"InitPrice": 1234
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```

### /v1/proposal/submit

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"Title": "Proposal title",
	"Description": "Proposal description",
	"ProposalType": 1,
	"InitialDeposit": 1234,
	"VotingPeriod": 1234
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```


### /v1/proposal/vote

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"ProposalID": 123456,
	"Option": 1
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```

### /v1/deposit/

Method: `POST`

Payload:
```
{
	"Wallet": "walletname",
	"ChainId": "ChainId",
	"AccountNumber": 1234,
	"Sequence": 123,
	"ProposalID": 123456,
	"Amount": 1234
}
```

Response:
```
{
	"Response": "HEX SIGNATURE"
}
```
