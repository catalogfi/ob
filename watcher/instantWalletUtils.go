package watcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func isBtcIwTx(url, txHash string) (bool, error) {
	resp, err := http.Get(url + "/validateTransaction/" + txHash)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return false, fmt.Errorf("failed to reach the server %v", 404)
		}
		errObj := struct {
			Error string `json:"error"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&errObj); err != nil {
			errMsg, err := io.ReadAll(resp.Body)
			if err != nil {
				return false, fmt.Errorf("failed to read the error message %v", err)
			}
			return false, fmt.Errorf("failed to decode the error %v", string(errMsg))
		}
		return false, fmt.Errorf("request failed %v", errObj.Error)
	}
	response := struct {
		Message bool `json:"message"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response.Message, err
}

func isBtcIwTxs(url string, initateTxHash string) (bool, error) {
	txHashes := strings.Split(initateTxHash, ",")
	n := len(txHashes)
	for _, txHash := range txHashes {
		isIw, err := isBtcIwTx(url, txHash)
		if err != nil {
			return false, err
		}
		if isIw {
			n--
		}

	}
	return n == 0, nil
}
