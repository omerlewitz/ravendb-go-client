package ravendb

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
)

var (
	_ RavenCommand = &PutSecretKeyCommand{}
)

type PutSecretKeyCommand struct {
	RavenCommandBase

	_name      string
	_base64key string
	_overwrite bool

	Result *PutResult
}

func NewPutSecretKeyCommand(name string, base64key string, overwrite bool) *PutSecretKeyCommand {
	cmd := &PutSecretKeyCommand{
		RavenCommandBase: NewRavenCommandBase(),
		_name:            name,
		_base64key:       base64key,
		_overwrite:       overwrite,
	}
	return cmd
}

func (c *PutSecretKeyCommand) CreateRequest(node *ServerNode) (*http.Request, error) {
	base, err := url.Parse(node.URL + "/admin/secrets")
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("name", c._name)
	params.Add("overwrite", strconv.FormatBool(c._overwrite))
	base.RawQuery = params.Encode()

	keyBytes := []byte(c._base64key)
	return http.NewRequest(http.MethodPost, base.String(), bytes.NewBuffer(keyBytes))
}

func (c *PutSecretKeyCommand) SetResponse(response []byte, fromCache bool) error {
	return jsonUnmarshal(response, &c.Result)
}
