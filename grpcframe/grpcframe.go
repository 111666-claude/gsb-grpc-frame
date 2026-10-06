// Package grpcframe 解析 gRPC 的长度前缀消息帧。
package grpcframe

import (
	"encoding/binary"
	"errors"
)

// MaxLen 是单帧负载上限。
const MaxLen = 1 << 20

// ErrTooLarge 表示单帧负载超过 MaxLen。
var ErrTooLarge = errors.New("too-large")

// Frame 是一条解析出来的消息帧。
type Frame struct {
	Compressed bool `json:"compressed"`
	Len        int  `json:"len"`
	Payload    int  `json:"payload"`
}

// Counter 记录本轮重复扫描过的缓冲字节数。
type Counter struct {
	Scanned int
}

// Parser 累积字节并吐出完整的帧。
type Parser struct {
	buf []byte
}

// Feed 吃进一段字节，返回这段字节里凑齐的帧。
func (p *Parser) Feed(chunk []byte, c *Counter) ([]Frame, error) {
	c.Scanned += len(p.buf)
	p.buf = append(p.buf, chunk...)
	out := make([]Frame, 0)
	for len(p.buf) >= 5 {
		n := int(binary.BigEndian.Uint32(p.buf[1:5]))
		if n > MaxLen {
			p.buf = p.buf[:0]
			return nil, ErrTooLarge
		}
		if len(p.buf) < 5+n {
			break
		}
		out = append(out, Frame{Compressed: p.buf[0] == 1, Len: n, Payload: n})
		p.buf = p.buf[5+n:]
	}
	if len(p.buf) == 0 {
		p.buf = nil
	}
	return out, nil
}

// Close 结束流，返回流尾残留的半帧数量。
func (p *Parser) Close() int {
	n := len(p.buf)
	p.buf = p.buf[:0]
	return n
}
