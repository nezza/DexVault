# DexVault - Binance Dex Signing Oracle

## Introduction

DexVault provides an easy-to-use RESTful API for signing API requests and for safekeeping the wallet credentials.

Main Features:
- Simple to install, backup and use
- Supports all functionality of the official SDK
- Fully signed API requests
- Secure TLS build-in
- Full data-at-rest protection, incl. Swap prevention
- IP whitelisting
- API almost identical to Binance's Go-SDK, making migration trivial.
- Based on official Go-SDK for guaranteed compatibility and future support.

To be supported soon:
- High availability
- Hardware Security Module support via Vault
- Hashicorp Vault integration for high-scalability and high security

## Warning

This code has not yet been reviewed for security. Here is hoping that we get some funds for financing a review once the service is stable!

## Quickstart - less than 5 minutes!

```
# Initialize dexvault
dexvault -command init

# Create a user
dexvault -command create-user --name MainUser
# Make sure to copy the JWT

# Give the user ALL permissions
dexvault -command add-permission --name MainUser --permission PermissionAll

# Create a wallet
dexvault -command create-wallet --wallet Mainwallet

# Start the server!
dexvault -command serve
```

Querying the server:
```
# Requesting an address
./examples/python/get_address.py JWT_SECRET_FROM_ABOVE Mainwallet
```

And that's all!

## Installation

### Building
```
# Clone dexvault into your local Go tree:
git clone github.com/nezza/dexvault.git
cd dexvault
# Download dependencies
go get ./...
# Build
go build
```

## [API Documentation](API.md)

## Security

### TLS

DexVault has trivial-to-configure, highly-secure TLS supoprt that achives an `A+` grade on the Qualys SSL Labs test.

### Permission management

DexVault provides detailed permissions for each user, ensuring adequate privilege separation.

### Data-at-rest protection

DexVault stores its wallet data using basic 256-bit AES GCM encryption. When DexVault is started, the password to unlock the datastore has to be either typed in or must be provided in the `DEXVAULT_SECRET` environment variable.

Please note that the data-at-rest protection only provides limited protection, the service needs to still be hosted on a secure machine.

### IP Whitelisting

DexVault supports IP whitelisting, ensuring that only certain machines are able to access the API.

### API Authentication

The API uses JWT (JSON Web Tokens) for authentication. All API requests are encapsulated into the claims field of the JWT, ensuring that they are fully signed.

The JWT can be supplied in the query-string or as `Authorization` header.

Simple Python example:

```
HOST = "http://127.0.0.1:1234"
SECRET = "randomly-generated-secret"

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
```


## Command-line interface

Note: All functions except init require the unseal password. It can either be entered interactively, or be provided in the DEXVAULT_SECRET environment variable.

Initialize datastore (required before first use):
```
$ dexvault -command init
```

Start-server:
```
$ dexvault -command serve
```

### User management
Create user:
```
$ dexvault -command create-user --name username
```

List users:
```
$ dexvault -command get-users
```

List single user:
```
$ dexvault -command get-user --name username
```

Delete user:
```
$ dexvault -command delete-user --name username
```

Give permission to user (see PERMISSIONS):
```
$ dexvault -command add-permission --name username --permission PermissionAll
```

Revoke permission:
```
$ dexvault -command revoke-permission --name username --permission PermissionAll
```

### Wallet management

Create wallet with locally generated key:
```
$ dexvault -command create-wallet --wallet Testwallet
```

Get wallets:
```
$ dexvault -command get-wallets
```

Export wallet (dangerous):
```
$ dexvault -command export-wallet --wallet Testwallet
```

Delete wallet (dangerous):
```
$ dexvault -command delete-wallet --wallet Testwallet
```

Import wallet using mnemonic:
```
$ dexvault -command import-wallet --wallet Testwallet
```

## Configuration

The configuration is in `yaml` format . The following options are currently supported:

- `listen_address` - `string` - The IP + port where DexVault should be listening. Defaults to `:1234`
- `tls_enabled` - `bool` - Whether TLS should be enabled. Defaults to: `false`
- `tls_certificate` - `string` - The path of the certificate that should be used for TLS. Defaults to: ""
- `tls_key` - `string` - The path of the key that should be used for TLS. Defaults to: ""

- `ip_whitelist` - `bool` - Whether the IP whitelist should be enabled. Defaults to: `false`
- `whitelist` - `string array` - The list of IPs that are whitelisted. Defaults to: []

Example configuration:

```
tls_enabled: true
tls_certificate: server.crt
tls_key: server.key

ip_whitelist: true
whitelist: ["192.168.1.100", "192.168.1.101"]
```

## Permissions

Available permissions:
- PermissionAll - Implies ALL permissions
- PermissionRead - Read data (such as wallet addresses, but no 'secret' data)
- PermissionCreateWallet - Allows to create wallets
- PermissionCreateOrder - Allows to create orders
- PermissionCancelOrder - Allows to cancel orders
- PermissionTokenBurn - Allows to burn tokens
- PermissionDeposit - Allows to sign deposit messages
- PermissionFreezeToken - Allows to sign freeze token messages
- PermissionIssueToken - Allows to sign issue token messages
- PermissionListPair - Allows to sign list pair message
- PermissionMintTokens - Allows to sign mint token messages
- PermissionSendToken  - Allows to sign send token messages
- PermissionSubmitProposal - Allows to sign submit messages
- PermissionUnfreezeToken - Allows to sign unfreeze token messages
- PermissionVoteProposal - Allows to sign vote proposal messages