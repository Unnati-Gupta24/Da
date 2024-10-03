package da

import (
    "gopkg.in/zeromq/goczmq.v4"
)

// MessageProcessor defines an interface for processing messages.
type MessageProcessor interface {
    Process(msg [][]byte, protocolId string) ([]byte, error)
}

// ChannelSubscriber defines an interface for subscribing to a ZeroMQ channel.
type ChannelSubscriber interface {
    Subscribe(endpoint string, typ string) bool
    Listen(processor MessageProcessor, protocolId string)
}

// ZmqChannelReader implements ChannelSubscriber for reading messages.
type ZmqChannelReader struct {
    channeler *goczmq.Channeler
}

func (zmqC *ZmqChannelReader) Subscribe(endpoint string, typ string) bool {
    zmqC.channeler = goczmq.NewSubChanneler(endpoint, typ)
    if zmqC.channeler == nil {
        return false
    }
    return true
}

func (zmqC *ZmqChannelReader) Listen(processor MessageProcessor, protocolId string) {
    for {
        msg, ok := <-zmqC.channeler.RecvChan
        if !ok || len(msg) != 3 {
            continue // Handle error or unexpected message format.
        }
        processor.Process(msg, protocolId)
    }
}
