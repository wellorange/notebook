

##bytes

  Buffer模块主要是使用在缓存读取的数据场景,就是先把数据缓存到buf中然后通过各种io接口花样读取
```
 type Buffer struct {
	buf       []byte 读取的缓存数据源
	off       int    读取的游标
}
```
   预定义的错误常量
```
// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("bytes.Buffer: too large")
var errNegativeRead = errors.New("bytes.Buffer: reader returned negative count from Read")
```


##buffer
##reader