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
3、如果缓存没有内容则，并且读取的内容小于缓存的大小，写入文件到缓存。然后重复1。

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

##### 实例化

`bufio` 包提供了两个实例化 `bufio.Reader` 对象的函数：`NewReader` 和 `NewReaderSize`。其中，`NewReader` 函数是调用 `NewReaderSize`。  
函数实现的：  

````
// NewReader returns a new Reader whose buffer has the default size.
func NewReader(rd io.Reader) *Reader {
    // defaultBufSize = 4096,默认的大小
	return NewReaderSize(rd, defaultBufSize)
}
````

调用的`NewReaderSize`

````go
// NewReaderSize returns a new Reader whose buffer has at least the specified
// size. If the argument io.Reader is already a Reader with large enough
// size, it returns the underlying Reader.
func NewReaderSize(rd io.Reader, size int) *Reader {
	// Is it already a Reader?
	b, ok := rd.(*Reader)
	if ok && len(b.buf) >= size {
		return b
	}
	if size < minReadBufferSize {
		size = minReadBufferSize
	}
	r := new(Reader)
	r.reset(make([]byte, size), rd)
	return r
}
````

##### ReadSlice


// ReadSlice reads until the first occurrence of delim in the input,  
// returning a slice pointing at the bytes in the buffer.  
// The bytes stop being valid at the next read.  
// If ReadSlice encounters an error before finding a delimiter,  
// it returns all the data in the buffer and the error itself (often io.EOF).  
// ReadSlice fails with error ErrBufferFull if the buffer fills without a delim.  
// Because the data returned from ReadSlice will be overwritten  
// by the next I/O operation, most clients should use  
// ReadBytes or ReadString instead.  
// ReadSlice returns err != nil if and only if line does not end in delim.  

`ReadSlice`需要放置一个界定符号，来分割

````go
	reader := bufio.NewReader(strings.NewReader("hello \n world"))
	line1, _ := reader.ReadSlice('\n')
	fmt.Printf("the line1:%s\n", line1)

	line2, _ := reader.ReadSlice('\n')
	fmt.Printf("the line2:%s\n", line2)
````

输出

````
the line1:hello 

the line2: world
````

##### ReadString

`ReadString`是通过调用`ReadBytes`来实现的，看下源码：

````go
func (b *Reader) ReadString(delim byte) (string, error) {
	bytes, err := b.ReadBytes(delim)
	return string(bytes), err
}
````

使用例子：

````go
	reader := bufio.NewReader(strings.NewReader("hello \n world"))
	line1, _ := reader.ReadString('\n')
	fmt.Printf("the line1:%s\n", line1)

	line2, _ := reader.ReadString('\n')
	fmt.Printf("the line2:%s\n", line2)
````

##### ReadLine

根据官方的解释这个是不推荐使用的，推荐使用`ReadBytes('\n')` or `ReadString('\n')`来替代。  

ReadLine尝试返回单独的行，不包括行尾的换行符。如果一行大于缓存，isPrefix会被设置为true，同时返回该行的开始部分（等于缓存大小的部分）。该行剩余的部分就会在下次调用的时候返回。当下次调用返回该行剩余部分时，isPrefix将会是false。跟ReadSlice一样，返回的line只是buffer的引用，在下次执行IO操作时，line会无效。  

````go
	reader := bufio.NewReader(strings.NewReader("hello \n world"))
	line1, _, _ := reader.ReadLine()
	fmt.Printf("the line1:%s\n", line1)

	line2, _, _ := reader.ReadLine()
	fmt.Printf("the line2:%s\n", line2)
````

##### Peek

`Peek`只是查看下`Reader`有没有读取的n个字节。相比于`ReadSlice`，是并发安全的。因为`ReadSlice`返回的[]byte只是buffer中的引用，在下次IO操作后会无效。  

````go
func main() {
	reader := bufio.NewReaderSize(strings.NewReader("hello world"), 12)
	go Peek(reader)
	go reader.ReadBytes('d')
	time.Sleep(1e8)
}

func Peek(reader *bufio.Reader) {
	line, _ := reader.Peek(5)
	fmt.Printf("%s\n", line)
	time.Sleep(1)
	fmt.Printf("%s\n", line)
}
````

#### Scanner

`bufio.Reader`结构体中所有读取数据的方法，都包含了`delim`分隔符，这个用起来很不方便，所以`Google`对此在go1.1版本中加入了`bufio.Scanner`结构体，用于读取数据。

````go
type Scanner struct {
    // 内含隐藏或非导出字段
}
````

Scanner类型提供了方便的读取数据的接口，如从换行符分隔的文本里读取每一行。  

`Scanner.Scan`方法默认是以换行符`\n`，作为分隔符。如果你想指定分隔符，`Go`语言提供了四种方法，`ScanBytes`(返回单个字节作为一个 `token`), `ScanLines`(返回一行文本), `ScanRunes`(返回单个 `UTF-8` 编码的 `rune` 作为一个 `token`)和`ScanWords`(返回通过“空格”分词的单词)。   

扫描会在抵达输入流结尾、遇到的第一个`I/O`错误、`token`过大不能保存进缓冲时，不可恢复的停止。当扫描停止后，当前读取位置可能会远在最后一个获得的`token`后面。需要更多对错误管理的控制或`token`很大，或必须从`reader`连续扫描的程序，应使用`bufio.Reader`代替。  

````go
	input := "hello world"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
````

##### Give me more data

缓冲区的默认 `size` 是 4096。如果我们指定了最小的缓存区的大小，当在读取的过程中，如果指定的最小缓冲区的大小不足以放置读取的内容，就会发生扩容，原则是新的长度是之前的两倍。
