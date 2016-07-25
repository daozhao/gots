package pes

import (
	"github.com/damienlevin/gots/ts"
)

type Reader struct {
	reader *ts.Reader

	//Packet      *Packet
    PrevPacket map[uint16]*Packet
	OnNewPacket func(*Packet)
}

func NewReader(r *ts.Reader, c func(*Packet)) *Reader {
	return &Reader{
		reader:      r,
		PrevPacket: make(map[uint16]*Packet),
		OnNewPacket: c}
}

func (r *Reader) Next() (*Packet, error) {
	var p *Packet
	var err error
	for p == nil {
		if _, err := r.reader.Next(); err != nil {
			return nil, err
		}
		t := r.reader.Packet
		pmt := r.reader.PMT
		if t != nil && pmt != nil && pmt.HasElementaryStream(t.PID) {
			if t.PayloadUnitStartIndicator {
				p, err = newPacket(t.Payload)
				if err != nil {
					return nil, err
				}
				r.OnNewPacket(p)
				return p, nil
			}
		}
	}
	return nil, nil
}

func (r *Reader) PesNext() (*Packet, error) {
	var p *Packet

	//pStream = make(map[uint16]*Packet)
	var err error

    for p == nil {
		if _, err := r.reader.Next(); err != nil {
            //TODO:  这里需要输出所有的pes包.
            for key,prev := range r.PrevPacket {
				if nil != prev {
					delete(r.PrevPacket,key)
					return prev,nil
				}
            }
			return nil, err
		}
		t := r.reader.Packet
		pmt := r.reader.PMT
		if t != nil && pmt != nil && pmt.HasElementaryStream(t.PID) {
			prev,exit := r.PrevPacket[t.PID]
			if t.PayloadUnitStartIndicator {
                if exit && (nil != prev) {
                    //fmt.Println("++++++++++++++++out put pes=====")
                    r.OnNewPacket(prev)
	                p = prev
                }
				prev, err = newPacket(t.Payload)
				if err != nil {
                    return nil,err
                }
                prev.Body = make([]byte,0)
                prev.Body = append(prev.Body,t.Payload...)

                prev.Ts   = make([]*ts.Packet,0)
                prev.Ts   = append(prev.Ts,t)
                r.PrevPacket[t.PID] = prev
			} else {
				if exit {
                    //fmt.Println("++++++++++++++++had exit=====")
					prev.Payload = append(prev.Payload,t.Payload...)
                    prev.Body = append(prev.Body,t.Payload...)
                    prev.Ts   = append(prev.Ts,t)
				} else {
					//fmt.Println("+++++++++++++++++TS pes reader error:no prev ts trunk=====")
				}


			}
		}
	}
	return p, nil
}
