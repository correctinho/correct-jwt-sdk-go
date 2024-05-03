package jwtsdk

import (
	"sync"
)

// Client - classe que representa o objeto do SDK
type Client struct {
	sync.RWMutex
	Context interface{}
}
