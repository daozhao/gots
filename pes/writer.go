package pes

import (
	"github.com/damienlevin/gots/ts"
)

type Writer struct {
	Writer *ts.Writer

	//Packet      *Packet
    PrevPacket map[uint16]*Packet
	OnNewPacket func(*Packet)
}

func NewWriter(w *ts.Writer, c func(*Packet)) *Writer{
	return &Writer{
		Writer:      w,
		PrevPacket: make(map[uint16]*Packet),
		OnNewPacket: c}
}

func (w *Writer) WritePacket(p *Packet)(data []byte) {
	//var p *Packet
	//var err error

    data = writePacket(p)

	return data
}


//func (w *Writer) WriteAVCRawData(payload []byte,DataAlignmentIndicator bool,pts uint64,dts uint64)(error) {
//
//
//
//	return nil
//}