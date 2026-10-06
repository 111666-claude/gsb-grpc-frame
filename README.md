# grpc-frame

按 gRPC 的长度前缀帧格式切分字节流。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/grpcframe --sample=split
go run ./cmd/grpcframe --sample=maxlen
go run ./cmd/grpcframe --sample=partial
go run ./cmd/grpcframe --sample=work
```

## 口径

- **帧头**：5 字节，第 1 字节是压缩标志，后 4 字节是负载长度，按大端读。
- **跨块**：帧可以被任意切分，头与负载分在不同块里也要接得上。
- **负载**：负载没到齐就继续等，不出帧。
- **上限**：单帧负载超过 1 MiB 直接报 `too-large`。
- **收尾**：流结束时残留的半帧丢弃，不计入结果。

## 不变量

- 出的帧负载长度与头里声明的一致。
- 切块方式不影响帧序列。
- `scanned`（重复扫描的缓冲字节数）不随帧数放大：两千帧不超过 9000。

## 输出契约

每个场景打印一行 JSON，`frames` 里每项含 `compressed`、`len`、`payload`；
出错时打印 `{"error": "..."}`；`--sample=work` 只打印 `{"frames": N, "scanned": N}`。
