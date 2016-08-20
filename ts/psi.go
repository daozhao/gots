package ts

import (
    "encoding/binary"
	"bytes"
	//"fmt"
	"hash/crc32"
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

var Crc32Tbl = crc32.Table{
	0x00000000, 0x04c11db7, 0x09823b6e, 0x0d4326d9, 0x130476dc, 0x17c56b6b,
	0x1a864db2, 0x1e475005, 0x2608edb8, 0x22c9f00f, 0x2f8ad6d6, 0x2b4bcb61,
	0x350c9b64, 0x31cd86d3, 0x3c8ea00a, 0x384fbdbd, 0x4c11db70, 0x48d0c6c7,
	0x4593e01e, 0x4152fda9, 0x5f15adac, 0x5bd4b01b, 0x569796c2, 0x52568b75,
	0x6a1936c8, 0x6ed82b7f, 0x639b0da6, 0x675a1011, 0x791d4014, 0x7ddc5da3,
	0x709f7b7a, 0x745e66cd, 0x9823b6e0, 0x9ce2ab57, 0x91a18d8e, 0x95609039,
	0x8b27c03c, 0x8fe6dd8b, 0x82a5fb52, 0x8664e6e5, 0xbe2b5b58, 0xbaea46ef,
	0xb7a96036, 0xb3687d81, 0xad2f2d84, 0xa9ee3033, 0xa4ad16ea, 0xa06c0b5d,
	0xd4326d90, 0xd0f37027, 0xddb056fe, 0xd9714b49, 0xc7361b4c, 0xc3f706fb,
	0xceb42022, 0xca753d95, 0xf23a8028, 0xf6fb9d9f, 0xfbb8bb46, 0xff79a6f1,
	0xe13ef6f4, 0xe5ffeb43, 0xe8bccd9a, 0xec7dd02d, 0x34867077, 0x30476dc0,
	0x3d044b19, 0x39c556ae, 0x278206ab, 0x23431b1c, 0x2e003dc5, 0x2ac12072,
	0x128e9dcf, 0x164f8078, 0x1b0ca6a1, 0x1fcdbb16, 0x018aeb13, 0x054bf6a4,
	0x0808d07d, 0x0cc9cdca, 0x7897ab07, 0x7c56b6b0, 0x71159069, 0x75d48dde,
	0x6b93dddb, 0x6f52c06c, 0x6211e6b5, 0x66d0fb02, 0x5e9f46bf, 0x5a5e5b08,
	0x571d7dd1, 0x53dc6066, 0x4d9b3063, 0x495a2dd4, 0x44190b0d, 0x40d816ba,
	0xaca5c697, 0xa864db20, 0xa527fdf9, 0xa1e6e04e, 0xbfa1b04b, 0xbb60adfc,
	0xb6238b25, 0xb2e29692, 0x8aad2b2f, 0x8e6c3698, 0x832f1041, 0x87ee0df6,
	0x99a95df3, 0x9d684044, 0x902b669d, 0x94ea7b2a, 0xe0b41de7, 0xe4750050,
	0xe9362689, 0xedf73b3e, 0xf3b06b3b, 0xf771768c, 0xfa325055, 0xfef34de2,
	0xc6bcf05f, 0xc27dede8, 0xcf3ecb31, 0xcbffd686, 0xd5b88683, 0xd1799b34,
	0xdc3abded, 0xd8fba05a, 0x690ce0ee, 0x6dcdfd59, 0x608edb80, 0x644fc637,
	0x7a089632, 0x7ec98b85, 0x738aad5c, 0x774bb0eb, 0x4f040d56, 0x4bc510e1,
	0x46863638, 0x42472b8f, 0x5c007b8a, 0x58c1663d, 0x558240e4, 0x51435d53,
	0x251d3b9e, 0x21dc2629, 0x2c9f00f0, 0x285e1d47, 0x36194d42, 0x32d850f5,
	0x3f9b762c, 0x3b5a6b9b, 0x0315d626, 0x07d4cb91, 0x0a97ed48, 0x0e56f0ff,
	0x1011a0fa, 0x14d0bd4d, 0x19939b94, 0x1d528623, 0xf12f560e, 0xf5ee4bb9,
	0xf8ad6d60, 0xfc6c70d7, 0xe22b20d2, 0xe6ea3d65, 0xeba91bbc, 0xef68060b,
	0xd727bbb6, 0xd3e6a601, 0xdea580d8, 0xda649d6f, 0xc423cd6a, 0xc0e2d0dd,
	0xcda1f604, 0xc960ebb3, 0xbd3e8d7e, 0xb9ff90c9, 0xb4bcb610, 0xb07daba7,
	0xae3afba2, 0xaafbe615, 0xa7b8c0cc, 0xa379dd7b, 0x9b3660c6, 0x9ff77d71,
	0x92b45ba8, 0x9675461f, 0x8832161a, 0x8cf30bad, 0x81b02d74, 0x857130c3,
	0x5d8a9099, 0x594b8d2e, 0x5408abf7, 0x50c9b640, 0x4e8ee645, 0x4a4ffbf2,
	0x470cdd2b, 0x43cdc09c, 0x7b827d21, 0x7f436096, 0x7200464f, 0x76c15bf8,
	0x68860bfd, 0x6c47164a, 0x61043093, 0x65c52d24, 0x119b4be9, 0x155a565e,
	0x18197087, 0x1cd86d30, 0x029f3d35, 0x065e2082, 0x0b1d065b, 0x0fdc1bec,
	0x3793a651, 0x3352bbe6, 0x3e119d3f, 0x3ad08088, 0x2497d08d, 0x2056cd3a,
	0x2d15ebe3, 0x29d4f654, 0xc5a92679, 0xc1683bce, 0xcc2b1d17, 0xc8ea00a0,
	0xd6ad50a5, 0xd26c4d12, 0xdf2f6bcb, 0xdbee767c, 0xe3a1cbc1, 0xe760d676,
	0xea23f0af, 0xeee2ed18, 0xf0a5bd1d, 0xf464a0aa, 0xf9278673, 0xfde69bc4,
	0x89b8fd09, 0x8d79e0be, 0x803ac667, 0x84fbdbd0, 0x9abc8bd5, 0x9e7d9662,
	0x933eb0bb, 0x97ffad0c, 0xafb010b1, 0xab710d06, 0xa6322bdf, 0xa2f33668,
	0xbcb4666d, 0xb8757bda, 0xb5365d03, 0xb1f740b4,
}

func makeCRC32(data []byte) uint32 {
    var crc uint32 = 0xFFFFFFFF

    for _,b := range data {
        crc = (crc << 8) ^ Crc32Tbl[ (byte(crc >> 24)^b) & 0xFF ]
    }
    return crc

}
var ieeeCrc32Tbl = crc32.Table{
	0x00000000, 0xB71DC104, 0x6E3B8209, 0xD926430D, 0xDC760413, 0x6B6BC517,
	0xB24D861A, 0x0550471E, 0xB8ED0826, 0x0FF0C922, 0xD6D68A2F, 0x61CB4B2B,
	0x649B0C35, 0xD386CD31, 0x0AA08E3C, 0xBDBD4F38, 0x70DB114C, 0xC7C6D048,
	0x1EE09345, 0xA9FD5241, 0xACAD155F, 0x1BB0D45B, 0xC2969756, 0x758B5652,
	0xC836196A, 0x7F2BD86E, 0xA60D9B63, 0x11105A67, 0x14401D79, 0xA35DDC7D,
	0x7A7B9F70, 0xCD665E74, 0xE0B62398, 0x57ABE29C, 0x8E8DA191, 0x39906095,
	0x3CC0278B, 0x8BDDE68F, 0x52FBA582, 0xE5E66486, 0x585B2BBE, 0xEF46EABA,
	0x3660A9B7, 0x817D68B3, 0x842D2FAD, 0x3330EEA9, 0xEA16ADA4, 0x5D0B6CA0,
	0x906D32D4, 0x2770F3D0, 0xFE56B0DD, 0x494B71D9, 0x4C1B36C7, 0xFB06F7C3,
	0x2220B4CE, 0x953D75CA, 0x28803AF2, 0x9F9DFBF6, 0x46BBB8FB, 0xF1A679FF,
	0xF4F63EE1, 0x43EBFFE5, 0x9ACDBCE8, 0x2DD07DEC, 0x77708634, 0xC06D4730,
	0x194B043D, 0xAE56C539, 0xAB068227, 0x1C1B4323, 0xC53D002E, 0x7220C12A,
	0xCF9D8E12, 0x78804F16, 0xA1A60C1B, 0x16BBCD1F, 0x13EB8A01, 0xA4F64B05,
	0x7DD00808, 0xCACDC90C, 0x07AB9778, 0xB0B6567C, 0x69901571, 0xDE8DD475,
	0xDBDD936B, 0x6CC0526F, 0xB5E61162, 0x02FBD066, 0xBF469F5E, 0x085B5E5A,
	0xD17D1D57, 0x6660DC53, 0x63309B4D, 0xD42D5A49, 0x0D0B1944, 0xBA16D840,
	0x97C6A5AC, 0x20DB64A8, 0xF9FD27A5, 0x4EE0E6A1, 0x4BB0A1BF, 0xFCAD60BB,
	0x258B23B6, 0x9296E2B2, 0x2F2BAD8A, 0x98366C8E, 0x41102F83, 0xF60DEE87,
	0xF35DA999, 0x4440689D, 0x9D662B90, 0x2A7BEA94, 0xE71DB4E0, 0x500075E4,
	0x892636E9, 0x3E3BF7ED, 0x3B6BB0F3, 0x8C7671F7, 0x555032FA, 0xE24DF3FE,
	0x5FF0BCC6, 0xE8ED7DC2, 0x31CB3ECF, 0x86D6FFCB, 0x8386B8D5, 0x349B79D1,
	0xEDBD3ADC, 0x5AA0FBD8, 0xEEE00C69, 0x59FDCD6D, 0x80DB8E60, 0x37C64F64,
	0x3296087A, 0x858BC97E, 0x5CAD8A73, 0xEBB04B77, 0x560D044F, 0xE110C54B,
	0x38368646, 0x8F2B4742, 0x8A7B005C, 0x3D66C158, 0xE4408255, 0x535D4351,
	0x9E3B1D25, 0x2926DC21, 0xF0009F2C, 0x471D5E28, 0x424D1936, 0xF550D832,
	0x2C769B3F, 0x9B6B5A3B, 0x26D61503, 0x91CBD407, 0x48ED970A, 0xFFF0560E,
	0xFAA01110, 0x4DBDD014, 0x949B9319, 0x2386521D, 0x0E562FF1, 0xB94BEEF5,
	0x606DADF8, 0xD7706CFC, 0xD2202BE2, 0x653DEAE6, 0xBC1BA9EB, 0x0B0668EF,
	0xB6BB27D7, 0x01A6E6D3, 0xD880A5DE, 0x6F9D64DA, 0x6ACD23C4, 0xDDD0E2C0,
	0x04F6A1CD, 0xB3EB60C9, 0x7E8D3EBD, 0xC990FFB9, 0x10B6BCB4, 0xA7AB7DB0,
	0xA2FB3AAE, 0x15E6FBAA, 0xCCC0B8A7, 0x7BDD79A3, 0xC660369B, 0x717DF79F,
	0xA85BB492, 0x1F467596, 0x1A163288, 0xAD0BF38C, 0x742DB081, 0xC3307185,
	0x99908A5D, 0x2E8D4B59, 0xF7AB0854, 0x40B6C950, 0x45E68E4E, 0xF2FB4F4A,
	0x2BDD0C47, 0x9CC0CD43, 0x217D827B, 0x9660437F, 0x4F460072, 0xF85BC176,
	0xFD0B8668, 0x4A16476C, 0x93300461, 0x242DC565, 0xE94B9B11, 0x5E565A15,
	0x87701918, 0x306DD81C, 0x353D9F02, 0x82205E06, 0x5B061D0B, 0xEC1BDC0F,
	0x51A69337, 0xE6BB5233, 0x3F9D113E, 0x8880D03A, 0x8DD09724, 0x3ACD5620,
	0xE3EB152D, 0x54F6D429, 0x7926A9C5, 0xCE3B68C1, 0x171D2BCC, 0xA000EAC8,
	0xA550ADD6, 0x124D6CD2, 0xCB6B2FDF, 0x7C76EEDB, 0xC1CBA1E3, 0x76D660E7,
	0xAFF023EA, 0x18EDE2EE, 0x1DBDA5F0, 0xAAA064F4, 0x738627F9, 0xC49BE6FD,
	0x09FDB889, 0xBEE0798D, 0x67C63A80, 0xD0DBFB84, 0xD58BBC9A, 0x62967D9E,
	0xBBB03E93, 0x0CADFF97, 0xB110B0AF, 0x060D71AB, 0xDF2B32A6, 0x6836F3A2,
	0x6D66B4BC, 0xDA7B75B8, 0x035D36B5, 0xB440F7B1,
}

func calcCRC32(crc uint32, data []byte) uint32 {
	for _, b := range data {
		crc = ieeeCrc32Tbl[b^byte(crc)] ^ (crc >> 8)
	}
	return crc
}

func MakePogramAssociationTable(pid uint16) *ProgramAssociationTable {

	a := &ProgramAssociationTable{
		TableId:                0,
		SectionSyntaxIndicator: true,
		TransportStreamId:      1,
		VersionNumber:          0,
		CurrentNextIndicator:   true,
		SectionNumber:          0,
		LastSectionNumber:      0}

	p := &Program{Number: 1,
		NetworkPID:0,
	    ProgramMapPID:pid}
	a.Programs = append(a.Programs, p)

	return a

}
func NewPogramAssociationTable(payload []byte) *ProgramAssociationTable {
	data := payload[1+payload[0]:]
	a := &ProgramAssociationTable{
		TableId:                data[0],
		//TODO: 确认是否移位7????
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
	//fmt.Printf("========================CRC:%X    myCRC32:%X \n",makeCRC32(data[0:i]),calcCRC32(0xffffffff,data[0:i]))
	a.CRC32 = uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
	return a
}

func WritePogramAssociationTable(pat *ProgramAssociationTable)(payload []byte) {

	//payload = make([]byte,PacketSize-4)
	payload = bytes.Repeat([]byte{0xFF},PacketSize-4)

	payload[0] = 0x00
	data := payload[1+payload[0]:]
	data[0] = pat.TableId
    data[1] = 0x00
    if pat.SectionSyntaxIndicator {
        data[1] |= 0x80
    }
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

	pat.SectionLenght = uint16(i + 1)

	data[1] |= uint8(pat.SectionLenght >> 8) &0x0F
	data[2] =  uint8(pat.SectionLenght & 0x000F)
    //fmt.Printf("XXXXXXXXXXXXXXXXPAT sectionLenght len:%d\n",pat.SectionLenght)
    //data2 := []byte{0x00, 0xb0, 0x0d, 0x00, 0x01, 0xc1, 0x00,  0x00, 0x00, 0x01, 0xe1, 0x00,}
    //copy(data,data2)
    //var crc32 uint32
    crc32 := makeCRC32(data[0:i])
	//TODO: 需要计算CRC32的值.
    binary.BigEndian.PutUint32(data[i:i+4],crc32)

	return payload
}

func MakeProgramMapTable(pid uint16) *ProgramMapTable {
	a := &ProgramMapTable{
		TableId:                2,
		SectionSyntaxIndicator: true,
		ProgramNumber:          1,
		VersionNumber:          0,
		CurrentNextIndicator:   true,
		SectionNumber:          0,
		LastSectionNumber:      0,
		PCRPID:                 257,}

	h264 := &ElementaryStream{
		StreamType:   0x1B ,
		ElementaryPID: pid+1,
		ESInfoLength:  0}

	aac := &ElementaryStream{
		StreamType:   0x0F ,
		ElementaryPID: pid+2 ,
		ESInfoLength:  0}

	a.ElementaryStreams = append(a.ElementaryStreams,h264)
	a.ElementaryStreams = append(a.ElementaryStreams,aac)

	return a
}

func NewProgramMapTable(payload []byte) *ProgramMapTable {
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

func WriteProgramMapTable(pmt *ProgramMapTable)(payload []byte)  {

	payload = bytes.Repeat([]byte{0xFF},PacketSize-4)

    payload[0] = 0x00
    data := payload[1+payload[0]:]

    data[0] = pmt.TableId

    data[1] = 0x00
    if pmt.SectionSyntaxIndicator {
        data[1] |= 0x80
    }


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

	i := writeStreamDescriptors(pmt.StreamDescriptors,data[12:])

    pmt.ProgramInfoLength = uint16(i)
    data[10] = uint8(pmt.ProgramInfoLength >> 8) & 0x0F
    data[11] = uint8(pmt.ProgramInfoLength & 0x00ff)
    //fmt.Printf("XXXXXXXXXXXXXXXXPMT sectionLenght len:%d\n",pmt.SectionLenght)

	i = writeElementaryStreams(pmt.ElementaryStreams,data[12+int(pmt.ProgramInfoLength):])

    pmt.SectionLenght = 12 + pmt.ProgramInfoLength + uint16(i) +1
    //fmt.Printf("XXXXXXXXXXXXXXXXPMT sectionLenght len:%d\n",pmt.SectionLenght)
    data[1] |= uint8(pmt.SectionLenght >> 8) & 0x0F
    data[2] = uint8(pmt.SectionLenght & 0x00FF)
//TODO: 需要计算CRC32的值.

    crc32 := makeCRC32(data[0:int(pmt.SectionLenght-1)])
    //fmt.Printf("data len:%d,sectionlength:%d\n",len(data),pmt.SectionLenght)
    binary.BigEndian.PutUint32(data[int(pmt.SectionLenght-1):],crc32)

	return payload
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

func writeStreamDescriptors(sdList []*StreamDescriptor, data []byte)(int) {
	i := 0
	for _,sd := range sdList {
		data[i] = sd.DescriptorTag
		data[i+1] = sd.DescriptorLength
		copy(data[i+2:i+2+int(sd.DescriptorLength)],sd.Data)
		i += 2+ int(sd.DescriptorLength)
	}

    return i

}

func newElementaryStreams(data []byte) []*ElementaryStream {
	var streams []*ElementaryStream
	for i := 0; i < len(data); {
		s := &ElementaryStream{
			StreamType:    data[i],
			ElementaryPID: uint16(data[i+1]&0x1F)<<8 | uint16(data[i+2]),
			ESInfoLength:  uint16(data[i+3]&0x0F)<<8 | uint16(data[i+4])}
		s.Data = data[i+5:i+5+int(s.ESInfoLength)]
		i += 5 + int(s.ESInfoLength)
		streams = append(streams, s)
	}
	return streams
}

func writeElementaryStreams(esList []*ElementaryStream,data []byte)(int)  {
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

    return i

}
