package ts

import (
    "encoding/binary"
)

type ProgramAssociationTable struct {
	TableId                uint8
	SectionSyntaxIndicator bool
	SectionLenght          uint16
	TransportStreamId      uint16
	VersionNumber          uint8
	CurrentNextIndicator   bool
	SectionNumber          uint8
	LastSectionNumber      uint8
	Programs               []*Program
	CRC32                  uint32
}

type Program struct {
	Number        uint16
	NetworkPID    uint16
	ProgramMapPID uint16
}

type ProgramMapTable struct {
	TableId                uint8
	SectionSyntaxIndicator bool
	SectionLenght          uint16
	ProgramNumber          uint16
	VersionNumber          uint8
	CurrentNextIndicator   bool
	SectionNumber          uint8
	LastSectionNumber      uint8
	PCRPID                 uint16
	ProgramInfoLength      uint16
	StreamDescriptors      []*StreamDescriptor
	ElementaryStreams      []*ElementaryStream
	CRC32                  uint32
}

type StreamDescriptor struct {
	DescriptorTag    uint8
	DescriptorLength uint8
	Data             []byte
}

type ElementaryStream struct {
	StreamType    uint8
	ElementaryPID uint16
	ESInfoLength  uint16
	Data          []byte
}

func newPogramAssociationTable(payload []byte) *ProgramAssociationTable {
	data := payload[1+payload[0]:]
	a := &ProgramAssociationTable{
		TableId:                data[0],
		SectionSyntaxIndicator: data[1]&0x80>>7 == 1,
		SectionLenght:          uint16(data[1]&0x0F)<<8 | uint16(data[2]),
		TransportStreamId:      uint16(data[3])<<8 | uint16(data[4]),
		VersionNumber:          data[5] & 0x3e >> 1,
		CurrentNextIndicator:   data[5]&0x01 == 1,
		SectionNumber:          data[6],
		LastSectionNumber:      data[7]}

	i := 8
	for i < int(a.SectionLenght)-4 {
        //TODO: 这里的移位,移动7是否有错
		p := &Program{Number: uint16(data[i])<<7 | uint16(data[i+1])}
		if p.Number == 0 {
			p.NetworkPID = uint16(data[i+2]&0x1F)<<8 | uint16(data[i+3])
		} else {
			p.ProgramMapPID = uint16(data[i+2]&0x1F)<<8 | uint16(data[i+3])
		}
		a.Programs = append(a.Programs, p)
		i += 4
	}
	a.CRC32 = uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
	return a
}

func writePogramAssociationTable(pat *ProgramAssociationTable,payload []byte) {

	payload[0] = 0x00
	data := payload[1+payload[0]:]
	data[0] = pat.TableId
    data[1] = 0x00
    if pat.SectionSyntaxIndicator {
        data[1] |= 0x80
    }
    data[1] |= uint8(pat.SectionLenght >> 8) &0x0F
    data[2] =  uint8(pat.SectionLenght & 0x000F)
    data[3] =  uint8((pat.TransportStreamId & 0xFF00) >> 8 )
    data[4] =  uint8(pat.TransportStreamId & 0x00FF)

    data[5] = (pat.VersionNumber << 1) & 0x3e
    if pat.CurrentNextIndicator {
        data[5] |= 0x01
    }

    data[6] = pat.SectionNumber
    data[7] = pat.LastSectionNumber

    i := 8
    for _, p := range pat.Programs {
        //TODO: 这里的移位,移动7是否有错
        data[i] = uint8(p.Number >> 7)
        data[i+1] = uint8(p.Number & 0x00FF)
        if 0 == p.Number {
            data[i+2] = uint8( (p.NetworkPID & 0x1F00) >> 8 )
            data[i+3] = uint8(p.NetworkPID & 0x00FF)
        } else {
            data[i+2] = uint8( (p.ProgramMapPID & 0x1F00) >> 8 )
            data[i+3] = uint8(p.ProgramMapPID & 0x00FF)
        }
        i += 4
    }

    binary.BigEndian.PutUint32(data[i:i+4],pat.CRC32)

}

