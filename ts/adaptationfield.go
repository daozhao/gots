package ts

type AdaptationField struct {
	AdaptationFieldLength             uint8
	DiscontinuityIndicator            bool
	RandomAccessIndicator             bool
	ElementaryStreamPriorityIndicator bool
	ContainsPCR                       bool
	ContainsOPCR                      bool
	ContainsSplicingPoint             bool
	ContainsTransportPrivateData      bool
	ContainsAdaptationFieldExtension  bool
	PCR                               []byte
	OPCR                              []byte
	SpliceCountdown                   uint8
	TransportPrivateDataLenght        uint8
	PrivateData                       []byte
	//TODO parse AdaptationFieldExtention
}

func newAdaptationField(data []byte) *AdaptationField {
	af := &AdaptationField{AdaptationFieldLength: uint8(data[0])}
	if af.AdaptationFieldLength == 0 {
		return af
	}
	af.DiscontinuityIndicator = data[1]&0x80>>7 == 1
	af.RandomAccessIndicator = data[1]&0x40>>6 == 1
	af.ElementaryStreamPriorityIndicator = data[1]&0x20>>5 == 1
	af.ContainsPCR = data[1]&0x10>>4 == 1
	af.ContainsOPCR = data[1]&0x08>>3 == 1
	af.ContainsSplicingPoint = data[1]&0x04>>2 == 1
	af.ContainsTransportPrivateData = data[1]&0x02>>1 == 1
	af.ContainsAdaptationFieldExtension = data[1]&0x01 == 1

	i := 2
	if af.ContainsPCR {
		af.PCR = data[i : i+6]
		i += 6
	}
	if af.ContainsOPCR {
		af.OPCR = data[i : i+6]
		i += 6
	}
	if af.ContainsSplicingPoint {
		af.SpliceCountdown = data[i]
		i += 1
	}
	//TODO: 这里估计有问题,重复了.
	if af.ContainsSplicingPoint {
		af.SpliceCountdown = data[i]
		i += 1
	}
	if af.ContainsTransportPrivateData {
		af.TransportPrivateDataLenght = data[i]
		af.PrivateData = data[i:i+int(af.TransportPrivateDataLenght)]
	}
	return af
}

func writeAdaptationField(data []byte,af *AdaptationField){

	data[0] = af.AdaptationFieldLength

	if af.AdaptationFieldLength == 0 {
		return
	}
	data[1] = 0
	if af.DiscontinuityIndicator {
		data[1] |= 0x80
	}
	if af.RandomAccessIndicator {
		data[1] |= 0x40
	}
	if af.ElementaryStreamPriorityIndicator {
		data[1] |= 0x20
	}
	if af.ContainsPCR {
		data[1] |= 0x10
	}
	if af.ContainsOPCR {
		data[1] |= 0x08
	}
	if af.ContainsSplicingPoint {
		data[1] |= 0x04
	}
	if af.ContainsTransportPrivateData {
		data[1] |= 0x02
	}
    if af.ContainsAdaptationFieldExtension {
        data[1] |= 0x01
    }
    i := 2
	if af.ContainsPCR {
		copy(data[i : i+6],af.PCR)
		i += 6
	}
	if af.ContainsOPCR {
		copy(data[i : i+6],af.OPCR)
		i += 6
	}
	if af.ContainsSplicingPoint {
		data[i] = af.SpliceCountdown
		i += 1
	}
	//TODO: 这里估计有问题,重复了.
	if af.ContainsSplicingPoint {
		data[i] = af.SpliceCountdown
		i += 1
	}
	if af.ContainsTransportPrivateData {
		data[i] = af.TransportPrivateDataLenght
		copy(data[i:i+int(af.TransportPrivateDataLenght)],af.PrivateData)
		i += int(af.TransportPrivateDataLenght)
	}

	for i <= int(af.AdaptationFieldLength) {
		data[i] = 0xFF
		i += 1
	}



}
