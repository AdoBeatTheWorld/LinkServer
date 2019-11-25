package LinkServer

import (
	"github.com/funny/link"
	"io"
	"reflect"
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
	ReqId    uint32
	RealSize uint16
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
