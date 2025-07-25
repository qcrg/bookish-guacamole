package v0

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/qcrg/bookish-guacamole/api/v0/token"
)

func auth(ctx iris.Context, deps *Deps) {
	deps.Log.Debug().Msg("auth")

	fields := strings.Fields(ctx.GetHeader("Authorization"))

	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(NewErrResp(ErrGeneral, "Token is invalid"))
		return
	}

	token_id, typ, err := token.Parse(fields[1])
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(NewErrResp(ErrGeneral, err.Error()))
		return
	}
	exists, err := deps.DB.Tokens().Exists(token_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to check existence of token")
		return
	}
	if !exists {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(NewErrResp(ErrGeneral, "Token is invalid"))
		deps.Log.Debug().Str("token_id", token_id).Msg("Token not found")
		return
	}
	if typ != token.TypeAccess {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(NewErrResp(ErrGeneral, "Type token must be access"))
		return
	}
	ctx.Values().SetImmutable("token_id", token_id)
	ctx.Next()
}
