package agave

import (
	hpfeeds "github.com/d1str0/go-hpfeeds"
)

type HpfeedsWriter struct {
	client      hpfeeds.Client
	publish     chan []byte
	channelName string
	currentErr  error
	errLock     RWMutex
}

func NewHpfeedsWriter(host string, port int, ident string, auth string, channel string) (*HpfeedsWriter, error) {
	p := make(chan []byte)
	c := hpfeeds.NewClient(host, port, ident, auth)
	err := c.Connect()
	if err != nil {
		return nil, error
	}

	c.Publish(channel, p)

	w := HpfeedsWriter{client: c, publish: p, channelName: channel}

	go func() {
		for {
			<-client.Disconnected
			fmt.Printf("Attempting to reconnect...\n")
			err = client.Connect()
			if err != nil {
				fmt.Printf("Error reconnecting: %s\n", err.Error())
				w.errLock.Lock()
				w.currentErr = err
				w.errLock.Unlock()
				recon = time.After(5 * time.Second)
				<-recon
			}
		}
	}()

	return w
}

func (w HpfeedsWriter) Write(p []byte) (n int, err error) {
	w.publish <- p
	w.errLock.RLock()
	err := w.currentErr
	w.errLock.RUnlock()

	return len(p), err
}
