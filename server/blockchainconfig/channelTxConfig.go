package blockchainconfig

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/sgururajan/hyperledger-tictactoe/server/common"
	"github.com/sgururajan/hyperledger-tictactoe/server/networkconfig"
	"github.com/sgururajan/hyperledger-tictactoe/server/pathUtil"
	"github.com/sgururajan/hyperledger-tictactoe/server/settings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Capability struct {
	V1 bool `yaml:"V1_1"`
}

type Capabilities struct {
	Global      Capability
	Orderer     Capability
	Application Capability
}

type BatchSize struct {
	MaxMessageCount   uint32 `yaml:"MaxMessageCount"`
	AbsoluteMaxBytes  uint32 `yaml:"AbsoluteMaxBytes"`
	PreferredMaxBytes uint32 `yaml:"PreferredMaxBytes"`
}

type Kafka struct {
	Brokers []string `yaml:"Brokers"`
}

type Orderer struct {
	OrdererType  string        `yaml:"OrdererType"`
	Addresses    []string      `yaml:"Addresses"`
	BatchTimeout time.Duration `yaml:"BatchTimeout"`
	BatchSize    BatchSize     `yaml:"BatchSize"`
	Kafka        Kafka         `yaml:"Kafka"`
}

type Organization struct {
	Name        string `yaml:"Name"`
	ID          string `yaml:"ID"`
	MSPDir      string `yaml:"MSPDir"`
	AnchorPeers []Peer `yaml:"AnchorPeers,omitempty"`
	isOrderer   bool
	endpoint    string
}

type Peer struct {
	Host string `yaml:"Host"`
	Port uint32 `yaml:"Port"`
}

type Application struct {
	Organizations []Organization `yaml:"Organizations"`
}

type OrdererGenesis struct {
	Orderer       Orderer        `yaml:"Orderer,inline"`
	Organizations []Organization `yaml:"Organizations"`
	Capabilities  Capability     `yaml:"Capabilities"`
}

type Consortium struct {
	name          string
	Organizations []Organization `yaml:"Organizations"`
}

type GenesisProfile struct {
	Capabilities   Capability            `yaml:"Capabilities"`
	OrdererGenesis OrdererGenesis        `yaml:"Orderer"`
	Consortiums    map[string]Consortium `yaml:"Consortiums"`
}

type ChannelApplication struct {
	Organizations []Organization `yaml:"Organizations"`
	Capabilities  Capability     `yaml:"Capabilities"`
	application   Application
}

type ChannelProfile struct {
	name               string
	Consortium         string             `yaml:"Consortium"`
	ChannelApplication ChannelApplication `yaml:"application"`
}

type Profile struct {
	GenesisProfiles map[string]GenesisProfile
	Channel         map[string]ChannelProfile
}

type ChannelOrg struct {
	Name        string
	AnchorPeers []string
	IsOrderer   bool
	Endpoint    string
	Consortium  string
}

type ChannelTxConfigHelper struct {
	Profiles map[string]interface{} `yaml:"Profiles"`
	rwLock   sync.RWMutex
}

var logger = logging.NewLogger("channelTxConfig")

func NewChannelTxConfigHelper() *ChannelTxConfigHelper {
	return &ChannelTxConfigHelper{}
}

func (m *ChannelTxConfigHelper) CreateChannelTxObject(config *networkconfig.FabNetworkConfiguration, chRequest common.CreateChannelRequest) error {

	capabilities := Capabilities{
		Application: Capability{V1: true},
		Global:      Capability{V1: true},
		Orderer:     Capability{V1: true},
	}

	orgs:= m.getChannelOrganizations(config, chRequest)
	ord, _:= m.getChannelOrderer(config, config.OrderersInfo[0].Endpoint)

	ordGenesisProfile := OrdererGenesis{
		Capabilities:  capabilities.Orderer,
		Orderer:       ord,
		Organizations: orgs,
	}

	consortium := Consortium{
		Organizations: orgs,
		name:          chRequest.ConsortiumName,
	}

	genesisProfiles := GenesisProfile{
		Capabilities:   capabilities.Global,
		OrdererGenesis: ordGenesisProfile,
	}

	genesisProfiles.Consortiums = make(map[string]Consortium)
	genesisProfiles.Consortiums[chRequest.ConsortiumName] = consortium

	chApplication := ChannelApplication{
		Organizations: orgs,
		Capabilities:  capabilities.Application,
	}

	chProfile := ChannelProfile{
		Consortium:         chRequest.ConsortiumName,
		ChannelApplication: chApplication,
		name:               chRequest.ChannelName,
	}

	m.Profiles = make(map[string]interface{})

	if !common.HasConsortium(config.Consortiums, chRequest.ConsortiumName) {
		//genesisKey:= fmt.Sprintf("%sGenesis", chRequest.ChannelName)
		//m.Profiles[genesisKey] = genesisProfiles
		// TODO: need to add the consortium to the orderer by updating the genesis block.
	}

	m.Profiles[chRequest.ChannelName] = chProfile

	return nil
}

