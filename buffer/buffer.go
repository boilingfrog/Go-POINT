package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

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
	file, err := os.Open("./buffer/test.txt") //test.txt的内容是“world”
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Println(file.Sync())
	buf := bytes.NewBufferString("hello ")
	buf.ReadFrom(file)        //将text.txt内容追加到缓冲器的尾部
	fmt.Println(buf.String()) //打印“hello world”
}
