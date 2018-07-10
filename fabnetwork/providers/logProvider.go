package providers

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/logging/api"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"log"
	"fmt"
)

type FabLogProvider struct {
	logger *utils.AppLogger
	networkName string
}

type fabLogger struct {
	logger *log.Logger
}

func GetAppLogProvider(nwName, module string) api.LoggerProvider {
	provider:= &FabLogProvider{
		logger:      utils.NewAppLogger(module, fmt.Sprintf("[%s] [%s] ", nwName, module)),
		networkName: nwName,
	}
	return provider
}

func (l *FabLogProvider) GetLogger(module string) api.Logger {
	//appLogger:= log.New(l.logger.Writer(), fmt.Sprintf("[%s] [%s] ", l.networkName, module), log.Ldate|log.Ltime|log.LUTC)
	appLogger:= utils.NewAppLogger("", fmt.Sprintf("[%s] #%s ", module, l.networkName))
	return appLogger
}