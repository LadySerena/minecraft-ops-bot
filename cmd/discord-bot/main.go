package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bsdlp/discord-interactions-go/interactions"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	serverBindEnv     = "BIND_ADDRESS"
	serverBindDefault = ":8080"
	publicKeyEnv      = "PUBLIC_KEY"
	publicKeyString   = "dc2a2ef24d22c445bd5a81bab30219e7f1ebbaa8035513457cac4b145b32cdc3"
)

var (
	publicKey   = []byte{}
	bindAddress = serverBindDefault
)

func ping(w http.ResponseWriter, r *http.Request) {
	zap.S().Debug("got request")
	verified := interactions.Verify(r, publicKey)
	if !verified {
		zap.S().Error("could not verify request")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()
	var data interactions.Data

	decodeErr := json.NewDecoder(r.Body).Decode(&data)

	if decodeErr != nil {
		zap.S().Errorf("could not decode the request: %s", decodeErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.Type == interactions.Ping {
		zap.S().Debug("responding to ping")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"type":1}`)) // return the body since we are supposed to reply with {"type": 1}
		return
	}

	zap.S().Debug("processing real request")
	w.WriteHeader(http.StatusOK)
	response := &interactions.InteractionResponse{
		Type: interactions.ChannelMessage,
		Data: &interactions.InteractionApplicationCommandCallbackData{
			Content: "meow",
		},
	}

	var responseBuffer bytes.Buffer
	encodeErr := json.NewEncoder(&responseBuffer).Encode(response)
	if encodeErr != nil {
		zap.S().Errorf("could not encode response: %s", encodeErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, postErr := http.Post(data.ResponseURL(), "application/json", &responseBuffer)
	if postErr != nil {
		zap.S().Errorf("could not post data to url: %s", postErr.Error())
	}
	
}

func main() {
	viper.SetDefault(serverBindEnv, serverBindDefault)
	viper.AutomaticEnv()
	bindAddress = viper.GetString(serverBindEnv)
	//publicKey = []byte(viper.GetString(publicKeyEnv)) // for the verify interaction function the key must be a byte slice
	logger, loggerInitErr := zap.NewDevelopment()
	if loggerInitErr != nil {
		log.Fatalf("could not initialize zap logger: %v", loggerInitErr)
	}
	zap.ReplaceGlobals(logger)

	zap.S().Debugf("public key for verifying requests: %v", publicKeyString)

	publicKeyTemp, decodeErr := hex.DecodeString(publicKeyString)
	if decodeErr != nil {
		zap.S().Errorf("could not decode public key hex string: %v", decodeErr)
		return
	}
	publicKey = publicKeyTemp
	router := mux.NewRouter()
	router.HandleFunc("/", ping).Methods(http.MethodPost)

	zap.S().Error(http.ListenAndServe(bindAddress, router))
}
