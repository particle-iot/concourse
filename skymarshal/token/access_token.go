package token

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/concourse/concourse/atc"

	"code.cloudfoundry.org/lager/v3"
	"github.com/concourse/concourse/atc/db"
	"github.com/go-jose/go-jose/v3/jwt"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Generator
type Generator interface {
	GenerateAccessToken(claims db.Claims) (string, error)
}

//counterfeiter:generate . Parser
type Parser interface {
	ParseExpiry(raw string) (time.Time, error)
}

//counterfeiter:generate . ClaimsParser
type ClaimsParser interface {
	ParseClaims(idToken string) (db.Claims, error)
}

func StoreAccessToken(
	logger lager.Logger,
	handler http.Handler,
	generator Generator,
	claimsParser ClaimsParser,
	accessTokenFactory db.AccessTokenFactory,
	userFactory db.UserFactory,
	displayUserIdGenerator atc.DisplayUserIdGenerator,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sky/issuer/token" {
			handler.ServeHTTP(w, r)
			return
		}
		logger := logger.Session("token-request")
		logger.Debug("start")
		defer logger.Debug("end")

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, r)

		var body io.Reader
		defer func() {
			copyResponseHeaders(w, rec.Result())
			if body != nil {
				io.Copy(w, body)
			}
		}()
		if rec.Code < 200 || rec.Code > 299 {
			body = rec.Body
			return
		}
		var resp struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token,omitempty"`
			IDToken      string `json:"id_token"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		if err != nil {
			logger.Error("unmarshal-response-from-dex", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		claims, err := claimsParser.ParseClaims(resp.IDToken)
		if err != nil {
			logger.Error("parse-id-token", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.AccessToken, err = generator.GenerateAccessToken(claims)
		if err != nil {
			logger.Error("generate-access-token", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = accessTokenFactory.CreateAccessToken(resp.AccessToken, claims)
		if err != nil {
			logger.Error("create-access-token-in-db", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		username := displayUserIdGenerator.DisplayUserId(
			claims.Connector,
			claims.UserID,
			claims.Username,
			claims.PreferredUsername,
			claims.Email,
		)
		err = userFactory.CreateOrUpdateUser(username, claims.Connector, claims.Subject)
		if err != nil {
			logger.Error("create-or-update-user", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newResp, err := json.Marshal(resp)
		if err != nil {
			logger.Error("marshal-new-response", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body = bytes.NewReader(newResp)
	})
}

func copyResponseHeaders(w http.ResponseWriter, res *http.Response) {
	for k, v := range res.Header {
		k = http.CanonicalHeaderKey(k)
		if k != "Content-Length" {
			w.Header()[k] = v
		}
	}
	w.WriteHeader(res.StatusCode)
}

func NewClaimsParser() ClaimsParser {
	return claimsParserNoVerify{}
}

type claimsParserNoVerify struct {
}

func (claimsParserNoVerify) ParseClaims(idToken string) (db.Claims, error) {
	token, err := jwt.ParseSigned(idToken)
	if err != nil {
		return db.Claims{}, err
	}

	var claims db.Claims
	err = token.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return db.Claims{}, err
	}
	return claims, nil
}

type Factory struct {
}

// GenerateAccessToken generates a token with 20 bytes of entropy with the
// unix timestamp appended.
func (Factory) GenerateAccessToken(claims db.Claims) (string, error) {
	b := [28]byte{}
	_, err := rand.Read(b[:20])
	if err != nil {
		return "", err
	}
	if claims.Expiry == nil {
		return "", errors.New("missing 'exp' claim")
	}
	binary.LittleEndian.PutUint64(b[20:], uint64(*claims.Expiry))
	return base64.RawStdEncoding.EncodeToString(b[:]), nil
}

func (Factory) ParseExpiry(accessToken string) (time.Time, error) {
	raw, err := base64.RawStdEncoding.DecodeString(accessToken)
	if err != nil {
		return time.Time{}, err
	}
	if len(raw) != 28 {
		return time.Time{}, errors.New("invalid access token length")
	}
	expiry := jwt.NumericDate(binary.LittleEndian.Uint64(raw[20:]))
	return expiry.Time(), nil
}
