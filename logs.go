package resingo

import (
	"errors"
	"fmt"
	"io"

	"github.com/antonholmquist/jason"
	"github.com/pubnub/go/messaging"
)

//Logs  streams resin device logs
//
// Resin uses pubnub service for logs. Unfortunate the API for pubnub sucks big
// time. This Limits the API for this struct too.
type Logs struct {
	nub     *messaging.Pubnub
	channel string
	ctx     *Context
	stop    chan struct{}
}

//NewLogs returns a new Logs instace which is initialized to support streaming
//logs from pubnub.
func NewLogs(ctx *Context) (*Logs, error) {
	cfg, err := ConfigGetAll(ctx)
	if err != nil {
		return nil, err
	}
	//pretty.Println(cfg)
	if cfg.PubNub.PubKey != "" && cfg.PubNub.SubKey != "" {
		n := messaging.NewPubnub(
			cfg.PubNub.PubKey,
			cfg.PubNub.SubKey, "", "", false, "",
		)
		return &Logs{nub: n, ctx: ctx, stop: make(chan struct{})}, nil
	}
	return nil, errors.New("resingo: no pubnub details found")
}

//Subscribe subscribe to device logs
func (l *Logs) Subscribe(uuid string) (
	chan []byte, chan []byte, error,
) {
	logChan, err := l.GetChannel(uuid)
	if err != nil {
		return nil, nil, err
	}
	schan, echan := messaging.CreateSubscriptionChannels()
	l.nub.Subscribe(logChan, "", schan, false, echan)
	return schan, echan, nil

}

//GetChannel returns the device logs channel for the device with given uuid..
//The value is not cached, so this is stateless. This allows to use the same
//Lofs instance to syvscribe to multiple devices in different goroutines without
//race conditions.
func (l *Logs) GetChannel(uuid string) (string, error) {
	logsChan := uuid
	dev, err := DevGetByUUID(l.ctx, uuid)
	if err != nil {
		return "", err
	}
	if dev.LogsChannel != "" {
		logsChan = dev.LogsChannel
	}
	return fmt.Sprintf("device-%s-logs", logsChan), nil
}

//Log streams logs go out. This is blocking opretation, you should run this in a
//gorouting and call Log.Close when you are done acceping writes to out.
func (l *Logs) Log(uuid string, out io.Writer) error {
	s, e, err := l.Subscribe(uuid)
	if err != nil {
		return err
	}
stop:
	for {
		select {
		case rcv := <-s:
			nerr := l.write(out, rcv)
			if nerr != nil {
				fmt.Println(nerr)
			}
		case errrcv := <-e:
			err = errors.New(string(errrcv))
			break stop
		case <-l.stop:
			fmt.Println("resingo: stopping streaming logs")
			break stop
		}
	}
	l.nub.Abort()
	return err
}

func (l *Logs) write(out io.Writer, src []byte) error {
	a, _, _, err := l.nub.ParseJSON(src, "")
	if err != nil {
		return err
	}
	v, err := jason.NewValueFromBytes([]byte(a))
	if err != nil {
		return err
	}
	va, err := v.Array()
	if err != nil {
		return err
	}
	for _, value := range va {
		na, err := value.ObjectArray()
		if err != nil {
			return err
		}
		for _, vn := range na {
			m, err := vn.GetString("m")
			if err != nil {
				return err
			}
			fmt.Fprintf(out, " %s \n", m)
		}
	}
	return nil
}

//Close stops streaming device logs.
func (l *Logs) Close() {
	l.stop <- struct{}{}
}
