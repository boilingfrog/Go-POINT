## buffer

### 前言

最近操作文件，进行优化使用到了`buffer`。好像也不太了解这个，那么就梳理下，`buffer`的使用。

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

### 了解下bytes.buffer

`bytes.buffer`是一个缓冲`byte`类型的缓冲器，这个缓冲器里存放着都是`byte`。

#### 如何创建bytes.buffer

放几种创建的方式

````go
	buf1 := bytes.NewBufferString("hello")
	fmt.Println(buf1)
	buf2 := bytes.NewBuffer([]byte("hello"))
	fmt.Println(buf2)
	buf3 := bytes.NewBuffer([]byte{byte('h'), byte('e'), byte('l'), byte('l'), byte('o')})
	fmt.Println(buf3)
    // 以上三者等效

	buf4 := bytes.NewBufferString("")
	fmt.Println(buf4)
	buf5 := bytes.NewBuffer([]byte{})
	fmt.Println(buf5)
    // 以上两者等效
````

查看源码可知

````go
func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }

func NewBufferString(s string) *Buffer {
	return &Buffer{buf: []byte(s)}
}
````
`NewBufferString`也是将参数转成 `[]byte()`。然后，初始化`Buffer`。

#### bytes.buffer的数据写入







### 参考

【go语言的bytes.buffer】https://my.oschina.net/u/943306/blog/127981  


 

