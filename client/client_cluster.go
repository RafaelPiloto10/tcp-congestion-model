package main

import "sync"

func main() {

	var totalPackets int64 = 100_000_000_000
	var sentPackets int64 = 0
	var sendPackets int64 = 100_000_000

	packetsLock := sync.Mutex{}	
	wg := sync.WaitGroup{}

	for i := 0; int64(i) < totalPackets / sendPackets; i++ {
		wg.Add(1)
		go func() {
			total := RunClient(sendPackets)
			packetsLock.Lock()
			defer packetsLock.Unlock()

			sentPackets += int64(total)

			wg.Done()
		}()
	}

	wg.Wait()
}
