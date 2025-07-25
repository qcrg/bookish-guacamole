package v0

import "github.com/kataras/iris/v12"

func me(ctx iris.Context, deps *Deps) {
	if !ctx.Values().Exists("token_id") {
		deps.Log.Fatal().Msg("Authorization is not passed")
	}

	token_id := ctx.Values().Get("token_id").(string)

	ss := deps.DB.Sessions()

	exists, err := ss.ExistsFromTokenId(token_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to check existence of session")
		return
	}
	if !exists {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Str("token_id", token_id).Msg("Session is not exists")
		return
	}
	user_id, err := ss.GetUserIdFromTokenId(token_id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		deps.Log.Error().Err(err).Msg("Failed to get user_id from session")
		return
	}
	ctx.JSON(iris.Map{"user_id": user_id})
}
