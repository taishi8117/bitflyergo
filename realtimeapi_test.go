package bitflyergo

import (
	"log"
	//"os"
	//"os/signal"
	//"syscall"
	"testing"
	//"time"
)

type C struct{}

func (c *C) OnReceiveExecutions(channelName string, executions []Execution) {
	log.Println("OnReceiveExecutions")
}

func (c *C) OnReceiveBoard(channelName string, board *Board) {
	log.Println("OnReceiveBoard")
}

func (c *C) OnReceiveBoardSnapshot(channelName string, board *Board) {
	log.Println("OnReceiveBoardSnapshot")
}

func (c *C) OnReceiveTicker(channelName string, ticker *Ticker) {
	log.Println("OnReceiveTicker")
}

func (c *C) OnReceiveChildOrderEvents(channelName string, event []ChildOrderEvent) {
	log.Println("OnReceiveExecutions")
}

func (c *C) OnReceiveParentOrderEvents(channelName string, event []ParentOrderEvent) {
	log.Println("OnReceiveParentOrderEvents")
}

func (c *C) OnErrorOccur(channelName string, err error) {
	log.Println("OnErrorOccur")
}

func TestReceive(t *testing.T) {
	//interrupt := make(chan os.Signal, 1)
	//	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	//	ws := WebSocketClient{
	//		Debug: false,
	//		Cb:    &C{},
	//	}
	//	err := ws.Connect()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	go ws.Receive()
	//	ws.SubscribeExecutions(ProductCodeFxBtcJpy)
	//	time.Sleep(2 * time.Second)
	//
	//LOOP:
	//	for {
	//		select {
	//		case _ = <-interrupt:
	//			break LOOP
	//		}
	//	}
}
