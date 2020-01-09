package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"

	terraapp "github.com/terra-project/core/app"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	core "github.com/terra-project/core/types"
)

// EncoderResponse is struct for sending encoded tx deta back to the caller
type EncoderResponse struct {
	EncodedTx []byte `json:"encoded_tx"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	cdc := terraapp.MakeCodec()

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var stdTx authtypes.StdTx

	if err = cdc.UnmarshalJSON(b, &stdTx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(stdTx)

	encResp := EncoderResponse{txBytes}

	js, err := json.Marshal(encResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func main() {

	config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(core.Bech32PrefixConsAddr, core.Bech32PrefixConsPub)
	config.Seal()

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
