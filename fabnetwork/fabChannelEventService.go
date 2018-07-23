package fabnetwork

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
//func (m *FabricNetwork) RegisterBlockEventListener(channelName, orgName string, receiver chan uint64) {
//	reg:= channelEventRegistration{
//		channelName: channelName,
//		orgName: orgName,
//	}
//
//	m.channelEventService[reg]= append(m.channelEventService[reg], receiver)
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