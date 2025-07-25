package v0

import (
	"crypto/sha512"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/qcrg/bookish-guacamole/api/v0/token"
	"github.com/rs/zerolog"
	"github.com/tdewolff/parse/v2/buffer"
	"golang.org/x/crypto/bcrypt"
)

type token_req struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type refresh_session_req struct {
	Token token_req `json:"token"`
}

func call_webhook(log zerolog.Logger, new_ip string) {
	url := os.Getenv("BHGL_WEBHOOK")
	if len(url) == 0 {
		log.Warn().Msg("Webhook is not defined")
	}
	data := fmt.Sprintf(
		`{"new_ip": "%s","msg":"Attempt to refresh session from new ip"}`,
		new_ip,
	)
	log.Debug().Str("webhook", url).Str("new_ip", new_ip).Msg("Call webhook")
	_, err := http.Post("", "application/json", buffer.NewReader([]byte(data)))
	if err != nil {
		log.Error().Err(err).Msg("Failed to call webhook")
		return
	}
}

func refresh_current_session(ctx iris.Context, deps *Deps) {
	var req refresh_session_req
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Request is invalid"))
		deps.Log.Error().Err(err).Msg("Body must be valid JSON")
		return
	}

	access_token_id, typ, err := token.Parse(req.Token.AccessToken)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Request is invalid"))
		deps.Log.Error().Err(err).Msg("Access token is invalid")
		return
	}
	if typ != token.TypeAccess {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Request is invalid"))
		deps.Log.Error().
			Str("token_type", typ).
			Msg("Access token with incorrect type")
		return
	}

	refresh_token_id, typ, err := token.Parse(req.Token.RefreshToken)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Request is invalid"))
		deps.Log.Error().Err(err).Msg("Refresh token is invalid")
		return
	}
	if typ != token.TypeRefresh {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Request is invalid"))
		deps.Log.Error().
			Err(err).
			Str("token_type", typ).
			Msg("Refresh token with incorrect type")
		return
	}

	if access_token_id != refresh_token_id {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Request is invalid"))
		deps.Log.Error().Msg("Tokens are not a pair")
		return
	}

	token_id := access_token_id

	exists, err := deps.DB.Tokens().Exists(token_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to check existence of token")
		return
	}
	if !exists {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Token is invalid"))
		deps.Log.Debug().Str("token_id", token_id).Msg("Token not found")
		return
	}

	old_refresh_hash, err := deps.DB.Tokens().GetRefreshHash(token_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to get refresh token hash from db")
		return
	}
	refresh_sha512 := sha512.Sum512([]byte(req.Token.RefreshToken))
	err = bcrypt.CompareHashAndPassword([]byte(old_refresh_hash), refresh_sha512[:])
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Token is invalid"))
		deps.Log.Debug().
			Err(err).
			Str("token_id", token_id).
			Msg("Hash from refresh tokens is not equal")
		return
	}

	session_id, err := deps.DB.Sessions().GetIdFromTokenId(token_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to get session id from db")
		return
	}

	req_user_agent_hash := get_user_agent_sha512(ctx)
	user_agent_hash, err := deps.DB.Sessions().GetUserAgent(session_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to get ua_hash from db")
		return
	}
	if user_agent_hash != req_user_agent_hash {
		deps.Log.Error().
			Str("user_agent", user_agent_hash).
			Str("req_user_agent", req_user_agent_hash).
			Msg("Initial User-Agent not equal to request User-Agent")
		err = deps.DB.PurgeSessionFromTokenId(token_id)
		if err != nil {
			deps.Log.Error().Err(err).Msg("Failed to purge session")
		}
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Request is invalid"))
		return
	}

	init_ip, err := deps.DB.Sessions().GetInitIp(session_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to get init_ip from db")
		return
	}
	cur_ip := get_client_ip_sha512(ctx)
	if cur_ip != init_ip {
		deps.Log.Warn().
			Str("init_ip", init_ip).
			Str("cur_ip", cur_ip).
			Msg("Initial IP not equal for current IP")
		go call_webhook(deps.Log, get_client_ip(ctx))
	}

	access, refresh, new_token_id, err := token.GenPair(
		token.DefaultAccessTokenTimeout,
		token.DefaultRefreshTokenTimeout,
	)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to create pair of tokens")
		return
	}

	refresh_hash, err := get_refresh_token_bcrypt(refresh)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().
			Err(err).
			Msg("Failed to create bcrypt hash from refresh token")
		return
	}

	err = deps.DB.UpdateSession(
		session_id,
		token_id,
		new_token_id,
		refresh_hash,
		time.Now().Add(token.DefaultRefreshTokenTimeout).Add(24*time.Hour),
	)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to update session in db")
		return
	}

	ctx.JSON(iris.Map{
		"tokens": iris.Map{
			"access":  access,
			"refresh": refresh,
		},
	})
	deps.Log.Debug().
		Str("access", access).
		Str("refresh", refresh).
		Send()
}

func close_current_session(ctx iris.Context, deps *Deps) {
	deps.Log.Debug().Msg("close_current_session")
	if !ctx.Values().Exists("token_id") {
		deps.Log.Fatal().Msg("Authorization is not passed")
	}

	token_id := ctx.Values().Get("token_id").(string)

	err := deps.DB.PurgeSessionFromTokenId(token_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Str("token_id", token_id).Msg("Failed to purge session")
		return
	}
	ctx.StatusCode(iris.StatusNoContent)
}
