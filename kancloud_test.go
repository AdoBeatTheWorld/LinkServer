package LinkServer

import "testing"

func TestSyncMap(t *testing.T) {
	m := make(map[int]int)
	go func() {
		for {
			_ = m [1]
		}
	}()
	go func() {
		for {
			m[2] = 2
		}
	}()
	select {

	}
}
