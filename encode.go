package jwtsdk

import "github.com/correctinho/correct-mlt-go/qlog"

// Encode - Chama o microserviço de JWT para realizar o encode
func (c *Client) Encode(data, extras any, seconds int) (PostEncodeResponse, *JwtError) {
	logger := qlog.NewProduction(c.Context)
	defer logger.Sync()
	var response PostEncodeResponse

	client := newHTTP(c.Context)
	request := PostEncodeRequest{}
	request.Data = data
	request.Extras = extras
	request.Seconds = seconds

	if ok, err := request.Validate(); !ok || err != nil {
		return response, err
	}

	if err := client.Do("POST", "/api/v1/jwt/encode", request, &response, true); err != nil {
		logger.Error(err.Error())
		return response, &ErrServiceUnavailable
	}
	return response, nil
}
