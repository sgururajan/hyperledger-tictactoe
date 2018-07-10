package main

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/log"
	"github.com/sgururajan/hyperledger-tictactoe/appServer"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/database"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/networkHandlers"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"github.com/gorilla/mux"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/apiHandlers"
	"net/http"
	"time"
	"os/signal"
	"context"
	"flag"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"github.com/gorilla/handlers"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork"
)

var networksDbFile = "networks.json"
var serverAddr = "0.0.0.0:4300"

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "shutdown-timeout", 15 * time.Second, "duration for which server will gracefully wait for request to complete before shutting down")
	flag.Parse()

	logger:= utils.NewAppLogger("main", "")

	viper.SetConfigFile("appsettings.json")
	verr := viper.ReadInConfig()

	if verr != nil {
		logger.Fatalf("error while reading config")
		panic(verr)
	}

	logger.Debugf("appsetting - configTxGenToolsPath: %s", viper.GetString("configTxGenToolsPath"))

	err := setupDefaultNetwork()
	if err != nil {
		panic(err)
	}
	repo := setupRepository()
	networkHandler, err := networkHandlers.NewNetworkHandler(repo)
	if err != nil {
		panic(err)
	}

	defer networkHandler.Close()

	t3Network, err:= networkHandler.GetNetwork("testnetwork")
	if err != nil {
		panic(err)
	}

	ensureTictactoeChannelAndChainCode(t3Network)

	//enable cors
	allowedHandlers:= handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowedOrigins:= handlers.AllowedOrigins([]string{"*"})
	allowedMethods:= handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "HEAD"})

	//testNetwork(networkHandler)
	router:= mux.NewRouter()
	apiHandler:= apiHandlers.NewNetworkAPIHandler(repo, networkHandler)
	apiHandler.RegisterRoutes(router.PathPrefix("/api").Subrouter())

	router.Use(loggingMiddleWare)

	server:= &http.Server{
		Addr: serverAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
		Handler: handlers.CORS(allowedHandlers,allowedMethods, allowedOrigins)(router),
	}

	go func(){
		logger.Infof("server starting and listening on %v", serverAddr)
		if err= server.ListenAndServe(); err!= nil {
			log.Fatal(err)
		}
	}()

	c:= make(chan os.Signal,1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx,cancel:= context.WithTimeout(context.Background(), wait)
	defer cancel()

	server.Shutdown(ctx)
	log.Info("tictactoe server signing off")
	os.Exit(0)
}

func setupRepository() database.NetworkRepository {
	return database.NewNetworkFileRepository(networksDbFile)
}

func loggingMiddleWare(next http.Handler) http.Handler {
	logger:= utils.NewAppLogger("http", "")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start:=time.Now()
		next.ServeHTTP(w,r)
		logger.Debugf("%s\t%s\t%s", r.Method, r.RequestURI, time.Since(start))
	})
}

func setupDefaultNetwork() error {
	if _, err := os.Stat(networksDbFile); err != nil {
		networks := getDefaultNetwork()

		cbytes, err := json.MarshalIndent(networks, "", "\t")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(networksDbFile, cbytes, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func getDefaultNetwork() map[string]database.Network {
	dNetwork := appServer.DefaultNetworkConfiguration()
	networks := make(map[string]database.Network)
	networks[dNetwork.Name] = dNetwork

	return networks
}

func ensureTictactoeChannelAndChainCode(fabNetwork *fabnetwork.FabricNetwork) {
	chReq:= entities.CreateChannelRequest{
		ChannelName:"tictactoechannel",
		OrganizationNames:[]string{
			"org1",
			"org2",
		},
		AnchorPeers: map[string][]string{
			"org1": {
				"peer0.org1.tictactoe.com",
				"peer1.org1.tictactoe.com",
			},
			"org2": {
				"peer0.org2.tictactoe.com",
				"peer1.org2.tictactoe.com",
			},
		},
		ConsortiumName:"tictactoechannelconsortium",
	}

	err:= fabNetwork.CreateChannel("org1", chReq)
	if err != nil {
		panic(err)
	}

	ccRequest := entities.InstallChainCodeRequest{
		ChainCodeName:    "tictactoe",
		ChainCodePath:    "github.com/sgururajan/hyperledger-tictactoe/chaincodes/tictactoe/",
		ChainCodeVersion: "0.0.4",
		ChannelName:      "tictactoechannel",
	}

	err = fabNetwork.InstallChainCode([]string{"org1", "org2"}, ccRequest)

	if err != nil {
		panic(err)
	}
}

func testNetwork(networkHandler *networkHandlers.NetworkHandler) {
	network, err := networkHandler.GetNetwork("testnetwork")
	if err != nil {
		log.Errorf("error while getting network. err: %v", err)
		os.Exit(0)
	}

	chReq := entities.CreateChannelRequest{
		ChannelName: "testchannel",
		OrganizationNames: []string{
			"org1",
			"org2",
		},
		AnchorPeers: map[string][]string{
			"org1": {
				"peer0.org1.tictactoe.com",
				"peer1.org1.tictactoe.com",
			},
			"org2": {
				"peer0.org2.tictactoe.com",
			},
		},
		ConsortiumName: "testconsortium",
	}

	err = network.CreateChannel("org1", chReq)

	if err != nil {
		panic(err)
	}

	ccRequest := entities.InstallChainCodeRequest{
		ChainCodeName:    "sample",
		ChainCodePath:    "github.com/sgururajan/hyperledger-tictactoe/chaincodes/sample/",
		ChainCodeVersion: "0.0.4",
		ChannelName:      "testchannel",
	}

	err = network.InstallChainCode([]string{"org1", "org2"}, ccRequest)

	if err != nil {
		panic(err)
	}
}
