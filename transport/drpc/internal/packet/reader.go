package packet

import (
	"encoding/binary"
	"github.com/dobyte/due/v2/errors"
	"io"
	"sync"
)

type Reader struct {
	sizePool sync.Pool
}

func NewReader() *Reader {
	return &Reader{sizePool: sync.Pool{New: func() any { return make([]byte, 4) }}}
}

// ReadMessage 读取消息
func (r *Reader) ReadMessage(reader io.Reader) (isHeartbeat bool, route int8, seq uint64, data []byte, err error) {
	buf := r.sizePool.Get().([]byte)

	if _, err = io.ReadFull(reader, buf); err != nil {
		r.sizePool.Put(buf)
		return
	}

	size := binary.BigEndian.Uint32(buf)

	if size == 0 {
		r.sizePool.Put(buf)
		err = errors.ErrInvalidMessage
		return
	}

	data = make([]byte, defaultSizeBytes+size)
	copy(data[:defaultSizeBytes], buf)

	r.sizePool.Put(buf)

	if _, err = io.ReadFull(reader, data[defaultSizeBytes:]); err != nil {
		return
	}

	header := data[defaultSizeBytes : defaultSizeBytes+defaultHeaderBytes][0]

	isHeartbeat = header&heartbeatBit == heartbeatBit

	if !isHeartbeat {
		route = int8(data[defaultSizeBytes+defaultHeaderBytes : defaultSizeBytes+defaultHeaderBytes+defaultRouteBytes][0])

		seq = binary.BigEndian.Uint64(data[defaultSizeBytes+defaultHeaderBytes+defaultRouteBytes : defaultSizeBytes+defaultHeaderBytes+defaultRouteBytes+8])
	}

	return
}
