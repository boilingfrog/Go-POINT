<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [buffer](#buffer)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [例子](#%E4%BE%8B%E5%AD%90)
  - [了解下bytes.buffer](#%E4%BA%86%E8%A7%A3%E4%B8%8Bbytesbuffer)
    - [如何创建bytes.buffer](#%E5%A6%82%E4%BD%95%E5%88%9B%E5%BB%BAbytesbuffer)
    - [bytes.buffer的数据写入](#bytesbuffer%E7%9A%84%E6%95%B0%E6%8D%AE%E5%86%99%E5%85%A5)
      - [写入string](#%E5%86%99%E5%85%A5string)
      - [写入[]byte](#%E5%86%99%E5%85%A5byte)
      - [写入byte](#%E5%86%99%E5%85%A5byte)
      - [写入rune](#%E5%86%99%E5%85%A5rune)
    - [数据写出](#%E6%95%B0%E6%8D%AE%E5%86%99%E5%87%BA)
      - [写出数据到io.Writer](#%E5%86%99%E5%87%BA%E6%95%B0%E6%8D%AE%E5%88%B0iowriter)
      - [Read](#read)
      - [ReadByte](#readbyte)
      - [ReadRune](#readrune)
      - [ReadBytes](#readbytes)
      - [ReadString](#readstring)
      - [Next](#next)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

##### 写入string

````go
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("小花猫")
	fmt.Println(buf.String())
````

##### 写入[]byte

````go
    buf := bytes.NewBuffer([]byte{})
	s := []byte("小黑猫")
	buf.Write(s)
	fmt.Println(buf.String())
````

##### 写入byte

````go
	var b byte = '?'
	buf.WriteByte(b)

	fmt.Println(buf.String())
````

##### 写入rune

````go
	var r rune = '小'
	buf.WriteRune(r)
	fmt.Println(buf.String())
````

#### 数据写出

##### 写出数据到io.Writer

````go
	file, _ := os.Open("text.txt")
	buf := bytes.NewBufferString("hello")
	buf.WriteTo(file) //hello写到text.txt文件中了
````

`os.File`就是实现`io.Writer`

##### Read

````go
	bufRead := bytes.NewBufferString("hello")
	fmt.Println(bufRead.String())
	var sRead = make([]byte, 3)   // 定义读出的[]byte为3，表示一次可读出3个byte
	bufRead.Read(sRead)           // 读出
	fmt.Println(bufRead.String()) // 打印结果为lo,因为前三个被读出了
	fmt.Println(string(sRead))    // 打印结果为hel,读取的是hello的前三个字母

	bufRead.Read(sRead)           // 接着读，但是bufRead之剩下lo，所以只有lo被读出了
	fmt.Println(bufRead.String()) // 打印结果为空
	fmt.Println(string(sRead))    // 大一结果lol，前两位的lo表示的本次的读出，因为bufRead只有两位，后面的l还是上次的读出结果
````

##### ReadByte

````go
    buf := bytes.NewBufferString("hello")
    fmt.Println(buf.String()) //buf.String()方法是吧buf里的内容转成string，>以便于打印
    b, _ := buf.ReadByte()    //读取第一个byte，赋值给b
    fmt.Println(buf.String()) //打印 ello，缓冲器头部第一个h被拿掉
    fmt.Println(string(b))    //打印 h
````

##### ReadRune

````go
    buf := bytes.NewBufferString("好hello")
    fmt.Println(buf.String()) //buf.String()方法是吧buf里的内容转成string，>以便于打印
    b, n, _ := buf.ReadRune() //读取第一个rune，赋值给b
    fmt.Println(buf.String()) //打印 hello
    fmt.Println(string(b))    //打印中文字： 好，缓冲器头部第一个“好”被拿掉
    fmt.Println(n)            //打印3，“好”作为utf8储存占3个byte
    b, n, _ = buf.ReadRune()  //再读取第一个rune，赋值给b
    fmt.Println(buf.String()) //打印 ello
    fmt.Println(string(b))    //打印h，缓冲器头部第一个h被拿掉
    fmt.Println(n)            //打印 1，“h”作为utf8储存占1个byte
````

##### ReadBytes

`ReadBytes`和`ReadByte`是有区别的。`ReadBytes`需要一个分隔符来对`buffer`进行分割读取。

````
    var d byte = 'e' //分隔符为e
	buf := bytes.NewBufferString("hello")
	fmt.Println(buf.String()) //buf.String()方法是吧buf里的内容转成string，以便于打印
	b, _ := buf.ReadBytes(d)  //读到分隔符，并返回给b
	fmt.Println(buf.String()) //打印 llo，缓冲器被取走一些数据
	fmt.Println(string(b))    //打印 he，找到e了，将缓冲器从头开始，到e的内容都返回给b
````

##### ReadString

`ReadString`和`ReadBytes`一样，也是需要一个分隔符进行，`buffer`

````go
	var d byte = 'e' //分隔符为e
	buf := bytes.NewBufferString("hello")
	fmt.Println(buf.String()) //buf.String()方法是吧buf里的内容转成string，以便于打印
	b, _ := buf.ReadString(d) //读到分隔符，并返回给b
	fmt.Println(buf.String()) //打印 llo，缓冲器被取走一些数据
	fmt.Println(b)            //打印 he，找到e了，将缓冲器从头开始，到e的内容都返回给b
````

##### Next

使用`Next`可依次读出固定长度的内容

````
	buf := bytes.NewBufferString("hello")
	fmt.Println(buf.String())
	b := buf.Next(2)          // 重头开始，取2个
	fmt.Println(buf.String()) // 变小了
	fmt.Println(string(b))    // 打印he
````



### 参考

【go语言的bytes.buffer】https://my.oschina.net/u/943306/blog/127981  


 

