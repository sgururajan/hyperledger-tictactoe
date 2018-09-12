package ex2

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type TicketingWindow struct {
	WindowKey        string
	GivesWaterBottle bool
}

type MovieShow struct {
	Key       string
	MovieName string
	ShowTime  time.Time
	Hall      Hall
}

type Hall struct {
	Key      string
	Capacity int
}

type Theater struct {
	TheaterKey       string
	Name             string
	Halls            []Hall
	CurrentlyShowing []MovieShow
	TicketingWindows []TicketingWindow
	Cafetaria        Cafetaria
}

type Ticket struct {
	TicketKey           string
	TheaterKey          string
	ShowKey             string
	MovieName           string
	ShowTime            time.Time
	NoOfTicket          int
	IssuedWindow        string
	CafetariaCouponCode int
	CouponRedeemed      bool
}

type Cafetaria struct {
	SodaAvailable int
}

type MovieTicket struct {
}

func main() {
	err := shim.Start(new(MovieTicket))
	if err != nil {
		fmt.Println("error starting tictactoe chaincode")
		fmt.Printf("%#v", err)
	}
}

func (m *MovieTicket) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// intialize just one theater during intialization

	theater := getTheaterByKey(1, stub)
	if theater == nil {
		halls := []Hall{
			Hall{Capacity: 100, Key: strconv.Itoa(1)},
			Hall{Capacity: 100, Key: strconv.Itoa(2)},
			Hall{Capacity: 100, Key: strconv.Itoa(3)},
			Hall{Capacity: 100, Key: strconv.Itoa(4)},
			Hall{Capacity: 100, Key: strconv.Itoa(5)},
		}
		shows:= []MovieShow{
			MovieShow{Key:strconv.Itoa(1), Hall:halls[0], ShowTime: getToday().Add(time.Hour*10), MovieName:"Avengers"},
			MovieShow{Key:strconv.Itoa(2), Hall:halls[0], ShowTime: getToday().Add(time.Hour*12), MovieName:"Avengers"},
			MovieShow{Key:strconv.Itoa(3), Hall:halls[0], ShowTime: getToday().Add(time.Hour*14), MovieName:"Avengers"},
			MovieShow{Key:strconv.Itoa(4), Hall:halls[0], ShowTime: getToday().Add(time.Hour*16), MovieName:"Avengers"},

			MovieShow{Key:strconv.Itoa(5), Hall:halls[1], ShowTime: getToday().Add(time.Hour*10), MovieName:"MI6"},
			MovieShow{Key:strconv.Itoa(6), Hall:halls[1], ShowTime: getToday().Add(time.Hour*12), MovieName:"MI6"},
			MovieShow{Key:strconv.Itoa(7), Hall:halls[1], ShowTime: getToday().Add(time.Hour*14), MovieName:"MI6"},
			MovieShow{Key:strconv.Itoa(8), Hall:halls[1], ShowTime: getToday().Add(time.Hour*16), MovieName:"MI6"},

			MovieShow{Key:strconv.Itoa(9), Hall:halls[2], ShowTime: getToday().Add(time.Hour*10), MovieName:"StarWars"},
			MovieShow{Key:strconv.Itoa(10), Hall:halls[2], ShowTime: getToday().Add(time.Hour*12), MovieName:"StarWars"},
			MovieShow{Key:strconv.Itoa(11), Hall:halls[2], ShowTime: getToday().Add(time.Hour*14), MovieName:"StarWars"},
			MovieShow{Key:strconv.Itoa(12), Hall:halls[2], ShowTime: getToday().Add(time.Hour*16), MovieName:"StarWars"},

			MovieShow{Key:strconv.Itoa(13), Hall:halls[3], ShowTime: getToday().Add(time.Hour*10), MovieName:"StarWars"},
			MovieShow{Key:strconv.Itoa(14), Hall:halls[3], ShowTime: getToday().Add(time.Hour*12), MovieName:"StarWars"},
			MovieShow{Key:strconv.Itoa(15), Hall:halls[3], ShowTime: getToday().Add(time.Hour*14), MovieName:"StarWars"},
			MovieShow{Key:strconv.Itoa(16), Hall:halls[3], ShowTime: getToday().Add(time.Hour*16), MovieName:"StarWars"},

			MovieShow{Key:strconv.Itoa(17), Hall:halls[4], ShowTime: getToday().Add(time.Hour*10), MovieName:"Gravity"},
			MovieShow{Key:strconv.Itoa(18), Hall:halls[4], ShowTime: getToday().Add(time.Hour*12), MovieName:"Gravity"},
			MovieShow{Key:strconv.Itoa(19), Hall:halls[4], ShowTime: getToday().Add(time.Hour*14), MovieName:"Gravity"},
			MovieShow{Key:strconv.Itoa(20), Hall:halls[4], ShowTime: getToday().Add(time.Hour*16), MovieName:"Gravity"},
		}
		theater = &Theater{
			TheaterKey: strconv.Itoa(1),
			TicketingWindows: []TicketingWindow{
				{strconv.Itoa(1), true},
				{strconv.Itoa(2), false},
				{strconv.Itoa(3), false},
				{strconv.Itoa(4), false},
			},
			Halls: halls,
			CurrentlyShowing:shows,
			Cafetaria:Cafetaria{
				SodaAvailable: 200,
			},
			Name: "Cinemark",
		}

		stubBytes,err:= json.Marshal(theater)
		if err != nil {
			return shim.Error("error creating default theatre")
		}

		stub.PutState(theater.TheaterKey, stubBytes)
	}

	return shim.Success(nil)
}

