## []byte、bytes.Buffer 缓存池

```golang
data := leakybuf.Bytes.Get(size) // 获取size大小的缓存
// do sth
leakybuf.Bytes.Put(data) // 将缓存放回缓存池
```
