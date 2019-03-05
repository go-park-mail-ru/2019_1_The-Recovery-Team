package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/mailru/easyjson"
)

func unmarshalJSONBodyToStruct(r *http.Request, strct easyjson.Unmarshaler) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = easyjson.Unmarshal(body, strct)
	if err != nil {
		return err
	}

	return nil
}
