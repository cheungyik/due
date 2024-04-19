package packet

import (
	"bytes"
	"encoding/binary"
	"github.com/dobyte/due/v2/errors"
	"io"
	"sync"
)

const (
	unbindReqBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + 8
	unbindResBytes = defaultSizeBytes + defaultHeaderBytes + defaultRouteBytes + defaultSeqBytes + defaultCodeBytes
)

type UnbindPacker struct {
	reqPool  *sync.Pool
	resPool  *sync.Pool
	reqPool2 *sync.Pool
}

func NewUnbindPacker() *UnbindPacker {
	p := &UnbindPacker{}
	p.reqPool = &sync.Pool{}
	p.reqPool.New = func() any { return NewBuffer(p.reqPool, unbindReqBytes) }
	p.resPool = &sync.Pool{}
	p.resPool.New = func() any { return NewBuffer(p.resPool, unbindResBytes) }
	p.reqPool2 = &sync.Pool{}
	p.reqPool2.New = func() any { return NewWriter(p.reqPool2, unbindReqBytes) }

	return p
}

// PackReq 打包请求
// 协议格式：size + header + route + seq + uid
func (p *UnbindPacker) PackReq(seq uint64, uid int64) (buf *Buffer, err error) {
	buf = p.reqPool.Get().(*Buffer)
	defer func() {
		if err != nil {
			buf.Recycle()
		}
	}()

	size := unbindReqBytes - defaultSizeBytes

	if err = binary.Write(buf, binary.BigEndian, int32(size)); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, dataBit); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, unbindReq); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, seq); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, uid); err != nil {
		return
	}

	return
}

func (p *UnbindPacker) PackReq2(seq uint64, uid int64) (writer *Writer, err error) {
	writer = p.reqPool2.Get().(*Writer)
	defer func() {
		if err != nil {
			writer.Recycle()
		}
	}()

	size := unbindReqBytes - defaultSizeBytes

	writer.WriteInt32s(binary.BigEndian, int32(size))
	writer.WriteUint8s(dataBit)
	writer.WriteInt8s(unbindReq)
	writer.WriteUint64s(binary.BigEndian, seq)
	writer.WriteInt64s(binary.BigEndian, uid)

	return
}

// UnpackReq 解包请求
// 协议格式：size + header + route + seq + uid
func (p *UnbindPacker) UnpackReq(data []byte) (seq uint64, uid int64, err error) {
	if len(data) != unbindReqBytes {
		err = errors.ErrInvalidMessage
		return
	}

	reader := bytes.NewReader(data)

	if _, err = reader.Seek(defaultSizeBytes+defaultHeaderBytes+defaultRouteBytes, io.SeekStart); err != nil {
		return
	}

	if err = binary.Read(reader, binary.BigEndian, &seq); err != nil {
		return
	}

	if err = binary.Read(reader, binary.BigEndian, &uid); err != nil {
		return
	}

	return
}

// PackRes 打包响应
// size + header + route + seq + code
func (p *UnbindPacker) PackRes(seq uint64, code int16) (buf *Buffer, err error) {
	buf = p.resPool.Get().(*Buffer)
	defer func() {
		if err != nil {
			buf.Recycle()
		}
	}()

	size := unbindResBytes - defaultSizeBytes

	if err = binary.Write(buf, binary.BigEndian, int32(size)); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, dataBit); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, unbindRes); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, seq); err != nil {
		return
	}

	if err = binary.Write(buf, binary.BigEndian, code); err != nil {
		return
	}

	return
}

// UnpackRes 解包响应
// size + header + route + seq + code
func (p *UnbindPacker) UnpackRes(data []byte) (code int16, err error) {
	if len(data) != unbindResBytes {
		err = errors.ErrInvalidMessage
		return
	}

	reader := bytes.NewReader(data)

	if _, err = reader.Seek(-defaultCodeBytes, io.SeekEnd); err != nil {
		return
	}

	if err = binary.Read(reader, binary.BigEndian, &code); err != nil {
		return
	}

	return
}
