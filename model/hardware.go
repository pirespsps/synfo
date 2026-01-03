package model

import (
	"encoding/json"
	"fmt"

	"github.com/pirespsps/synfo/utils"
)

type Hardware interface {
	Overall() (Response, error)
	Extensive() (Response, error)
	getData() (any, error) //depois desse método, repensar response...
}

// transformar tudo em []byte ao invés de response, e dar unmarshal usando o json
type Response struct { //com o GetData fica inútil
	Data []any
}

func (r Response) Print() {
	fmt.Print("\n")
	fmt.Print(utils.PrintStruct(r.Data))
	fmt.Print("\n")
}

func (r Response) Json() ([]byte, error) {
	js, err := json.Marshal(r.Data)
	if err != nil {
		return nil, fmt.Errorf("error in marshal:%v", err)
	}

	return js, nil
}
