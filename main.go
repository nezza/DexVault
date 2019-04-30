package main

import (
	"github.com/binance-chain/go-sdk/keys"

	// "context"
	"encoding/json"
	"encoding/hex"
	// "errors"
	"fmt"
	// jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	// "github.com/go-chi/render"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"net/http"
	"strings"
	"os"
	"bufio"
	"flag"
	"crypto/rand"
)

func GetRequestConfig(r *http.Request) *DexVaultConfiguration {
	return r.Context().Value(ConfigurationCtxKey).(*DexVaultConfiguration)
}

func GetRequestDatastore(r *http.Request) *DexVaultDatastore {
	return r.Context().Value(DatastoreCtxKey).(*DexVaultDatastore)
}

func GetRequestUser(r *http.Request) string {
	return *r.Context().Value(NameCtxKey).(*string)
}

// Configuration structure
type DexVaultConfiguration struct {
	ListenAddr string `yaml:"listen_address"`
	// Vault parameters
	VaultToken string `yaml:"vault_token"`
	VaultAddr  string `yaml:"vault_address"`
	// TLS
	TlsEnabled     bool   `yaml:"tls_enabled"`
	TlsCertificate string `yaml:"tls_certificate"`
	TlsKey         string `yaml:"tls_key"`
	// Whitelist
	IpWhitelist bool     `yaml:"ip_whitelist"`
	Whitelist   []string `yaml:"whitelist"`
}


func newAuthToken(name string, secret string) DexVaultAuth {
	// tokenAuth := jwtauth.New("HS256", []byte(secret), nil)
	tokenAuth2 := DexVaultAuth{
		// JWTAuth: *tokenAuth,
		Name:    name,
		Secret: secret,
		Permissions: []Permission{PermissionAll,},
	}
	return tokenAuth2
}

func (b *DexVaultAuth)GetJwtAuth() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(b.Secret), nil)
}

func readSecret() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
    text = strings.Replace(text, "\n", "", -1)
    return text
}

func unseal() DexVaultDatastore {
	// Protect memory from swapping.
	Mlock()

	var secret = ""
	secret = os.Getenv("DEXVAULT_SECRET")
	if secret == "" {
		fmt.Println("DEXVAULT_SECRET environment variable not set.")
		fmt.Println("Please enter unseal secret:")
		secret = readSecret()
	} else {
		fmt.Println("Unsealing datastore using DEXVAULT_SECRET env var.")
	}
	
    contents := decryptFile("datastore.bin", secret)

    datastore := DexVaultDatastore{}
    err := json.Unmarshal(contents, &datastore)
    if err != nil {
    	panic("Failed to load datastore!")
    }
    fmt.Println("Successfully unsealed.")
    datastore.Secret = secret

    return datastore
}

func (b *DexVaultDatastore)Save() {
	fmt.Println("Updating datastore.")
	bin, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	encryptFile("datastore.bin", bin, b.Secret)
}


func commandServe() {
	// Load configuration
	cfg_data, err := ioutil.ReadFile("dexvault.conf")
	if err != nil {
		panic("Failed to read configuration data.")
	}

	var cfg DexVaultConfiguration
	err = yaml.Unmarshal(cfg_data, &cfg)
	if err != nil {
		panic("Failed to decode configuration.")
	}
	
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":1234"
	}
	// Load and unseal datastore
	datastore := unseal()

	// Configure router
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		// Attach datastore to request
		r.Use(DatastoreContext(&datastore, &cfg))

		// First check: IP whitelist
		r.Use(IPWhitelist)

		// Second check: JWT
		r.Use(Verifier(datastore.Users))
		r.Use(Authenticator)

		r.Post("/v1/address", getAddressHandler)
		r.Get("/v1/wallet/", getWalletsHandler)
		r.Post("/v1/wallet/", getWalletHandler)
		r.Post("/v1/wallet/create", createWalletHandler)
		r.Post("/v1/order/create", createOrderHandler)
		r.Post("/v1/order/cancel", cancelOrderHandler)
		r.Post("/v1/token/burn", tokenBurnHandler)
		r.Post("/v1/token/freeze", freezeTokenHandler)
		r.Post("/v1/token/unfreeze", unfreezeTokenHandler)
		r.Post("/v1/token/issue", issueTokenHandler)
		r.Post("/v1/token/mint", mintTokenHandler)
		r.Post("/v1/token/send", sendTokenHandler)
		r.Post("/v1/listPair", listPairHandler)
		r.Post("/v1/proposal/submit", submitProposalHandler)
		r.Post("/v1/proposal/vote", voteProposalHandler)
		r.Post("/v1/deposit/", depositHandler)
	})

	fmt.Println("Starting server on: " + cfg.ListenAddr)
	if cfg.TlsEnabled {
		err = http.ListenAndServeTLS(cfg.ListenAddr, cfg.TlsCertificate, cfg.TlsKey, r)
	} else {
		err = http.ListenAndServe(cfg.ListenAddr, r)
	}
	fmt.Println("Server quit: ")
	fmt.Println(err)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func commandInit() {
	if _, err := os.Stat("datastore.bin"); os.IsNotExist(err) {
		fmt.Println("Welcome to the initial DexVault configuration.")
		fmt.Println("Please enter a password to use as sealing secret:")
		secret := readSecret()
		datastore := &DexVaultDatastore{}
		datastore.Secret = secret
		datastore.Save()
	} else {
		fmt.Println("datastore.bin already exists. Cancelling init.")
	}		
}

func existingUser(datastore *DexVaultDatastore, u string) *DexVaultAuth {
	if u == "" {
		fmt.Println("No user supplied.")
		os.Exit(1)
	}
	user := datastore.GetUser(u)
	if user == nil {
		fmt.Println("User not found.")
		os.Exit(1)
	}
	return user
}



