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

    OnNewPacket func(*Packet)
    OnNewPAT    func(*ProgramAssociationTable)
    OnNewPMT    func(*ProgramMapTable)
}

func NewWriter(r io.Writer, c func(*Packet), a func(*ProgramAssociationTable), m func(*ProgramMapTable)) *Writer{
    return &Writer{
        writer:      r,
        OnNewPacket: c,
        OnNewPAT:    a,
        OnNewPMT:    m}
}

func (w *Writer) WritePacket(p *Packet) (error) {

    data := writePacket(p)
    _,err := w.writer.Write(data)

    return err

}

