package message

const BufferLength = 256

type Message struct {
	Data [BufferLength]byte
}

// Create a new Message with Empty Buffer
func NewEmptyMessage() *Message {
	return &Message{}
}

// Create a new Message Buffer filled with 0..BufferLength
func NewMessage() *Message {
	buffer := [BufferLength]byte{}

	for i := 0; i < BufferLength; i++ {
		buffer[i] = byte(uint16(i))
	}

	return &Message{
		Data: buffer,
	}
}

// Create a new Message using a provided buffer
func NewMessageFromBuffer(buffer [BufferLength]byte) *Message {
	return &Message{
		Data: buffer,
	}
}

// Verifies that a provided Message contains valid buffer data
func (m *Message) Checksum() bool {
	var total uint16 = 0
	for i := 0; i < BufferLength; i++ {
		total += uint16(m.Data[i])
	}

	// BufferLength - 1; 256 overflows on uint16, 0...255
	return total == ((BufferLength - 1) * ((BufferLength - 1) + 1) / 2)
}

func (m *Message) GetChecksum() uint16 {
	var total uint16 = 0
	for i := 0; i < BufferLength; i++ {
		total += uint16(m.Data[i])
	}

	return total
}
