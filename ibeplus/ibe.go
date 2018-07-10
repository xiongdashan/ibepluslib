package ibeplus

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/otwdev/galaxylib"
)

type IBEPlus struct {
	RequestURL      string
	RequestParms    string
	RequestUserID   string
	RequestUsername string
}

func NewIBE(url string, data interface{}) *IBEPlus {
	dataByte, _ := xml.Marshal(data)
	return &IBEPlus{
		RequestURL:   url,
		RequestParms: string(dataByte),
	}
}

func (i *IBEPlus) Reqeust() (buffer []byte, retErr *galaxylib.GalaxyError) {

	URL := galaxylib.GalaxyCfgFile.MustValue("ibe", "transurl")
	client := http.DefaultClient

	buffer = nil

	// fmt.Println(i.RequestParms)
	paramsJSON, err := json.Marshal(i)
	if err != nil {
		retErr = galaxylib.DefaultGalaxyError.FromError(1, err)
		return
	}

	request, err := http.NewRequest("POST", URL, bytes.NewReader(paramsJSON))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		retErr = galaxylib.DefaultGalaxyError.FromError(1, err)
		return
		//return nil, utils.FromError(1, err)
	}

	response, err := client.Do(request)
	if err != nil {
		retErr = galaxylib.DefaultGalaxyError.FromError(1, err)
		return
		//return nil, utils.FromError(1, err)
	}

	defer response.Body.Close()
	var rev []byte
	if rev, err = ioutil.ReadAll(response.Body); err != nil {
		retErr = galaxylib.DefaultGalaxyError.FromError(1, err)
		return
	}
	rev = bytes.Replace(rev, []byte("\\\""), []byte("'"), -1)
	return rev, nil

	//return response.
}