func (m *MovieTicket) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, _ := stub.GetFunctionAndParameters()
	switch strings.ToLower(fn) {
	case "issueticket":
		return issueTicket(stub)
	default:
		shim.Error("invalid invoke function")
	}
	return shim.Success(nil)
}

func issueTicket(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	// expected arguments: theaterKey, movieName, showTime, noOfTickets, windowKey
	if len(args) < 5 {
		return shim.Error("minimum 5 arguments expected")
	}

	theaterKey := args[0]
	movieName := args[1]
	showTime, err := time.Parse("yyyy/MM/dd:HH:mm", args[2])
	if err != nil {
		return shim.Error(fmt.Sprintf("error while parsing show time. err: %#v", err))
	}

	noOfTickets, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("invalid no of tickets")
	}

	windowkey := args[4]

	stubBytes, err := stub.GetState(theaterKey)
	if err != nil {
		return shim.Error(fmt.Sprintf("error while getting theater state. Err: %#v", err))
	}

	theater, err := getTheater(stubBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("error while parsing theater information. Err: %#v", err))
	}

	var show *MovieShow
	for _, s := range theater.CurrentlyShowing {
		if strings.EqualFold(s.MovieName, movieName) && s.ShowTime == showTime {
			show = &s
			break
		}
	}

	if show.Hall.Capacity < noOfTickets {
		return shim.Error("not enough seats available")
	}

	var window TicketingWindow
	for _, w := range theater.TicketingWindows {
		if w.WindowKey == windowkey {
			window = w
			break
		}
	}

	ticket := Ticket{}
	ticket.MovieName = movieName
	ticket.ShowTime = show.ShowTime
	ticket.IssuedWindow = windowkey
	ticket.NoOfTicket = noOfTickets
	ticket.TheaterKey = theaterKey
	ticket.ShowKey = show.Key
	ticket.TicketKey = getRandomString()

	if window.GivesWaterBottle {
		ticket.CafetariaCouponCode = rand.Intn(9999-1001) + 1001
	}

	show.Hall.Capacity = show.Hall.Capacity - noOfTickets

	ticketBytes, err := json.Marshal(ticket)
	if err != nil {
		return shim.Error(fmt.Sprintf("error while issuing ticket. err: %#v", err))
	}
	stub.PutState(ticket.TicketKey, ticketBytes)

	stubBytes, err = json.Marshal(theater)
	if err != nil {
		return shim.Error(fmt.Sprintf("error while update theater state. err: %#v", err))
	}
	stub.PutState(theaterKey, stubBytes)

	return shim.Success(ticketBytes)

}

func getTheaterByKey(key int, stub shim.ChaincodeStubInterface) *Theater {
	tKey := strconv.Itoa(key)
	ticketBytes, err := stub.GetState(tKey)
	if err != nil {
		return nil
	}

	if len(ticketBytes) == 0 {
		return nil
	}

	theater, err := getTheater(ticketBytes)
	if err != nil {
		return nil
	}
	return &theater
}

func getTheater(data []byte) (Theater, error) {
	result := Theater{}
	err := json.Unmarshal(data, &result)
	return result, err
}

func getRandomString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := 7

	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func getToday() time.Time {
	year,month,day:= time.Now().Date()
	return time.Date(year,month,day,0,0,0,0,time.Local)
}