//func (m *ChannelTxConfigHelper) CreateChannelTxObject1(config *networkconfig.FabNetworkConfiguration, channelName string, chOrgs []ChannelOrg, ordName string) error {
//	orgs := m.getChannelOrganizations(config, chOrgs)
//	//appDefaults:= application{Organizations: []Organization{} }
//	capabilities := Capabilities{
//		Application: Capability{V1: true},
//		Global:      Capability{V1: true},
//		Orderer:     Capability{V1: true},
//	}
//	ord, _ := m.getChannelOrderer(config, ordName)
//
//	var ordOrg Organization
//	for _, o := range orgs {
//		if o.endpoint == ordName {
//			ordOrg = o
//		}
//		//cOrg, ok:= config.Orderers[ordName]
//		//if ok && o.ID == cOrg.GRPCOptions["ssl-target-Name-override"].(string) {
//		//	ordOrg = o
//		//}
//	}
//
//	noOrdererOrgs := []Organization{}
//
//	for _, o := range orgs {
//		if !o.isOrderer {
//			noOrdererOrgs = append(noOrdererOrgs, o)
//		}
//	}
//
//	ordGenesisProfile := OrdererGenesis{
//		Capabilities:  capabilities.Orderer,
//		Orderer:       ord,
//		Organizations: []Organization{ordOrg},
//	}
//
//	consortiumName := fmt.Sprintf("%sConsortium", channelName)
//	consortium := Consortium{
//		Organizations: noOrdererOrgs,
//		name:          consortiumName,
//	}
//
//	genesisProfiles := GenesisProfile{
//		Capabilities:   capabilities.Global,
//		OrdererGenesis: ordGenesisProfile,
//	}
//
//	genesisProfiles.Consortiums = make(map[string]Consortium)
//	genesisProfiles.Consortiums[consortiumName] = consortium
//
//	chApplication := ChannelApplication{
//		Organizations: noOrdererOrgs,
//		Capabilities:  capabilities.Application,
//	}
//
//	chProfile := ChannelProfile{
//		Consortium:         consortiumName,
//		ChannelApplication: chApplication,
//		name:               channelName,
//	}
//
//	m.Profiles = make(map[string]interface{})
//
//	/*genesisKey:= fmt.Sprintf("%sGenesis", channelName)
//	m.Profiles[genesisKey] = genesisProfiles*/
//
//	m.Profiles[channelName] = chProfile
//
//	//ybytes, err:= yaml.Marshal(m)
//	//ioutil.WriteFile("yamlTest.yaml", ybytes, os.ModePerm)
//
//	return nil
//}

