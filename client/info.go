package client

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
)

type InfoRes struct {
    TotalNumEnvs int `json:"total_environments"`
}

func Info(hostStr string, env string, apiKey string, apiSecret string) (InfoRes, error) {
    emptyRes := InfoRes{}
    //TODO: sign the request
    endPt := fmt.Sprintf("%s/info", stripEndSlash(hostStr))
    if len(env) > 0 {
        endPt = fmt.Sprintf("%s?env=%s", endPt, env)
    }
    res, err := http.Get(endPt)
    if err != nil {
        return emptyRes, err
    }
    bytes, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        return emptyRes, err
    }

    infoRes := InfoRes{}
    err = json.Unmarshal(bytes, &infoRes)
    if err != nil {
        return emptyRes, err
    }
    return infoRes, nil
}
