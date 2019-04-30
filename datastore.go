package main

import (
	"errors"
	"fmt"
	"github.com/binance-chain/go-sdk/keys"
)

// JWT Authentication struct (User)
type DexVaultAuth struct {
	// jwtauth.JWTAuth
	Name        string
	Secret      string
	Permissions []Permission
}

type Wallet struct {
	Name string
	Seed string
}

type DexVaultDatastore struct {
	Secret  string `json:"-"`
	Wallets []Wallet
	Users   []*DexVaultAuth
}

func (b *DexVaultDatastore) CreateWallet(wallet string) (*Wallet, error) {
	fmt.Println("Creating new wallet: " + wallet)
	old_w := b.GetWallet(wallet)
	if old_w != nil {
		fmt.Println("Wallet with name already exists.")
		return nil, errors.New("Wallet with name already exists.")
	}

	newKey, err := keys.NewKeyManager()
	if err != nil {
		fmt.Println("Key generation failed:")
		fmt.Println(err)
		return nil, err
	}
	mnemonic, err := newKey.ExportAsMnemonic()
	if err != nil {
		fmt.Println("Mnemonic export failed:")
		fmt.Println(err)
		return nil, err
	}

	w := Wallet{
		Name: wallet,
		Seed: mnemonic,
	}

	b.Wallets = append(b.Wallets, w)
	b.Save()
	return &w, nil
}

func (u *DexVaultAuth) HasPermission(p Permission) bool {
	for _, per := range u.Permissions {
		if per == PermissionAll {
			return true
		}
		if per == p {
			return true
		}
	}
	return false
}

func (u *DexVaultAuth) HasSpecificPermission(p Permission) bool {
	for _, per := range u.Permissions {
		if per == p {
			return true
		}
	}
	return false
}

func (u *DexVaultAuth) AddPermission(p Permission) {
	if u.HasSpecificPermission(p) {
		return
	}
	u.Permissions = append(u.Permissions, p)
}

func (u *DexVaultAuth) RevokePermission(p Permission) {
	if !u.HasSpecificPermission(p) {
		return
	}
	for i, per := range u.Permissions {
		if per == p {
			u.Permissions = append(u.Permissions[:i], u.Permissions[i+1:]...)
		}
	}
}

func (w *Wallet) GetKeyManager() (keys.KeyManager, error) {
	return keys.NewMnemonicKeyManager(w.Seed)
}

func (w *Wallet) GetAddress() (*string, error) {
	km, err := w.GetKeyManager()
	if err != nil {
		return nil, err
	}
	straddr := km.GetAddr().String()
	return &straddr, nil
}

func (b *DexVaultDatastore) GetWallet(wallet string) *Wallet {
	for _, w := range b.Wallets {
		if w.Name == wallet {
			return &w
		}
	}
	return nil
}

func (b *DexVaultDatastore) DeleteWallet(w string) error {
	for i, wallet := range b.Wallets {
		if wallet.Name == w {
			b.Wallets = append(b.Wallets[:i], b.Wallets[i+1:]...)
			b.Save()
			return nil
		}
	}
	return errors.New("Wallet not found.")
}

func (b *DexVaultDatastore) GetUser(user string) *DexVaultAuth {
	for _, u := range b.Users {
		if u.Name == user {
			return u
		}
	}
	return nil
}

func (b *DexVaultDatastore) CreateUser(u *DexVaultAuth) {
	b.Users = append(b.Users, u)
	b.Save()
}

func (b *DexVaultDatastore) DeleteUser(u string) error {
	for i, user := range b.Users {
		if user.Name == u {
			b.Users = append(b.Users[:i], b.Users[i+1:]...)
			b.Save()
			return nil
		}
	}
	return errors.New("User not found.")
}

func (b *DexVaultDatastore) IsPermitted(user string, wallet string, action Permission) bool {
	fmt.Println("IsPermitted: " + user + "for action: " + string(action))
	u := b.GetUser(user)
	w := b.GetWallet(wallet)

	if u == nil || w == nil {
		fmt.Println("No user or wallet found.")
		return false
	}

	for _, p := range u.Permissions {
		if p == PermissionAll {
			fmt.Println("User has ALL permission.")
			return true
		}
		if p == action {
			fmt.Println("User has permission.")
			return true
		}
	}

	fmt.Println("User does NOT have permission.")

	return false
}
