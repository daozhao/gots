package pes

type Header struct {
	ScramblingControl          uint8
	Priority                   bool
	DataAlignmentIndicator     bool
	Copyright                  bool
	Original                   bool
	ContainsPTS                bool
	ContainsDTS                bool
	ContainsESCR               bool
	ContainsESRate             bool
	ContainsDSMTrickMode       bool
	ContainsAdditionalCopyInfo bool
	ContainsCRC                bool
	ContainsExtension          bool
	HeaderLength               uint8

	PTS uint64
	DTS uint64
	//TODO Parse me
	ESCR                 uint64
	ESRate               uint32
	TrickModeControl     uint8
	FieldId              uint8
	IntraSliceRefresh    bool
	FrequencyTruncation  uint8
	RepCntrl             uint8
	AdditionalCopyInfo   uint8
	PreviousPESPacketCRC uint16
}

func newHeader(data []byte) *Header {
	h := &Header{
		//TODO: 这里怀疑是错误的.
		//ScramblingControl:          data[6] & 0x30 >> 4,
		ScramblingControl:          data[6] & 0xF0 >> 4,
		Priority:                   data[6]&0x08>>3 == 1,
		DataAlignmentIndicator:     data[6]&0x04>>2 == 1,
		Copyright:                  data[6]&0x02>>1 == 1,
		Original:                   data[6]&0x01 == 1,

		ContainsPTS:                data[7]&0x80>>7 == 1,
		ContainsDTS:                data[7]&0x40>>6 == 1,
        //TODO: 这里怀疑错误
		ContainsESCR:               data[7]&0x10>>4 == 1,
		ContainsESRate:             data[7]&0x20>>4 == 1,
        //ContainsESRate:             data[7]&0x10>>4 == 1,
		ContainsDSMTrickMode:       data[7]&0x08>>3 == 1,
		ContainsAdditionalCopyInfo: data[7]&0x04>>2 == 1,
		ContainsCRC:                data[7]&0x02>>1 == 1,
		ContainsExtension:          data[7]&0x01 == 1,

		HeaderLength:               data[8]}

	if h.ContainsPTS && !h.ContainsDTS {
		h.PTS = uint64(data[9]&0x0E>>1) << 30
		h.PTS = h.PTS | (uint64(data[10])<<8|uint64(data[11]))>>1<<15
		h.PTS = h.PTS | (uint64(data[12])<<8|uint64(data[13]))>>1
	}

	if h.ContainsPTS && h.ContainsDTS {
		h.PTS = uint64(data[9]&0x0E>>1) << 30
		h.PTS = h.PTS | (uint64(data[10])<<8|uint64(data[11]))>>1<<15
		h.PTS = h.PTS | (uint64(data[12])<<8|uint64(data[13]))>>1

		h.DTS = uint64(data[14]&0x0E>>1) << 30
		h.DTS = h.DTS | (uint64(data[15])<<8|uint64(data[16]))>>1<<15
		h.DTS = h.DTS | (uint64(data[17])<<8|uint64(data[18]))>>1
	}

	return h
}

func writeHeader(h *Header,data []byte){

	data[0] = 0x00;
	data[1] = 0x00;

	data[0] = (h.ScramblingControl << 4) & 0xF0
	if h.Priority {
		data[0] |= 0x08
	}
	if h.DataAlignmentIndicator {
		data[0] |= 0x04
	}
	if h.Copyright {
		data[0] |= 0x02
	}
	if h.Original {
		data[0] |= 0x01
	}

    if h.ContainsPTS {
        data[1] |= 0x80
    }
    if h.ContainsDTS {
        data[1] |= 0x40
    }
    if h.ContainsESCR {
        data[1] |= 0x10
    }
    if h.ContainsESRate {
        data[1] |= 0x20
    }
    if h.ContainsDSMTrickMode {
        data[1] |= 0x08
    }
    if h.ContainsAdditionalCopyInfo {
        data[1] |= 0x04
    }
    if h.ContainsCRC {
        data[1] |= 0x02
    }
    if h.ContainsExtension {
        data[1] |= 0x01
    }
//TODO: 这个太怪了,检查协议,看看是否有问题.
    if h.ContainsPTS && !h.ContainsDTS {

	    //h.PTS = uint64(data[9]&0x0E>>1) << 30
	    //h.PTS = h.PTS | (uint64(data[10])<<8|uint64(data[11]))>>1<<15
	    //h.PTS = h.PTS | (uint64(data[12])<<8|uint64(data[13]))>>1
        data[9-6]  = (0x02<<4) |  (uint8((h.PTS >> 30) << 1) & 0x0E) | 0x01
        data[10-6] = uint8( ( ((h.PTS >> 15 << 1) & 0x000000000000FFFF) >> 8) &  0x00000000000000FF)
        data[11-6] = uint8(   ((h.PTS >> 15 << 1) & 0x000000000000FFFF)  &  0x00000000000000FF) | 0x01

        data[12-6] = uint8( ( ((h.PTS << 1) & 0x000000000000FFFF) >> 8) &  0x00000000000000FF)
        data[13-6] = uint8(   ((h.PTS << 1) & 0x000000000000FFFF) &  0x00000000000000FF) | 0x01


	    h.HeaderLength = 5
	    data[2] = h.HeaderLength;
        //h.PTS = uint64(data[9]&0x0E>>1) << 30
        //h.PTS = h.PTS | (uint64(data[10])<<8|uint64(data[11]))>>1<<15
        //h.PTS = h.PTS | (uint64(data[12])<<8|uint64(data[13]))>>1
    }

    if h.ContainsPTS && h.ContainsDTS {
        //TODO: 需要高清楚这个0x03是什么意思.参考srs
        data[9-6]  = (0x03<<4) |  (uint8((h.PTS >> 30) << 1) & 0x0E) | 0x01
        data[10-6] = uint8( ( ((h.PTS >> 15 << 1) & 0x000000000000FFFF) >> 8) &  0x00000000000000FF)
        data[11-6] = uint8(   ((h.PTS >> 15 << 1) & 0x000000000000FFFF)  &  0x00000000000000FF) | 0x01

        data[12-6] = uint8( ( ((h.PTS << 1) & 0x000000000000FFFF) >> 8) &  0x00000000000000FF)
        data[13-6] = uint8(   ((h.PTS << 1) & 0x000000000000FFFF) &  0x00000000000000FF) |  0x01

        //TODO: 需要高清楚这个0x01是什么意思.参考srs
        data[14-6]  = (0x01<<4) |  (uint8((h.DTS >> 30) << 1) & 0x0E) | 0x01
        data[15-6] = uint8( ( ((h.DTS >> 15 << 1) & 0x000000000000FFFF) >> 8) &  0x00000000000000FF)
        data[16-6] = uint8(   ((h.DTS >> 15 << 1) & 0x000000000000FFFF)  &  0x00000000000000FF) | 0x01

        data[17-6] = uint8( ( ((h.DTS << 1) & 0x000000000000FFFF) >> 8) &  0x00000000000000FF)
        data[18-6] = uint8(   ((h.DTS << 1) & 0x000000000000FFFF) &  0x00000000000000FF) |  0x01

	    h.HeaderLength = 10
	    data[2] = h.HeaderLength;
        //h.DTS = uint64(data[14]&0x0E>>1) << 30
        //h.DTS = h.DTS | (uint64(data[15])<<8|uint64(data[16]))>>1<<15
        //h.DTS = h.DTS | (uint64(data[17])<<8|uint64(data[18]))>>1
    }
}
