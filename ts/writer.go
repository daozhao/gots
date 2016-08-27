package ts

import (
    "io"
)

type Writer struct {
    writer io.Writer
    //offset int

    Packet *Packet
    PAT    *ProgramAssociationTable
    PMT    *ProgramMapTable

    continuityCounter  uint8

    OnNewPacket func(*Packet)
    OnNewPAT    func(*ProgramAssociationTable)
    OnNewPMT    func(*ProgramMapTable)
}

func NewWriter(r io.Writer, c func(*Packet), a func(*ProgramAssociationTable), m func(*ProgramMapTable)) *Writer{
    return &Writer{
        writer:      r,
        OnNewPacket: c,
        OnNewPAT:    a,
        OnNewPMT:    m,
        continuityCounter:1}
}

func (w *Writer) WritePacket(p *Packet) (error) {

    data := WritePacket(p)
    _,err := w.writer.Write(data)

    return err

}

func (w *Writer) GetContinuitycounter()(uint8){
    c := w.continuityCounter
    w.continuityCounter += 1
    return c
}
func (w *Writer) MakePATPacket()(p *Packet){

    pat := MakePogramAssociationTable(256);
    w.PAT = pat

    p = MakePacket(ProgramAssociationTableID,true,true,w.continuityCounter,nil,0)
    //TODO: 这个counter计算有点问题.
    w.continuityCounter += 1
    p.Payload = WritePogramAssociationTable(pat)

    return p

}

func (w *Writer) MakePMTPacket()(p *Packet){
    pmt := MakeProgramMapTable(w.PAT.Programs[0].ProgramMapPID)
    w.PMT = pmt

    p = MakePacket(w.PAT.Programs[0].ProgramMapPID,true,true,w.continuityCounter,nil,0)
    w.continuityCounter += 1
    p.Payload = WriteProgramMapTable(pmt)

    return p
}

//func (w *Writer) MakePlayloadPacket()(p *Packet){
//
//}