func main() {
	command := flag.String("command", "", "The command to run. Example: serve, init, create-user.")
	name := flag.String("name", "", "Username to create/modify")
	permission := flag.String("permission", "", "A permission to add/revoke")
	wallet := flag.String("wallet", "", "Wallet to work on")
	flag.Parse()


	if *command == "" {
		fmt.Println("Command required.")
		return
	}
	switch *command {
	case "serve":
		commandServe()
	case "init":
		commandInit()
	}

	// User management
	if *command == "create-user" {
		datastore := unseal()

		if *name == "" {
			fmt.Println("Username required.")
			return
		}
		if datastore.GetUser(*name) != nil {
			fmt.Println("User with name already exists.")
			return
		}
		fmt.Println("Creating user: " + *name)
		bytes, err := GenerateRandomBytes(20)
		if err != nil {
			panic("Failed to generate random!")
		}
		s := hex.EncodeToString(bytes)
		fmt.Println("JWT Secret for user: " + s)

		u := DexVaultAuth {
			Name: *name,
			Secret: s,
			Permissions: []Permission{},
		}
		datastore.CreateUser(&u)
	}
	if *command == "delete-user" {
		datastore := unseal()
		_ = existingUser(&datastore, *name)
		fmt.Println("Deleting user: " + *name)
		err := datastore.DeleteUser(*name)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *command == "get-users" {
		datastore := unseal()
		fmt.Println("Users:")
		for _, u := range datastore.Users {
			fmt.Println("- " + u.Name)
		}
	}
	if *command == "get-user" {
		datastore := unseal()
		user := existingUser(&datastore, *name)
		fmt.Println("User: " + user.Name)
		fmt.Print("Permissions: ")
		fmt.Println(user.Permissions)
	}
	if *command == "add-permission" {
		datastore := unseal()
		user := existingUser(&datastore, *name)
		if *permission == "" {
			fmt.Println("No permission supplied.")
			return
		}
		user.AddPermission(Permission(*permission))
		datastore.Save()
	}
	if *command == "revoke-permission" {
		datastore := unseal()
		user := existingUser(&datastore, *name)
		if *permission == "" {
			fmt.Println("No permission supplied.")
			return
		}
		user.RevokePermission(Permission(*permission))
		datastore.Save()
	}

	// Wallet management
	if *command == "create-wallet" {
		datastore := unseal()
		if *wallet == "" {
			fmt.Println("Wallet name required.")
			return
		}

		manager, err := keys.NewKeyManager()
		if err != nil {
			fmt.Println("Failed to generate seed.")
			return
		}

		mnemonic, err := manager.ExportAsMnemonic()
		if err != nil {
			fmt.Println("Failed to acquire mnemonic.")
			return
		}

		old_w := datastore.GetWallet(*wallet)
		if old_w != nil {
			fmt.Println("Wallet with name already exists.")
			return
		}
		w := Wallet {
			Name: *wallet,
			Seed: mnemonic,
		}
		datastore.Wallets = append(datastore.Wallets, w)
		datastore.Save()
		addr := manager.GetAddr().String()
		fmt.Println("New wallet generated: " + addr)
		for {
			fmt.Println("Do you want to display the seed? (YES/NO)")
			t := readSecret()
			if t == "YES" {
				fmt.Println("Seed: " + mnemonic)
				break
			} else if t == "NO" {
				break
			}
		}
	}
	if *command == "get-wallets" {
		datastore := unseal()
		fmt.Println("Wallets:")
		for _, w := range datastore.Wallets {
			addr, err := w.GetAddress()
			if err != nil {
				fmt.Println("Failed to retrieve wallet address.")
				continue
			}
			fmt.Println("- " + w.Name + *addr)
		}
	}
	if *command == "export-wallet" {
		datastore := unseal()
		fmt.Println("ARE YOU SURE? THIS WILL DISPLAY YOUR SEED.")
		fmt.Println("Type 'YES I KNOW' to continue.")
		text := readSecret()
		if text == "YES I KNOW" {
			w := datastore.GetWallet(*wallet)
			if w == nil {
				fmt.Println("Wallet not found.")
				return
			}
			fmt.Println("Seed: " + w.Seed)
		} else {
			fmt.Println("Cancelled.")
		}
	}
	if *command == "delete-wallet" {
		datastore := unseal()
		fmt.Println("ARE YOU SURE? THE KEY WILL BE GONE.")
		fmt.Println("Type 'YES I KNOW' to continue.")
		text := readSecret()
		if text == "YES I KNOW" {
			w := datastore.GetWallet(*wallet)
			if w == nil {
				fmt.Println("Wallet not found.")
				return
			}
			datastore.DeleteWallet(w.Name)
			fmt.Println("Wallet deleted.")
		} else {
			fmt.Println("Cancelled.")
		}
	}
	if *command == "import-wallet" {
		datastore := unseal()
		if *wallet == "" {
			fmt.Println("Wallet name required.")
			return
		}

		fmt.Println("Please enter the mnemonic seed:")
		seed := readSecret()
		manager, err := keys.NewMnemonicKeyManager(seed)
		if err != nil {
			fmt.Println("Failed to parse mnemonic.")
			return
		}

		old_w := datastore.GetWallet(*wallet)
		if old_w != nil {
			fmt.Println("Wallet with name already exists.")
			return
		}
		w := Wallet {
			Name: *wallet,
			Seed: seed,
		}
		datastore.Wallets = append(datastore.Wallets, w)
		datastore.Save()
		addr := manager.GetAddr().String()
		fmt.Println("New wallet imported: " + addr)
	}
}
