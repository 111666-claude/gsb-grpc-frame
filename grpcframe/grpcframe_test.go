package grpcframe

import "testing"

func build(flag byte, payload []byte) []byte {
	out := []byte{flag, 0, 0, 0, byte(len(payload))}
	return append(out, payload...)
}

func TestSingleFrame(t *testing.T) {
	p := &Parser{}
	out, err := p.Feed(build(0, []byte("abcd")), &Counter{})
	if err != nil {
		t.Fatalf("解析失败：%v", err)
	}
	if len(out) != 1 || out[0].Len != 4 || out[0].Payload != 4 {
		t.Fatalf("应该出一条 len=4/payload=4 的帧：%+v", out)
	}
}

func TestCompressedFlag(t *testing.T) {
	p := &Parser{}
	out, _ := p.Feed(build(1, []byte("x")), &Counter{})
	if len(out) != 1 || !out[0].Compressed {
		t.Fatalf("标志位应解析成压缩：%+v", out)
	}
}

func TestShortChunkGivesNothing(t *testing.T) {
	p := &Parser{}
	out, _ := p.Feed([]byte{0, 0, 0}, &Counter{})
	if len(out) != 0 {
		t.Fatalf("不足一个头应不出帧：%+v", out)
	}
}

func TestEmptyChunkGivesNothing(t *testing.T) {
	p := &Parser{}
	out, _ := p.Feed(nil, &Counter{})
	if len(out) != 0 {
		t.Fatalf("空块应不出帧：%+v", out)
	}
}

func TestCloseReportsLeftover(t *testing.T) {
	p := &Parser{}
	p.Feed(build(0, []byte("abcd")), &Counter{})
	if got := p.Close(); got != 0 {
		t.Fatalf("完整帧后不该有残留，得到 %d", got)
	}
}
