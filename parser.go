package main

import (
	"flag"
	"fmt"
	"github.com/damienlevin/gots/pes"
	"github.com/damienlevin/gots/ts"
	"io"
	"net/http"
	"os"
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

    h264,err2 := os.Create("./h264output.raw")
    if  nil != err2 {
        fmt.Println(err2)
    }
    defer h264.Close()
    aacfile,err3 := os.Create("./AACoutput.aac")
    if  nil != err3 {
        fmt.Println(err3)
    }
    defer aacfile.Close()
	t := ts.NewReader(r, displayTSPacket, displayPAT, displayPMT)
	p := pes.NewReader(t, displayPES)
	for {
		//_, err := p.Next()
		pes, err := p.PesNext()
		if err != nil {
			fmt.Println(err)
			return
		}
        if 0xe0 == pes.StreamID {
            h264.Write(pes.Payload)
        }
		if 0xc0 == pes.StreamID {
			aacfile.Write(pes.Payload)
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

func displayPES(m *pes.Packet) {
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("PES packet [%d]\n", PESIndex)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("%s\n", m)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++")
	for i,p := range m.Ts {
		displayPESTSPacket(p,i)
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++")

	PESIndex++
}

func displayPESTSPacket(p *ts.Packet,TSIndex int) {
    fmt.Printf("TS packet [%d]\n", TSIndex)
    fmt.Printf("%s\n", p)
    fmt.Println("------------------------------------------------------------")
    fmt.Printf("Payload (%d) \n", len(p.Payload))
    fmt.Println(hex.Dump(p.Payload))
    fmt.Println("------------------------------------------------------------")
    TSIndex++
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
