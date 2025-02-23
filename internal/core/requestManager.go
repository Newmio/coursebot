package core

import (
	"cbot/pkg"
	"io"
	"net/http"
)

type RequestManagerImpl struct {
	mClient *http.Client
}

func CreateRequestManager() pkg.RequestManager {
	obj := &RequestManagerImpl{}
	obj.mClient = &http.Client{}
	return obj
}

func (obj *RequestManagerImpl) Do(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, pkg.Trace(err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := obj.mClient.Do(req)
	if err != nil {
		return nil, pkg.Trace(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pkg.Trace(err)
	}

	return body, nil
}
