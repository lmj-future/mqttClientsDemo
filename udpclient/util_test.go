package udpclient

import (
 
	"testing"
	"time"
)

func TestXxx(t *testing.T) {
	res := make(chan bool)
	go func() {
		for {
			select {
			case <-res:
				return
			default:
				continue
			}
		}

	}()
	time.Sleep(time.Second * 15)
	res <- true
}
