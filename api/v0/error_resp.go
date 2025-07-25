package v0

const (
	ErrGeneral int = -1
)

type ErrorData struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

type ErrorResp struct {
	Response
	Error ErrorData `json:"error"`
}

func NewErrResp(code int, reason string) ErrorResp {
	return ErrorResp{
		Error: ErrorData{
			Code:   code,
			Reason: reason,
		},
	}
}
