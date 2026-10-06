// Command grpcframe 是消息帧解析的场景入口。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.com/grpc-frame/grpcframe"
)

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func frame(flag byte, payload []byte) []byte {
	out := []byte{flag, 0, 0, 0, byte(len(payload))}
	return append(out, payload...)
}

func main() {
	sample := flag.String("sample", "", "split|maxlen|partial|work")
	flag.Parse()

	c := &grpcframe.Counter{}
	p := &grpcframe.Parser{}
	switch *sample {
	case "split":
		data := frame(0, []byte("abcd"))
		p.Feed(data[:3], c)
		out, err := p.Feed(data[3:], c)
		if err != nil {
			emit(map[string]string{"error": err.Error()})
			return
		}
		emit(map[string]any{"frames": out})
	case "maxlen":
		head := []byte{0, 0x00, 0x20, 0x00, 0x00}
		out, err := p.Feed(head, c)
		if err != nil {
			emit(map[string]string{"error": err.Error()})
			return
		}
		emit(map[string]any{"frames": out})
	case "partial":
		out, err := p.Feed([]byte{0, 0, 0, 0, 10}, c)
		if err != nil {
			emit(map[string]string{"error": err.Error()})
			return
		}
		emit(map[string]any{"frames": out})
	case "work":
		data := []byte{0, 0, 0, 0, 1, 'x'}
		n := 0
		for i := 0; i < 2000; i++ {
			for j := 0; j < len(data); j += 3 {
				end := j + 3
				if end > len(data) {
					end = len(data)
				}
				out, _ := p.Feed(data[j:end], c)
				n += len(out)
			}
		}
		emit(map[string]int{"frames": n, "scanned": c.Scanned})
	default:
		fmt.Fprintln(os.Stderr, "未知场景："+*sample)
		os.Exit(2)
	}
}
