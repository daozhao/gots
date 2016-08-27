package pes

import (
	"github.com/damienlevin/gots/ts"
	//"github.com/quarnster/util/container"
	"fmt"
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

func (w *Writer) WritePacketToTS(p *Packet) {

	data := writePacket(p)

    l := len(data)

    pos := 0

    var pid uint16
    pid = 0
    if STREAMID_AAC_CODE == p.StreamID {
       pid = 258
    } else if STREAMID_AVC_H264_CODE == p.StreamID {
        pid = 257
    }

    for pos < l {
        //这里的判断应该还有点问题,可能临界数值会有问题.
        if ( pos+ts.PacketSize-4 <= l ){
            tsPacket := ts.MakePacket(pid,0==pos,true,w.Writer.GetContinuitycounter(),nil,0)
            tsPacket.Payload = data[pos:pos+ts.PacketSize-4]
            pos += ts.PacketSize-4
            w.Writer.WritePacket(tsPacket)
        } else {
            tsPacket := ts.MakePacket(pid,0==pos,true,w.Writer.GetContinuitycounter(),nil,uint8(ts.PacketSize-4-1-(l-pos)))
	        //fmt.Println("end of playload len:",l-pos," adap len",ts.PacketSize-4-1-(l-pos))
            tsPacket.Payload = data[pos:l]
            pos = l
            w.Writer.WritePacket(tsPacket)
        }
    }
}


func (w *Writer) WriteAVCRawData(payload []byte,DataAlignmentIndicator bool,pts uint64,dts uint64)(*Packet,error) {

    return  w.WriteRawData(payload,STREAMID_AVC_H264_CODE,DataAlignmentIndicator,pts,dts)
}

func (w *Writer) WriteAACRawData(payload []byte,DataAlignmentIndicator bool,pts uint64)(*Packet,error) {

    return  w.WriteRawData(payload,STREAMID_AAC_CODE,DataAlignmentIndicator,pts,0)
}

func (w *Writer) WriteRawData(payload []byte,streamID uint8,DataAlignmentIndicator bool,pts uint64,dts uint64)(*Packet,error) {

	containsDTS := false
	var headerLength uint8
	headerLength = 5

	if (dts > 0 || ( 0 == pts && 0 == dts) ) {
		containsDTS = true
		headerLength = 10
	}

	hd := &Header{
		ScramblingControl         : 8,
		Priority                  : false,
		DataAlignmentIndicator    : DataAlignmentIndicator,
		Copyright                 : false,
		Original                  : false,
		ContainsPTS               : true,
		ContainsDTS               : containsDTS,
		ContainsESCR              : false,
		ContainsESRate            : false,
		ContainsDSMTrickMode      : false,
		ContainsAdditionalCopyInfo: false,
		ContainsCRC               : false,
		ContainsExtension         : false,
		HeaderLength              : headerLength,

		PTS                       : pts,
		DTS                       : dts }

    var packetLength uint16
    packetLength = 0

    packetLength = uint16(len(payload) + 3 + int(hd.HeaderLength))
    if !(len(payload) < 0xFFFF && STREAMID_AAC_CODE == streamID) {
        packetLength = 0
    }

	pk := &Packet {
		CodePrefix   : 1,
		StreamID     : streamID,
		PacketLength : packetLength,
		Header       : hd,
		Payload      : payload,
	}

	return pk, nil
}