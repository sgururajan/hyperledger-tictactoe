package fabnetwork

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/sgururajan/hyperledger-tictactoe/domainModel"
)

//import (
//	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
//	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
//	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/client"
//	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
//)
//
////type ChannelEvent struct {
////	ChannelName string
////	OrgName string
////}
////
////type ChannelEventService struct {
////	FabRegistration fab.Registration
////	BlockEventChannel <-chan *fab.FilteredBlockEvent
////}
////
////func NewChannelEventService(registration fab.Registration, eventChannel <-chan *fab.FilteredBlockEvent) *ChannelEventService {
////	return &ChannelEventService{
////		BlockEventChannel:eventChannel,
////		FabRegistration:registration,
////	}
////}
//
//type channelEventRegistration struct {
//	channelName string
//	orgName string
//}
//
//type channelEventService struct {
//	register chan channelEventRegistration
//	unregister chan channelEventRegistration
//}
//
//func newChannelEventService() *channelEventService {
//	return &channelEventService{
//		register: make(chan channelEventRegistration),
//		unregister: make(chan channelEventRegistration),
//	}
//}
//
//func (m *FabricNetwork) RunEventService(channelName, orgName string) {
//	log:= logging.NewLogger("ChannelEventService-RunEventService")
//
//	org:= m.orgsByName[orgName]
//	chContextProvider:= m.sdk.ChannelContext(channelName, fabsdk.WithOrg(orgName), fabsdk.WithUser(org.AdminUser))
//	chContext,err:= chContextProvider()
//	if err != nil {
//		log.Errorf("error while creating channel context. err: %#v", err)
//		return
//	}
//
//	eventService,err:= chContext.ChannelService().EventService(client.WithBlockEvents())
//	if err != nil {
//		log.Errorf("error while creating eventService. err: %#v", err)
//		return
//	}
//
//	reg, eventch, err:= eventService.RegisterFilteredBlockEvent()
//	if err != nil {
//		log.Errorf("error while registering for event. err: %#v", err)
//		return
//	}
//
//	defer eventService.Unregister(reg)
//	clientReg:= channelEventRegistration{channelName:channelName, orgName:orgName}
//	for {
//		select{
//		case blockEvent:= <- eventch:
//			for _,r:= range m.channelEventService[clientReg] {
//				r <- blockEvent.FilteredBlock.Number
//			}
//		}
//
//	}
//}
//
//func (m *FabricNetwork) RegisterBlockEventListener(channelName, orgName string, Receiver chan uint64) {
//	reg:= channelEventRegistration{
//		channelName: channelName,
//		orgName: orgName,
//	}
//
//	m.channelEventService[reg]= append(m.channelEventService[reg], Receiver)
//}
//
//func (m *FabricNetwork) UnregisterBlockEventListener() {
//
//}
//
//func (m *FabricNetwork) registerEventServiceForChannel(channelName, orgName string) {
//	org:= m.orgsByName[orgName]
//	chContextProvider:= m.sdk.ChannelContext(channelName, fabsdk.WithUser(org.AdminUser), fabsdk.WithOrg(orgName))
//	chContext, _:= chContextProvider()
//	eventService,_:= chContext.ChannelService().EventService(client.WithBlockEvents())
//	reg, eventCh, _:= eventService.RegisterFilteredBlockEvent()
//	channelEvent:= ChannelEvent{
//		OrgName:orgName,
//		ChannelName:channelName,
//	}
//
//
//	m.channelEventService[channelEvent] = NewChannelEventService(reg, eventCh)
//
//}
//
//func runEventServiceChannel(fabEventService fab2.EventService, eventService ChannelEventService) {
//	defer fabEventService.Unregister(eventService.FabRegistration)
//
//	for {
//		select {
//		case blockEvent:= <- eventService.BlockEventChannel
//
//		}
//	}
//}

type BlockEventListener struct {
	Receiver chan domainModel.BlockInfo
}

type BlockEventService struct {
	channelName string
	orgName     string
	register    chan BlockEventListener
	unRegister  chan BlockEventListener
}

func (m *FabricNetwork) registerEventServiceForChannel(channelName, orgName string) {
	log := logging.NewLogger("registerEventServiceForChannel")
	listeners := make(map[BlockEventListener]bool)
	bEventService := BlockEventService{
		channelName: channelName,
		orgName:     orgName,
		register:    make(chan BlockEventListener),
		unRegister:  make(chan BlockEventListener),
	}
	serviceCompositeKey := fmt.Sprintf("%s-%s", channelName, orgName)

	defer close(bEventService.register)
	defer close(bEventService.unRegister)
	defer delete(m.channelBlockEventService, serviceCompositeKey)

	m.channelBlockEventService[serviceCompositeKey] = bEventService

	org := m.orgsByName[orgName]
	contextProvider := m.sdk.ChannelContext(channelName, fabsdk.WithUser(org.AdminUser), fabsdk.WithOrg(orgName))
	chContext, err := contextProvider()
	if err != nil {
		log.Errorf("error while creating channel context. Err: %#v", err)
		return
	}

	eventService, err := chContext.ChannelService().EventService(client.WithBlockEvents())
	if err != nil {
		log.Errorf("error while creating event service. err: %#v", err)
		return
	}

	reg, eventCh, err := eventService.RegisterFilteredBlockEvent()
	if err != nil {
		log.Errorf("error while registering for filtered block event. err: %#v", err)
		return
	}

	defer eventService.Unregister(reg)

	for {
		select {
		case listener := <-bEventService.register:
			listeners[listener] = true
		case listener := <-bEventService.unRegister:
			if _, ok := listeners[listener]; ok {
				delete(listeners, listener)
				close(listener.Receiver)
			}
		case bEvent := <-eventCh:
			blockInfo:= getDomainBlockInfo(bEvent)
			for l, _ := range listeners {
				select {
				case l.Receiver <- blockInfo:
				default:
					close(l.Receiver)
					delete(listeners, l)
				}
			}
		}
	}
}

func getDomainBlockInfo(bEvent *fab.FilteredBlockEvent) domainModel.BlockInfo {
	var result domainModel.BlockInfo
	result.BlockNumber = bEvent.FilteredBlock.Number
	result.ChannelId = bEvent.FilteredBlock.ChannelId
	result.Source = bEvent.SourceURL
	result.NoOfTransactions = len(bEvent.FilteredBlock.FilteredTransactions)
	result.Transactions = []domainModel.BlockTransaction{}

	for _,t:= range bEvent.FilteredBlock.FilteredTransactions {
		var bTran domainModel.BlockTransaction
		bTran.Type = t.Type.String()
		bTran.TxId = t.Txid
		bTran.ValidationCode=t.TxValidationCode.String()
		result.Transactions = append(result.Transactions, bTran)
	}

	return result
}

func (m *FabricNetwork) RegisterBlockEventListener(channelName, orgName string, eventListener BlockEventListener) error {
	serviceCompositeKey := fmt.Sprintf("%s-%s", channelName, orgName)
	if _,ok:= m.channelBlockEventService[serviceCompositeKey];!ok {
		return errors.New(fmt.Sprintf("no event service running for channel %s and org %s", channelName, orgName))
	}

	m.channelBlockEventService[serviceCompositeKey].register <- eventListener
	return nil
}

func (m *FabricNetwork) UnRegisterBlockEventListener(channelName, orgName string, eventListener BlockEventListener) {
	serviceCompositeKey := fmt.Sprintf("%s-%s", channelName, orgName)
	if _,ok:= m.channelBlockEventService[serviceCompositeKey];ok {
		m.channelBlockEventService[serviceCompositeKey].unRegister <- eventListener
	}
}
