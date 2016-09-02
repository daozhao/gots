package pes

import (
"fmt"
"github.com/damienlevin/gots/ts"
)

const (
	StartCodePrefix uint32 = 0x000001
	MarkerBits      uint16 = 0x02
	STREAMID_AVC_H264_CODE uint8 = 0xE0
	STREAMID_AAC_CODE uint8 = 0xC0
)

type Packet struct {
	CodePrefix   uint32
	StreamID     uint8
	PacketLength uint16
	Header       *Header
	Payload      []byte
	Body         []byte
    Ts           []*ts.Packet
	HasKeyFrame   bool
}

func newPacket(data []byte) (*Packet, error) {
	//TODO: 怀疑移位有错误
	cp := uint32(data[0])<<8 | uint32(data[1])<<4 | uint32(data[2])
	if cp != StartCodePrefix {
		return nil, fmt.Errorf("Invalid PES packet, must start with proper code prefix, got %x expect %x", cp, StartCodePrefix)
	}
	p := &Packet{
		CodePrefix:   cp,
		StreamID:     data[3],
		PacketLength: uint16(data[4])<<8 | uint16(data[5])}

	p.Payload = []byte{}
	if hasHeader(p.StreamID) {
		p.Header = newHeader(data)
		p.Payload = append(p.Payload,data[9 + p.Header.HeaderLength:]...)
	} else {
		p.Payload = append(p.Payload,data[6:]...)
	}
	return p, nil
}

//func MakePESPacket(streamId uint8, pts,dts uint64,rawData []byte)(*Packet){
//    p := &Packet{
//        CodePrefix:   0x01,
//        StreamID:     streamId,
//        PacketLength: 0,
//        Payload:      rawData,
//        Ts:           make([]*ts.Packet,0),
//    }
//
//    p.Header = &Header{
//        ScramblingControl         : 0x00,
//        Priority                  : false,
//        DataAlignmentIndicator    : false,
//        Copyright                 : false,
//        Original                  : false,
//        ContainsPTS               : true,
//        ContainsDTS               : false,
//        ContainsESCR              : false,
//        ContainsESRate            : false,
//        ContainsDSMTrickMode      : false,
//        ContainsAdditionalCopyInfo: false,
//        ContainsCRC               : false,
//        ContainsExtension         : false,
//
//        HeaderLength              : 5,
//        PTS                       : pts,
//        DTS                       : dts,
//    }
//
//    if dts > 0 {
//        p.Header.ContainsDTS  = true
//        p.Header.HeaderLength = 10
//    }
//}

func writePacket(packet *Packet)(data []byte){

	if hasHeader(packet.StreamID) {
		data = make([]byte, len(packet.Payload)+ int(packet.Header.HeaderLength) + 9)
	} else {
		data = make([]byte, len(packet.Payload) + 6)
	}
//TODO: 怀疑移位有错误
	data[0] = uint8((packet.CodePrefix & 0x00FF0000) >> 16)
	data[1] = uint8((packet.CodePrefix & 0x0000FF00) >> 8)
	data[2] = uint8((packet.CodePrefix & 0x000000FF) )

    data[3] = packet.StreamID

    data[4] = uint8((packet.PacketLength & 0xFF00) >> 8)
    data[5] = uint8(packet.PacketLength & 0x00FF)

    if hasHeader(packet.StreamID) {
        writeHeader(packet.Header,data[6:])
	    copy(data[9 + packet.Header.HeaderLength:],packet.Payload)
    } else {
	    copy(data[6:],packet.Payload)
    }
	return  data
}

func hasHeader(id uint8) bool {
	return id != 0xBC && id != 0xBE && id != 0xBF &&
		id != 0xF0 && id != 0xF1 && id != 0xFF &&
		id != 0xF2 && id != 0xF8
}

func IsAudio(id uint8) bool {
	return 0xC0 <= id && id <= 0xDF
}

func IsVideo(id uint8) bool {
	return 0xE0 <= id && id <= 0xEF
}
