package v0

import (
	"crypto/sha512"
	"encoding/base64"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/qcrg/bookish-guacamole/api/v0/token"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type auth_req struct {
	UserID string `json:"id"`
}

func get_ip_from_addr(addr string) string {
	ip, _, err := net.SplitHostPort(addr)
	if err != nil {
		ip = addr
	}
	return ip
}

func get_client_ip(ctx iris.Context) string {
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP",
		"True-Client-IP",
	}
	for _, header := range headers {
		ip := ctx.GetHeader(header)
		if ip != "" {
			return get_ip_from_addr(strings.Split(ip, ",")[0])
		}
	}
	return get_ip_from_addr(ctx.RemoteAddr())
}

func get_client_ip_sha512(ctx iris.Context) string {
	ip := get_client_ip(ctx)
	buf := sha512.Sum512([]byte(ip))
	return base64.StdEncoding.EncodeToString(buf[:])
}

func get_user_agent_sha512(ctx iris.Context) string {
	user_agent := ctx.GetHeader("User-Agent")
	user_agent_sha512 := sha512.Sum512([]byte(user_agent))
	hashed_user_agent := base64.StdEncoding.EncodeToString(user_agent_sha512[:])
	return hashed_user_agent
}

func get_refresh_token_bcrypt(refresh_token string) (string, error) {
	hashed_sha512_refresh := sha512.Sum512([]byte(refresh_token))
	bcrypted_refresh, err := bcrypt.GenerateFromPassword(
		hashed_sha512_refresh[:],
		bcrypt.DefaultCost,
	)
	return string(bcrypted_refresh), err
}

func create_session(ctx iris.Context, deps *Deps) {
	req := auth_req{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Body of request must be JSON"))
		return
	}
	err = uuid.Validate(req.UserID)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(NewErrResp(ErrGeneral, "Invalid ID format"))
		return
	}
	users := deps.DB.Users()
	exists, err := users.Exists(req.UserID)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("postgres error")
		return
	}
	if !exists {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(NewErrResp(ErrGeneral, "User with ID not found"))
		return
	}

	access, refresh, token_id, err := token.GenPair(
		token.DefaultAccessTokenTimeout,
		token.DefaultRefreshTokenTimeout,
	)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to create pair of tokens")
		return
	}

	bcrypted_refresh, err := get_refresh_token_bcrypt(refresh)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().
			Err(err).
			Msg("Failed to create bcrypt hash from refresh token")
		return
	}
	log.Debug().
		Int("size", len(bcrypted_refresh)).
		Str("bcrypted", bcrypted_refresh).
		Send()

	session_id, err := deps.DB.CreateSession(
		req.UserID,
		token_id,
		bcrypted_refresh,
		get_user_agent_sha512(ctx),
		get_client_ip_sha512(ctx),
		time.Now().Add(token.DefaultRefreshTokenTimeout).Add(24*time.Hour),
	)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().
			Err(err).
			Msg("Failed to create session")
		return
	}
	deps.Log.Debug().Str("session_id", session_id).Msg("New session")

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
