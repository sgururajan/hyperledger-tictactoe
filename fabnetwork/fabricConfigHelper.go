package fabnetwork

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork/entities"
	utils2 "github.com/sgururajan/hyperledger-tictactoe/utils"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"fmt"
	"sync"
)

type fabCapability struct {
	V1 bool `yaml:"V1_1"`
}

type fabCapabilities struct {
	Global      fabCapability
	Orderer     fabCapability
	Application fabCapability
}

type fabBatchSize struct {
	MaxMessageCount   uint32 `yaml:"MaxMessageCount"`
	AbsoluteMaxBytes  uint32 `yaml:"AbsoluteMaxBytes"`
	PreferredMaxBytes uint32 `yaml:"PreferredMaxBytes"`
}

type fabKafka struct {
	Brokers []string `yaml:"Brokers"`
}

type fabOrderer struct {
	OrdererType  string        `yaml:"OrdererType"`
	Addresses    []string      `yaml:"Addresses"`
	BatchTimeout time.Duration `yaml:"BatchTimeout"`
	BatchSize    fabBatchSize  `yaml:"BatchSize"`
	Kafka        fabKafka      `yaml:"Kafka"`
}

type fabOrganization struct {
	Name        string    `yaml:"Name"`
	ID          string    `yaml:"ID"`
	MSPDir      string    `yaml:"MSPDir"`
	AnchorPeers []fabPeer `yaml:"AnchorPeers,omitempty"`
	isOrderer   bool
	endpoint    string
}

type fabPeer struct {
	Host string `yaml:"Host"`
	Port uint32 `yaml:"Port"`
}
type fabApplication struct {
	Organizations []fabOrganization `yaml:"Organizations"`
}

type fabOrdererGenesis struct {
	Orderer       fabOrderer        `yaml:"Orderer,inline"`
	Organizations []fabOrganization `yaml:"Organizations"`
	Capabilities  fabCapability     `yaml:"Capabilities"`
}

type fabConsortium struct {
	name          string
	Organizations []fabOrganization `yaml:"Organizations"`
}

type fabGenesisProfile struct {
	Capabilities   fabCapability            `yaml:"Capabilities"`
	OrdererGenesis fabOrdererGenesis        `yaml:"Orderer"`
	Consortiums    map[string]fabConsortium `yaml:"Consortiums"`
}

type fabChannelApplication struct {
	Organizations []fabOrganization `yaml:"Organizations"`
	Capabilities  fabCapability     `yaml:"Capabilities"`
	application   fabApplication
}

type fabChannelProfile struct {
	name               string
	Consortium         string                `yaml:"Consortium"`
	ChannelApplication fabChannelApplication `yaml:"application"`
}

type fabProfile struct {
	GenesisProfiles map[string]fabGenesisProfile
	Channel         map[string]fabChannelProfile
}


type fabChannelOrg struct {
	Name        string
	AnchorPeers []string
	IsOrderer   bool
	Endpoint    string
	Consortium  string
}


type fabConfigTxHelper struct {
	Profiles map[string]interface{} `yaml:"Profiles"`
	rwLock   sync.RWMutex
}

func newFabConfixTxHelper() *fabConfigTxHelper {
	return &fabConfigTxHelper{}
}

func (m *fabConfigTxHelper) createChannelTransaction(network *FabricNetwork, chRequest entities.CreateChannelRequest, txFileName string) error {

	capabilities := getCapabilities()
	chOrganizations := getChannelOrganizations(network, chRequest)
	//consortium := fabConsortium{
	//	Organizations: chOrganizations,
	//	name:          chRequest.ConsortiumName,
	//}

	chApplication := fabChannelApplication{
		Organizations: chOrganizations,
		Capabilities:  capabilities.Application,
	}

	chProfile := fabChannelProfile{
		Consortium:         chRequest.ConsortiumName,
		ChannelApplication: chApplication,
		name:               chRequest.ChannelName,
	}

	m.Profiles = make(map[string]interface{})

	m.Profiles[chRequest.ChannelName] = chProfile
	err:= m.createChannelTransactionBlock(txFileName)
	if err != nil {
		return err
	}

	return nil
}

func (m *fabConfigTxHelper) createChannelTransactionBlock(txFileName string) error {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	ybytes, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	txConfigFilePath := filepath.Join(utils2.Substitute(viper.GetString("configTxGenToolsPath")), "configtx.yaml")
	err = ioutil.WriteFile(txConfigFilePath, ybytes, os.ModePerm)
	if err != nil {
		return err
	}

	toolsDir := utils2.Substitute(viper.GetString("configTxGenToolsPath"))
	err = os.Chdir(toolsDir)
	if err != nil {
		return err
	}

	for k, v := range m.Profiles {
		if _, ok := v.(fabChannelProfile); ok {
			args := []string{"-profile", k, "-outputCreateChannelTx", txFileName, "-channelID", k}
			log.Infof("cmd: ./configtxgen %#v", args)
			cmd := exec.Command("./configtxgen", args...)
			output, err := cmd.Output()
			fmt.Println(string(output))
			if err != nil {
				return err
			}

			log.Info(output)
		}
	}

	return nil
}

func getChannelOrganizations(network *FabricNetwork, chRequest entities.CreateChannelRequest) []fabOrganization {
	result := []fabOrganization{}
	for _, orgName := range chRequest.OrganizationNames {
		org, ok := network.orgsByName[orgName]
		if !ok {
			continue
		}

		orgPeers, err := getOrgPeersAsMapByEndpoint(network, orgName)
		if err != nil {
			continue
		}

		fabOrgPeers := []fabPeer{}
		for _, pv := range chRequest.AnchorPeers[orgName] {
			ep, ok := orgPeers[pv]
			if !ok {
				continue
			}

			fabOrgPeers = append(fabOrgPeers, fabPeer{
				Host: ep.EndPoint,
				Port: utils2.GetPortFromUrl(ep.URL),
			})
		}

		result = append(result, fabOrganization{
			MSPDir:      utils2.Substitute(org.MSPDir),
			Name:        org.MSPID,
			ID:          org.MSPID,
			AnchorPeers: fabOrgPeers,
		})
	}

	return result
}

func getOrgPeersAsMapByEndpoint(network *FabricNetwork, orgName string) (map[string]entities.Peer, error) {
	orgPeers, err := network.providers.peerProvider.GetPeersForOrgId(orgName)
	if err != nil {
		return nil, err
	}

	result := make(map[string]entities.Peer)
	for _, v := range orgPeers {
		result[v.EndPoint] = v
	}

	return result, nil
}

func getCapabilities() fabCapabilities {
	return fabCapabilities{
		Application: fabCapability{V1: true},
		Global:      fabCapability{V1: true},
		Orderer:     fabCapability{V1: true},
	}
}