func (m *ChannelTxConfigHelper) CreateConfigurationBlocks(txFileName string) error {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	ybytes, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(settings.ConfigTxGenToolsPath, "configtx.yaml")
	err = ioutil.WriteFile(pathUtil.Substitute(configFilePath), ybytes, os.ModePerm)
	if err != nil {
		panic(err)
	}
	//toolPath := "./configtxgen"
	//toolPath:= pathUtil.Substitute(filepath.Join(settings.ConfigTxGenToolsPath, "configtxgen"))
	toolPath := pathUtil.Substitute(settings.ConfigTxGenToolsPath)
	//artifactsPath := pathUtil.Substitute(settings.ChannelArtifactsPath)

	//genesisBlockName := txFilePath + "/%s.block"
	//channelTxBlockName := txFilePath + "/%s.tx"

	//genesisBlockFormat:= "%s -profile %s -outputBlock %s/%s.block"
	//channelTxBlockFormat:= "%s -profile %s -outputCreateChannelTx %s/%s.tx -channelID %s"
	//anchorTxBlockFormat:= "%s -profile %s -outputAnchorPeersUpdate %s/%sAnchors.tx -channelID %s"

	err = os.Chdir(toolPath)
	if err != nil {
		return err
	}

	for key, val := range m.Profiles {
		/*if _, ok := val.(GenesisProfile); ok {
			args := []string{"-profile", key, "-outputBlock", gbFileName}
			cmd := exec.Command("./configtxgen", args...)
			//args:= []string{"-c", fmt.Sprintf("cd %s && ./configtxgen -profile %s -outputBlock %s", toolPath, key, fmt.Sprintf(genesisBlockName, key))}
			//cmd:= exec.Command("/bin/sh", args...)
			//cmd.Stdout=os.Stdout
			//err:= cmd.Run()
			output, err := cmd.CombinedOutput()
			fmt.Println(string(output))
			if err != nil {
				return err
			}
			fmt.Printf("Created genesis block for profile %s\n", key)
		} else*/
		if _, ok := val.(ChannelProfile); ok {

			args := []string{"-profile", key, "-outputCreateChannelTx", txFileName, "-channelID", key}
			cmd := exec.Command("./configtxgen", args...)

			//args:= []string{"-c", fmt.Sprintf("cd %s && ./configtxgen -profile %s -outputCreateChannelTx %s -channelID %s", toolPath, key, fmt.Sprintf(channelTxBlockName, key), key)}
			//cmd:= exec.Command("/bin/sh", args...)
			//cmd:= exec.Command(toolPath, "-profile", key, "-outputCreateChannelTx", fmt.Sprintf(channelTxBlockName, key), "-channelID", key)
			output, err := cmd.Output()
			fmt.Println(string(output))
			if err != nil {
				return err
			}
			logger.Infof("Created channel tx block for profile %s\n", key)
		}
	}

	return nil
}

func (m *ChannelTxConfigHelper) getChannelOrganizations(config *networkconfig.FabNetworkConfiguration, chRequest common.CreateChannelRequest) []Organization {
	orgs := []Organization{}

	for _, o := range chRequest.OrganizationNames {
		orgConfig, ok := config.Organizations[strings.ToLower(o)]
		if ok {
			orgPeers := []Peer{}
			_, pok:= chRequest.AnchorPeers[o]
			if pok {
				for _, p := range chRequest.AnchorPeers[o] {
					aPeer, ok := config.Peers[p]
					if ok {
						orgPeers = append(orgPeers, Peer{
							Host: aPeer.GRPCOptions["ssl-target-name-override"].(string),
							Port: common.GetPortFromUrl(aPeer.URL),
						})
					}
				}
			}

			orgs = append(orgs, Organization{
				MSPDir:      pathUtil.Substitute(config.OrgsByName[o].MSPDir),
				Name:        orgConfig.MSPID,
				ID:          orgConfig.MSPID,
				AnchorPeers: orgPeers,
			})
		}
	}

	return orgs
}

func (m *ChannelTxConfigHelper) getChannelOrderer(config *networkconfig.FabNetworkConfiguration, ordName string) (Orderer, bool) {
	Kafka := Kafka{
		Brokers: []string{
			"127.0.0.1:9092",
		},
	}

	ord, ok := config.Orderers[ordName]
	if !ok {
		return Orderer{}, false
	}

	ordAddress := []string{ord.URL}

	return Orderer{
		Kafka:       Kafka,
		Addresses:   ordAddress,
		OrdererType: "solo", // only one orderer supported for now
		BatchSize: BatchSize{
			MaxMessageCount:   20,
			AbsoluteMaxBytes:  99,
			PreferredMaxBytes: 512,
		},
		BatchTimeout: 2 * time.Second,
	}, true
}

//func (m *ChannelTxConfigHelper) getChannelOrganizations1(config *networkconfig.FabNetworkConfiguration, chOrgs []ChannelOrg) []Organization {
//	orgs := []Organization{}
//
//	for _, o := range chOrgs {
//		orgConfig, ok := config.Organizations[strings.ToLower(o.Name)]
//		if ok {
//			orgPeers := []Peer{}
//			for _, p := range o.AnchorPeers {
//				aPeer, ok := config.Peers[p]
//				if ok {
//					orgPeers = append(orgPeers, Peer{
//						Host: aPeer.GRPCOptions["ssl-target-name-override"].(string),
//						Port: common.GetPortFromUrl(aPeer.URL),
//					})
//				}
//			}
//			orgs = append(orgs, Organization{
//				MSPDir:      pathUtil.Substitute(config.OrgsByName[o.Name].MSPDir),
//				Name:        orgConfig.MSPID,
//				ID:          orgConfig.MSPID,
//				AnchorPeers: orgPeers,
//				isOrderer:   o.IsOrderer,
//				endpoint:    o.Endpoint,
//			})
//		}
//	}
//
//	return orgs
//}
