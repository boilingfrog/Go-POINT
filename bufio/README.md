## bufio

### 前言

最近操作文件，进行优化使用到了`bufio`。好像也不太了解这个，那么就梳理下，`bufio`的使用。

### 例子

我的场景：使用`xml`拼接了`office2003`的文档。写入到`buffer`，然后处理完了，转存到文件里面。

````go
type Buff struct {
	Buffer *bytes.Buffer
	Writer *bufio.Writer
}

// 初始化
func NewBuff() *Buff {
	b := bytes.NewBuffer([]byte{})
	return &Buff{
		Buffer: b,
		Writer: bufio.NewWriter(b),
	}
}

func (b *Buff) WriteString(str string) error {
	_, err := b.Writer.WriteString(str)
	return err
}

func (b *Buff) SaveAS(name string) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := b.Writer.Flush(); err != nil {
		return nil
	}

	_, err = b.Buffer.WriteTo(file)
	return err
}

func main() {
	var b = NewBuff()

	b.WriteString("haah")
}
````

### bufio

> Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer object, creating another object (Reader or Writer) that also implements the interface but provides buffering and some help for textual I/O.  

`bufio`包实现了有缓冲的`I/O`。它包装一个`io.Reader`或`io.Writer`接口对象，创建另一个也实现了该接口，且同时还提供了缓冲和一些文本`I/O`的帮助函数的对象。  

简单的说就是`bufio`会把文件内容读取到缓存中（内存），然后再取读取需要的内容的时候，直接在缓存中读取，避免文件的`i/o`操作。同样，通过`bufio`写入内容，也是先写入到缓存中（内存），然后由缓存写入到文件。避免多次小内容的写入操作`I/O`。  

读取  

![bufio](images/bufio-read.png)  

写入  

![bufio](images/bufio-write.png?raw=true)


#### 源码解析

#### Reader对象

bufio.Reader 是bufio中对io.Reader 的封装

````go
// Reader implements buffering for an io.Reader object.
type Reader struct {
	buf          []byte 
	rd           io.Reader // 底层的io.Reader
	r, w         int       // r:从buf中读走的字节（偏移）；w:buf中填充内容的偏移； 
	                       // w - r 是buf中可被读的长度（缓存数据的大小），也是Buffered()方法的返回值
	err          error
	lastByte     int // 最后一次读到的字节（ReadByte/UnreadByte)
	lastRuneSize int // 最后一次读到的Rune的大小(ReadRune/UnreadRune)
}
````

bufio.Read(p []byte) 的思路如下：  

1、如果缓冲区有内容,直接读取内容到p中。
2、如果缓存区没有内容，并且读取的内容大于缓存的大小，直接不经过缓存直接读取文件。
3、如果缓存没有内容则，并且读取的内容小于缓存的大小，写入文件到缓存。当写完

````go
// Read reads data into p.
// It returns the number of bytes read into p.
// The bytes are taken from at most one Read on the underlying Reader,
// hence n may be less than len(p).
// To read exactly len(p) bytes, use io.ReadFull(b, p).
// At EOF, the count will be zero and err will be io.EOF.
func (b *Reader) Read(p []byte) (n int, err error) {
	n = len(p)
	if n == 0 {
		if b.Buffered() > 0 {
			return 0, nil
		}
		return 0, b.readErr()
	}
	// r:从buf中读走的字节（偏移）；w:buf中填充内容的偏移；
	// w - r 是buf中可被读的长度（缓存数据的大小），也是Buffered()方法的返回值
	// b.r == b.w 表示，当前缓冲区里面没有内容
	if b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		// 如果p的大小大于等于缓冲区大小，则直接将数据读入p，然后返回
		if len(p) >= len(b.buf) {
			// Large read, empty buffer.
			// Read directly into p to avoid copy.
			n, b.err = b.rd.Read(p)
			if n < 0 {
				panic(errNegativeRead)
			}
			if n > 0 {
				b.lastByte = int(p[n-1])
				b.lastRuneSize = -1
			}
			return n, b.readErr()
		}
		// One read.
		// Do not use b.fill, which will loop.
		b.r = 0
		b.w = 0
		n, b.err = b.rd.Read(b.buf)
		if n < 0 {
			panic(errNegativeRead)
		}
		if n == 0 {
			return 0, b.readErr()
		}
		b.w += n
	}

	// copy缓存区的内容到p中
	// copy as much as we can
	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	b.lastByte = int(b.buf[b.r-1])
	b.lastRuneSize = -1
	return n, nil
}
````








