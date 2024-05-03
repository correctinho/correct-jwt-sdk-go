package jwtsdk

import stg "github.com/correctinho/correct-util-sdk-go/stg"

// PostEncodeRequest é a estrutura de solicitação para a codificação de dados.
type PostEncodeRequest struct {
	Data    any `json:"data"`              // Dados a serem codificados.
	Extras  any `json:"extras,omitempty"`  // Informações extras a serem codificadas.
	Seconds int `json:"seconds,omitempty"` // Número de segundos que o token será válido.
}

// Validate - valida os campos do payload
func (p *PostEncodeRequest) Validate() (bool, *JwtError) {
	if p.Data == nil {
		return false, &ErrDataRequired
	}
	return true, nil
}

// PostDecodeRequest é a estrutura de solicitação para a decodificação de um token.
type PostDecodeRequest struct {
	Token string `json:"token"` // Token a ser decodificado.
}

// Validate - valida os campos do payload
func (p *PostDecodeRequest) Validate() (bool, *JwtError) {
	if stg.IsEmpty(&p.Token) {
		return false, &ErrTokenRequired
	}
	return true, nil
}
