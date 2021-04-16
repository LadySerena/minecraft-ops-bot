package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	serverBindEnv     = "BIND_ADDRESS"
	serverBindDefault = ":8080"
	publicKeyEnv      = "PUBLIC_KEY"
	publicKeyString   = "dc2a2ef24d22c445bd5a81bab30219e7f1ebbaa8035513457cac4b145b32cdc3"
)

var (
	publicKey = []byte{}
	bindAddress = serverBindDefault
)

type discordPing struct {
	MessageType int `json:"type"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	zap.S().Info("got request")
	if !discordgo.VerifyInteraction(r, publicKey) {
		zap.S().Info("could not verify request")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		zap.S().Errorf("could not read the request body: %v", readErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	parsedBody := discordPing{}

	unmarshalErr := json.Unmarshal(body, &parsedBody)

	if unmarshalErr != nil {
		zap.S().Errorf("could not unmarshal the response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if parsedBody.MessageType == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write(body) // return the body since we are supposed to reply with {"type": 1}
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}



func main() {
	viper.SetDefault(serverBindEnv, serverBindDefault)
	viper.AutomaticEnv()
	bindAddress = viper.GetString(serverBindEnv)
	//publicKey = []byte(viper.GetString(publicKeyEnv)) // for the verify interaction function the key must be a byte slice
	logger, loggerInitErr := zap.NewProduction()
	if loggerInitErr != nil {
		log.Fatalf("could not initialize zap logger: %v", loggerInitErr)
	}
	zap.ReplaceGlobals(logger)

	zap.S().Infof("public key for verifying requests: %v", publicKeyString)

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