package main

import (
	"flag"
	"fmt"
	"github.com/damienlevin/gots/pes"
	"github.com/damienlevin/gots/ts"
	"io"
	"net/http"
	"os"
	"bytes"
	"encoding/hex"
)

var TSIndex = 1
var PESIndex = 1

func main() {
	r, err := newReader()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

    pesfile,err2 := os.Create("./pes.file")
    if  nil != err2 {
        fmt.Println(err2)
    }
    defer pesfile.Close()

    rePesfile,err3 := os.Create("./repes.file")
    if  nil != err3 {
        fmt.Println(err3)
    }

	wt := ts.NewWriter(rePesfile, displayTSPacket, displayPAT, displayPMT)
	repes := pes.NewWriter(wt,displayPES)

    defer rePesfile.Close()
	t := ts.NewReader(r, displayTSPacket, displayPAT, displayPMT)
	p := pes.NewReader(t, displayPESTS)

	for {
		var pesRemake *pes.Packet
		//_, err := p.Next()
		pes, err := p.PesNext()
		if err != nil {
			fmt.Println(err)
			return
		}
		pesfile.Write(pes.Body)
		//rePesfile.Write(pes.Payload)


		if 0xE0 == pes.StreamID {
			pesRemake,_ = repes.WriteAVCRawData(pes.Payload,pes.Header.DataAlignmentIndicator,pes.Header.PTS,pes.Header.DTS)
		}
		if 0xC0 == pes.StreamID {
            pesRemake,_ = repes.WriteAACRawData(pes.Payload,pes.Header.DataAlignmentIndicator,pes.Header.PTS)
		}
		data := repes.WritePacket(pesRemake)
		rePesfile.Write(data)

		if !bytes.Equal(data,pes.Body) {
			fmt.Println("============================================================")
			displayPES(pes)
			fmt.Println("============================================================")
			fmt.Println(hex.Dump(data))
			fmt.Println("============================================================")
		}
	}
}

func newReader() (io.ReadCloser, error) {
	u := flag.String("u", "https://devimages.apple.com.edgekey.net/streaming/examples/bipbop_4x3/gear1/fileSequence179.ts", "TS URL to parse.")
	f := flag.String("f", "", "TS file to parse.")
	flag.Parse()

	switch {
	case *f != "":
		fmt.Printf("Parsing file %s\n", *f)
		r, err := os.Open(*f)
		if err != nil {
			return nil, err
		}
		return r, nil
	default:
		fmt.Printf("Parsing URL %s\n", *u)
		client := &http.Client{}
		rsp, err := client.Get(*u)
		if err != nil {
			return nil, err
		}
		return rsp.Body, nil
	}
}

func displayTSPacket(p *ts.Packet) {
	//fmt.Println("============================================================")
	//fmt.Printf("TS packet [%d]\n", TSIndex)
	//fmt.Println("============================================================")
	//fmt.Printf("%s\n", p)
	//fmt.Println("------------------------------------------------------------")
	//fmt.Printf("Payload (%d) \n", len(p.Payload))
	//fmt.Println("------------------------------------------------------------")
	//displayPayload(p.Payload)
	TSIndex++
}

func displayPAT(m *ts.ProgramAssociationTable) {
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("PAT\n")
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%s\n", m)
}

func displayPMT(m *ts.ProgramMapTable) {
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("PMT\n")
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%s\n", m)
}

func displayPESTS(m *pes.Packet) {
	//fmt.Println("------------------------------------------------------------")
	//fmt.Printf("PES packet [%d]\n", PESIndex)
	//fmt.Println("------------------------------------------------------------")
	//fmt.Printf("%s\n", m)
	//fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++")
	//for i,p := range m.Ts {
	//	displayPESTSPacket(p,i)
	//}
	//fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++")
    //
	//PESIndex++
}
func displayPES(m *pes.Packet) {
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("PES packet [%d]\n", PESIndex)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%s\n", m)
	//fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++")
	//for i,p := range m.Ts {
	//	displayPESTSPacket(p,i)
	//}
	//fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++")

	PESIndex++
}
func displayPESTSPacket(p *ts.Packet,TSIndex int) {
    //fmt.Printf("TS packet [%d]\n", TSIndex)
    //fmt.Printf("%s\n", p)
    //fmt.Println("------------------------------------------------------------")
    //fmt.Printf("Payload (%d) \n", len(p.Payload))
    //fmt.Println(hex.Dump(p.Payload))
    //fmt.Println("------------------------------------------------------------")
    //TSIndex++
}


func displayPayload(bytes []byte) {
	for i, b := range bytes {
		if (i+1)%16 == 0 || i+1 == len(bytes) {
			fmt.Printf("%02x \n", b)
			continue
		}
		fmt.Printf("%02x ", b)
	}
}
