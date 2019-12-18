package linkserver

import (
	"bytes"
	"encoding/binary"
)

type Header struct {
	Len      uint16
	Crc      uint16 //hash/crc32
	Ver      uint16
	Sign     uint16
	MainId   uint8
	SubId    uint8
	EncType  uint8
	Reserve  uint8
	ReqId    uint32 //sequence id
	RealSize uint16
}

func (h *Header) ToBuf() *bytes.Buffer {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, h.Len)
	binary.Write(buf, binary.LittleEndian, h.Crc)
	binary.Write(buf, binary.LittleEndian, h.Ver)
	binary.Write(buf, binary.LittleEndian, h.Sign)
	binary.Write(buf, binary.LittleEndian, h.MainId)
	binary.Write(buf, binary.LittleEndian, h.SubId)
	binary.Write(buf, binary.LittleEndian, h.EncType)
	binary.Write(buf, binary.LittleEndian, h.Reserve)
	binary.Write(buf, binary.LittleEndian, h.ReqId)
	binary.Write(buf, binary.LittleEndian, h.RealSize)
	return buf
}

func (h *Header) FromBuf(buf *bytes.Buffer) {
	binary.Read(buf, binary.LittleEndian, h.Len)
	binary.Read(buf, binary.LittleEndian, h.Crc)
	binary.Read(buf, binary.LittleEndian, h.Ver)
	binary.Read(buf, binary.LittleEndian, h.Sign)
	binary.Read(buf, binary.LittleEndian, h.MainId)
	binary.Read(buf, binary.LittleEndian, h.SubId)
	binary.Read(buf, binary.LittleEndian, h.EncType)
	binary.Read(buf, binary.LittleEndian, h.Reserve)
	binary.Read(buf, binary.LittleEndian, h.ReqId)
	binary.Read(buf, binary.LittleEndian, h.RealSize)
}

type InternalPreHeader struct {
	Len      uint16
	BKicking int16
	NOK      int32
	UserId   int64
	Ip       uint32
	Session  string
	AesKey   string
	CheckSum uint16 //crypto/md5
}

func (ih *InternalPreHeader) ToBuf() *bytes.Buffer {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ih.Len)
	binary.Write(buf, binary.LittleEndian, ih.BKicking)
	binary.Write(buf, binary.LittleEndian, ih.NOK)
	binary.Write(buf, binary.LittleEndian, ih.UserId)
	binary.Write(buf, binary.LittleEndian, ih.Ip)
	binary.Write(buf, binary.LittleEndian, ih.Session)
	binary.Write(buf, binary.LittleEndian, ih.AesKey)
	binary.Write(buf, binary.LittleEndian, ih.CheckSum)
	return buf
}

func (ih *InternalPreHeader) FromBuf(buf *bytes.Buffer) {
	binary.Read(buf, binary.LittleEndian, ih.Len)
	binary.Read(buf, binary.LittleEndian, ih.BKicking)
	binary.Read(buf, binary.LittleEndian, ih.NOK)
	binary.Read(buf, binary.LittleEndian, ih.UserId)
	binary.Read(buf, binary.LittleEndian, ih.Ip)
	binary.Read(buf, binary.LittleEndian, ih.Session)
	binary.Read(buf, binary.LittleEndian, ih.AesKey)
	binary.Read(buf, binary.LittleEndian, ih.CheckSum)
}

/*
type ProtoProtocol struct {
	types map[string]reflect.Type
	names map[reflect.Type]string
}

func (p *ProtoProtocol) Register(t interface{})  {
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	name := rt.PkgPath()+"/"+rt.Name()
	p.types[name] = rt
	p.names[rt] = name
}

func (p *ProtoProtocol) RegisterName(name string, t interface{})  {
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	p.types[name] = rt
	p.names[rt] = name
}

func (p *ProtoProtocol) NewCodec(rw io.ReadWriter) (link.Codec, error) {
	return nil, nil
}

func Proto() *ProtoProtocol {
	return &ProtoProtocol{
		types:make(map[string]reflect.Type),
		names:make(map[reflect.Type]string),
	}
}

type protoIn struct {

}

type protoCodec struct {
	p *ProtoProtocol
	closer io.Closer
}

func (p *protoCodec) Receive() (interface{}, error) {
	var in protoIn
}
*/
