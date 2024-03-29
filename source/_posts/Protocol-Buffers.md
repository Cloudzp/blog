---
title: Protocol Buffers
categories:
  - 后端
tags:
  - Protocol Buffers
comments: true
abbrlink: 15539
date: 2019-04-18 21:50:18
img:
---

## 一、 Protocol Buffers 是什么？protoc是什么？ 
 - Protocol Buffers (ProtocolBuffer/ protobuf )是Google公司开发的一种数据描述语言，类似于XML能够将结构化数据序列化，可用于数据存储、通信协议等方面
 - protoc是Protocol Buffers的一个工具，用来支持将Protocol Buffers定义的文件转换为各种语言的客户端代码。
 
## 二、 为什么要有Protocol Buffers？
1、通过它，可以定义我们的数据的结构，并生成基于各种语言的代码。这些你定义的数据流可以轻松地在传递并不破坏我们原有的程序。并且也可以更新这些数据而现有的程序也不会受到任何的影响。
2、而且同XML相比，Protocol buffers在序列化结构化数据方面有许多优点：
- （1）更简单
- （2）数据描述文件只需原来的1/10至1/3
- （3）解析速度是原来的20倍至100倍
- （4）减少了二义性
- （5）生成了更容易在编程中使用的数据访问类
- （6）支持多种编程语言
> 参考：
// 阿里技术文档
1. https://yq.aliyun.com/ziliao/580483 
2. https://www.ibm.com/developerworks/cn/linux/l-cn-gpb/

### 三、如何安装使用？
1.protor的安装 

>参考： https://blog.csdn.net/JustinSeraph/article/details/70171331

2.使用
2.1 使用命令生成proto文件的go语言代码，然后根据go语言代码去做开发即可；
```
protoc --go_out=./go/ ./proto/helloworld.proto
```

[点击查看示例代码](src/main.go)
```proto文件
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

package example

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FOO int32

const (
	FOO_X FOO = 17
)

var FOO_name = map[int32]string{
	17: "X",
}
var FOO_value = map[string]int32{
	"X": 17,
}

func (x FOO) Enum() *FOO {
	p := new(FOO)
	*p = x
	return p
}
func (x FOO) String() string {
	return proto.EnumName(FOO_name, int32(x))
}
func (x *FOO) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FOO_value, data, "FOO")
	if err != nil {
		return err
	}
	*x = FOO(value)
	return nil
}
func (FOO) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_test_c26727442f671233, []int{0}
}

type Test struct {
	Label                *string  `protobuf:"bytes,1,req,name=label" json:"label,omitempty"`
	Type                 *int32   `protobuf:"varint,2,opt,name=type,def=77" json:"type,omitempty"`
	Reps                 []int64  `protobuf:"varint,3,rep,name=reps" json:"reps,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Test) Reset()         { *m = Test{} }
func (m *Test) String() string { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()    {}
func (*Test) Descriptor() ([]byte, []int) {
	return fileDescriptor_test_c26727442f671233, []int{0}
}
func (m *Test) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Test.Unmarshal(m, b)
}
func (m *Test) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Test.Marshal(b, m, deterministic)
}
func (dst *Test) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Test.Merge(dst, src)
}
func (m *Test) XXX_Size() int {
	return xxx_messageInfo_Test.Size(m)
}
func (m *Test) XXX_DiscardUnknown() {
	xxx_messageInfo_Test.DiscardUnknown(m)
}

var xxx_messageInfo_Test proto.InternalMessageInfo

const Default_Test_Type int32 = 77

func (m *Test) GetLabel() string {
	if m != nil && m.Label != nil {
		return *m.Label
	}
	return ""
}

func (m *Test) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Default_Test_Type
}

func (m *Test) GetReps() []int64 {
	if m != nil {
		return m.Reps
	}
	return nil
}

func init() {
	proto.RegisterType((*Test)(nil), "example.Test")
	proto.RegisterEnum("example.FOO", FOO_name, FOO_value)
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_test_c26727442f671233) }

var fileDescriptor_test_c26727442f671233 = []byte{
	// 122 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f, 0xad, 0x48, 0xcc, 0x2d, 0xc8, 0x49, 0x55,
	0xf2, 0xe0, 0x62, 0x09, 0x49, 0x2d, 0x2e, 0x11, 0x12, 0xe1, 0x62, 0xcd, 0x49, 0x4c, 0x4a, 0xcd,
	0x91, 0x60, 0x54, 0x60, 0xd2, 0xe0, 0x0c, 0x82, 0x70, 0x84, 0xc4, 0xb8, 0x58, 0x4a, 0x2a, 0x0b,
	0x52, 0x25, 0x98, 0x14, 0x18, 0x35, 0x58, 0xad, 0x98, 0xcc, 0xcd, 0x83, 0xc0, 0x7c, 0x21, 0x21,
	0x2e, 0x96, 0xa2, 0xd4, 0x82, 0x62, 0x09, 0x66, 0x05, 0x66, 0x0d, 0xe6, 0x20, 0x30, 0x5b, 0x8b,
	0x87, 0x8b, 0xd9, 0xcd, 0xdf, 0x5f, 0x88, 0x95, 0x8b, 0x31, 0x42, 0x40, 0x10, 0x10, 0x00, 0x00,
	0xff, 0xff, 0x71, 0xcd, 0x01, 0xd7, 0x6d, 0x00, 0x00, 0x00,
}

```

```go
package main

import (
	"protocolBuffers/example"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
)

func main(){
	test := &example.Test{
		Label:proto.String("Hello"),
		Type: proto.Int32(17),
		Reps: []int64{1,2,3},
	}

	data, err := proto.Marshal(test)
	if err != nil {
		log.Fatal("Marshal", err)
	}

	newTest := example.Test{}
	err = proto.Unmarshal(data, &newTest)
	if err != nil {
		log.Fatal("Unmarshal:",err)
	}
	// Now test and newTest contain the same data.
	if test.GetLabel() != newTest.GetLabel() {
		log.Fatalf("data mismatch %q != %q", test.GetLabel(), newTest.GetLabel())
	}else {
		fmt.Println(newTest.GetLabel(),newTest.GetType(),newTest.GetReps())
	}
}

```


> 参考：
1. https://www.cnblogs.com/chenyangyao/p/5422044.html (protobuf 介绍)
2. https://github.com/golang/protobuf (go protobuf库)
