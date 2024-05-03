package jwtsdk

import (
	"github.com/correctinho/correct-mlt-go/qlog"
	types "github.com/correctinho/correct-types-sdk-go"
)

// Decode - Chama o microserviço de JWT para realizar o decode
func (c *Client) Decode(token string) (types.JwtToken, *JwtError) {
	logger := qlog.NewProduction(c.Context)
	defer logger.Sync()
	var response types.JwtToken

	client := newHTTP(c.Context)
	request := PostDecodeRequest{}
	request.Token = token

	if ok, err := request.Validate(); !ok || err != nil {
		return response, err
	}

	if err := client.Do("POST", "/api/v1/jwt/decode", request, &response, true); err != nil {
		logger.Error(err.Error())
		return response, &ErrServiceUnavailable
	}
	return response, nil
}
