package counter

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

// Counter struct represents a virtual monotonic counter
type Counter struct {
	CounterID uint64
	Value     uint64
}

// Certificate struct represents the certificate returned by various operations
type Certificate struct {
	CounterID uint64
	Value     uint64
	Nonce     uint64
}

// VirtualMonotonicCounter struct represents the virtual monotonic counter mechanism
type VirtualMonotonicCounter struct {
	Counters map[uint64]*Counter
}

// CreateNewCounter creates a new virtual monotonic counter and returns a create certificate
func (vmc *VirtualMonotonicCounter) CreateNewCounter(nonce uint64) *Certificate {
	counterID := generateCounterID()
	newCounter := &Counter{CounterID: counterID, Value: 0}
	vmc.Counters[counterID] = newCounter
	return &Certificate{CounterID: counterID, Value: newCounter.Value, Nonce: nonce}
}

// ReadCounter returns a read certificate containing the current value of the virtual counter
func (vmc *VirtualMonotonicCounter) ReadCounter(counterID, nonce uint64) *Certificate {
	counter, exists := vmc.Counters[counterID]
	if !exists {
		return nil
	}
	return &Certificate{CounterID: counter.CounterID, Value: counter.Value, Nonce: nonce}
}

// IncrementCounter increments the specified virtual counter and returns an increment certificate
func (vmc *VirtualMonotonicCounter) IncrementCounter(counterID, nonce uint64) *Certificate {
	counter, exists := vmc.Counters[counterID]
	if !exists {
		return nil
	}
	counter.Value++
	return &Certificate{CounterID: counter.CounterID, Value: counter.Value, Nonce: nonce}
}

// DestroyCounter destroys the specified virtual counter and returns a destroy certificate
func (vmc *VirtualMonotonicCounter) DestroyCounter(counterID, nonce uint64) *Certificate {
	counter, exists := vmc.Counters[counterID]
	if !exists {
		return nil
	}
	delete(vmc.Counters, counterID)
	return &Certificate{CounterID: counter.CounterID, Nonce: nonce}
}

// generateCounterID generates a unique counter ID using crypto/rand
func generateCounterID() uint64 {
	var idBytes [8]byte
	_, err := rand.Read(idBytes[:])
	if err != nil {
		panic(err)
	}
	return binary.LittleEndian.Uint64(idBytes[:])
}

func main() {
	// Example usage of the virtual monotonic counter mechanism
	vmc := VirtualMonotonicCounter{Counters: make(map[uint64]*Counter)}

	nonce := uint64(123)

	// Create a new counter
	createCert := vmc.CreateNewCounter(nonce)
	fmt.Printf("Create Certificate: %+v\n", createCert)

	// Read the value of the created counter
	readCert := vmc.ReadCounter(createCert.CounterID, nonce)
	fmt.Printf("Read Certificate: %+v\n", readCert)

	// Increment the counter
	incrementCert := vmc.IncrementCounter(createCert.CounterID, nonce)
	fmt.Printf("Increment Certificate: %+v\n", incrementCert)

	// Destroy the counter
	destroyCert := vmc.DestroyCounter(createCert.CounterID, nonce)
	fmt.Printf("Destroy Certificate: %+v\n", destroyCert)
}
