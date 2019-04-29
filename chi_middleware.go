package main

import (
	"context"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"net/http"
	"strings"
	// "fmt"
)

const NameCtxKey = "namectxkey"
const ConfigurationCtxKey = "configurationctxkey"
const DatastoreCtxKey = "datastorectxkey"

// Validate the request
func Verifier(ja []*DexVaultAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Verify(ja, jwtauth.TokenFromQuery, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie)(next)
	}
}

func Verify(ja []*DexVaultAuth, findTokenFns ...func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var token *jwt.Token
			var err error
			var name *string = nil
			for _, v := range ja {
				j := v.GetJwtAuth()
				token, err = jwtauth.VerifyRequest(j, r, findTokenFns...)
				if err == nil {
					name = &v.Name
					// If verification was successful then break
					break
				}
			}
			ctx = jwtauth.NewContext(ctx, token, err)
			ctx = context.WithValue(ctx, NameCtxKey, name)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

// Implement authenticator
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			fmt.Println("Authenticator failed:")
			fmt.Println(err)
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid {
			fmt.Println("Authenticator failed: Invalid/empty token.")
			http.Error(w, http.StatusText(401), 401)
			return
		}

		fmt.Println("Authenticator user: " + *r.Context().Value(NameCtxKey).(*string))

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// Implements a simple IP whitelist
func IPWhitelist(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := strings.Split(r.RemoteAddr, ":")[0]
		fmt.Println("Requesting IP: " + ip)

		cfg := GetRequestConfig(r)
		if cfg.IpWhitelist {
			fmt.Println("Whitelist enabled.")
			var allow = false
			for _, wip := range cfg.Whitelist {
				if wip == ip {
					allow = true
					break
				}
			}

			if !allow {
				fmt.Println("IP not found in whitelist.")
				http.Error(w, http.StatusText(401), 401)
				return
			}
			fmt.Println("IP allowed.")
		} else {
			fmt.Println("Whitelist disabled.")
		}
		next.ServeHTTP(w, r)
	})
}

// The following two functions attach the datastore and
// the config to the request.
func DatastoreContextHandler(datastore *DexVaultDatastore, config *DexVaultConfiguration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, DatastoreCtxKey, datastore)
			ctx = context.WithValue(ctx, ConfigurationCtxKey, config)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func DatastoreContext(datastore *DexVaultDatastore, config *DexVaultConfiguration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return DatastoreContextHandler(datastore, config)(next)

	}
}
