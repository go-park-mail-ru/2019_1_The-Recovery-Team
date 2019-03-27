package unmarshaler

import (
	"io/ioutil"
	"net/http"

	"github.com/mailru/easyjson"
)

func UnmarshalJSONBodyToStruct(r *http.Request, s easyjson.Unmarshaler) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	err = easyjson.Unmarshal(body, s)
	if err != nil {
		return err
	}

	return nil
}
