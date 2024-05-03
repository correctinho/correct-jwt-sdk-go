package jwtsdk

import (
	"errors"
	"net/http"
	"strings"
)

// JwtError - tipo de error
type JwtError struct {
	Message  error
	Code     string
	httpCode int
}

// Error implements error
func (e JwtError) Error() string {
	var r strings.Builder
	r.WriteString(e.Code + ": " + e.Message.Error() + ".")
	return r.String()
}

// Mensagens de erro
var (
	ErrServiceUnavailable = JwtError{
		Message:  errors.New("Serviço indisponível no momento. Por favor, tente novamente em alguns instantes"),
		Code:     "00000",
		httpCode: http.StatusServiceUnavailable,
	}

	ErrDataRequired = JwtError{
		Message:  errors.New("Os dados são obrigatórios"),
		Code:     "00001",
		httpCode: http.StatusBadRequest,
	}

	ErrTokenRequired = JwtError{
		Message:  errors.New("O token é obrigatório"),
		Code:     "00002",
		httpCode: http.StatusBadRequest,
	}
)
