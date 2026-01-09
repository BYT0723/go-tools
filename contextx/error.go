package ctxx

var (
	ErrApiKeyNotFound    = ctxNotFoundErr{"api key"}
	ErrWaitGroupNotFound = ctxNotFoundErr{"waitgroup"}
	ErrListenerNotFound  = ctxNotFoundErr{"listener"}
	ErrLoggerNotFound    = ctxNotFoundErr{"logger"}
)

type ctxNotFoundErr struct {
	key string
}

func (e ctxNotFoundErr) Error() string {
	return e.key + " not found"
}