func newProgramMapTable(payload []byte) *ProgramMapTable {
	data := payload[1+payload[0]:]
	a := &ProgramMapTable{
		TableId:                data[0],
		SectionSyntaxIndicator: data[1]&0x80>>7 == 1,
		SectionLenght:          uint16(data[1]&0x0F)<<8 | uint16(data[2]),
		ProgramNumber:          uint16(data[3])<<8 | uint16(data[4]),
		VersionNumber:          data[5] & 0x3e >> 1,
		CurrentNextIndicator:   data[5]&0x01 == 1,
		SectionNumber:          data[6],
		LastSectionNumber:      data[7],
		PCRPID:                 uint16(data[8]&0x1F)<<8 | uint16(data[9]),
		ProgramInfoLength:      uint16(data[10]&0x0F)<<8 | uint16(data[11])}

	a.StreamDescriptors = newStreamDescriptors(data[12 : int(a.ProgramInfoLength)+12])
	a.ElementaryStreams = newElementaryStreams(data[int(a.ProgramInfoLength)+12 : int(a.SectionLenght)-1])
	a.CRC32 = uint32(data[int(a.SectionLenght)-1])<<24 | uint32(data[int(a.SectionLenght)])<<16 |
		uint32(data[int(a.SectionLenght)+1])<<8 | uint32(data[int(a.SectionLenght)+2])
	return a
}

func writeProgramMapTable(pmt *ProgramMapTable,payload []byte)  {
    payload[0] = 0x00
    data := payload[1+payload[0]:]

    data[0] = pmt.TableId

    data[1] = 0x00
    if pmt.SectionSyntaxIndicator {
        data[1] |= 0x80
    }

	data[1] |= uint8(pmt.SectionLenght >> 8) & 0x0F
	data[2] = uint8(pmt.SectionLenght & 0x00FF)

    data[3] = uint8(pmt.ProgramNumber >> 8)
	data[4] = uint8(pmt.ProgramNumber & 0x00FF)

	data[5] = (pmt.VersionNumber << 1 ) & 0x3e
	if pmt.CurrentNextIndicator {
		data[5] |= 0x01
	}

	data[6] = pmt.SectionNumber
	data[7] = pmt.LastSectionNumber

	data[8] = uint8(pmt.PCRPID >> 8) & 0x1F
	data[9] = uint8(pmt.PCRPID & 0x00ff)

	data[10] = uint8(pmt.ProgramInfoLength >> 8) & 0x0F
	data[11] = uint8(pmt.ProgramInfoLength & 0x00ff)

	writeStreamDescriptors(pmt.StreamDescriptors,data[12:int(pmt.ProgramInfoLength)+12])
	writeElementaryStreams(pmt.ElementaryStreams,data[int(pmt.ProgramInfoLength)+12 : int(pmt.SectionLenght)-1])

    binary.BigEndian.PutUint32(data[int(pmt.SectionLenght-1):int(pmt.SectionLenght+2)],pmt.CRC32)

}

func (p ProgramMapTable) HasElementaryStream(pid uint16) bool {
	for _, e := range p.ElementaryStreams {
		if e.ElementaryPID == pid {
			return true
		}
	}
	return false
}

func newStreamDescriptors(data []byte) []*StreamDescriptor {
	var descriptors []*StreamDescriptor
	for i := 0; i < len(data); {
		d := &StreamDescriptor{
			DescriptorTag:    data[i],
			DescriptorLength: data[i+1]}
		d.Data = data[i+2 : i+2+int(d.DescriptorLength)]
		descriptors = append(descriptors, d)
		i += 2 + int(d.DescriptorLength)
	}
	return descriptors
}

func writeStreamDescriptors(sdList []*StreamDescriptor, data []byte) {
	i := 0
	for _,sd := range sdList {
		data[i] = sd.DescriptorTag
		data[i+1] = sd.DescriptorLength
		copy(data[i+2:i+2+int(sd.DescriptorLength)],sd.Data)
		i += 2+ int(sd.DescriptorLength)
	}

}

func newElementaryStreams(data []byte) []*ElementaryStream {
	var streams []*ElementaryStream
	for i := 0; i < len(data); {
		s := &ElementaryStream{
			StreamType:    data[i],
			ElementaryPID: uint16(data[i+1]&0x1F)<<8 | uint16(data[i+2]),
			ESInfoLength:  uint16(data[i+3]&0x0F)<<8 | uint16(data[i+4])}
		s.Data = data[i+5:]
		i += 5 + int(s.ESInfoLength)
		streams = append(streams, s)
	}
	return streams
}

func writeElementaryStreams(esList []*ElementaryStream,data []byte)  {
	i := 0
	for _,es := range esList {
		data[i] = es.StreamType

		data[i+1] = uint8(es.ElementaryPID >> 8) & 0x1F
		data[i+2] = uint8(es.ElementaryPID & 0x00ff)

		data[i+3] = uint8(es.ESInfoLength >> 8) & 0x0F
		data[i+4] = uint8(es.ESInfoLength & 0x00FF)

		copy(data[i+5:i+5+int(es.ESInfoLength)],es.Data)
		i += 5 + int(es.ESInfoLength)
	}

}
