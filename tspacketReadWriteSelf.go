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

    w,err2 := os.Create("./outputself.ts")
    if err2 != nil {
        fmt.Println("Create output.ts error!!!!")
        return
    }
    defer w.Close()

	t := ts.NewReader(r, displayTSPacket, displayPAT, displayPMT)
	p := pes.NewReader(t, displayPES)

    wt := ts.NewWriter(w, displayTSPacket, displayPAT, displayPMT)
    wp := pes.NewWriter(wt,displayPES)

	wtp := wt.MakePATPacket()
	wt.WritePacket(wtp)

	wtp = wt.MakePMTPacket()
	wt.WritePacket(wtp)
	for {
		//_, err := p.Next()
		pesPacket, err := p.PesNext()
		if err != nil {
			fmt.Println(err)
			return
		}
        if pes.STREAMID_AAC_CODE == pesPacket.StreamID {
            wpacket,_:= wp.WriteAACRawData(pesPacket.Payload,false,pesPacket.Header.PTS)
            //wpacket,_:= wp.WriteAACRawData(pesPacket.Payload,false,0)
            wp.WritePacketToTS(wpacket)


        } else if pes.STREAMID_AVC_H264_CODE == pesPacket.StreamID {
            wpacket,_ := wp.WriteAVCRawData(pesPacket.Payload,false,pesPacket.Header.PTS,pesPacket.Header.DTS)
            //wpacket,_ := wp.WriteAVCRawData(pesPacket.Payload,false,0,0)
            wp.WritePacketToTS(wpacket)

        } else {
            fmt.Println("Unknow streamID...............")
        }
		//pesfile.Write(pes.Body)
		//rePesfile.Write(pes.Payload)

		//data := repes.WritePacket(pes)
		//rePesfile.Write(data)

		//if !bytes.Equal(data,pes.Body) {
		//	fmt.Println("============================================================")
		//	displayPES(pes)
		//	fmt.Println("============================================================")
		//	fmt.Println(hex.Dump(data))
		//	fmt.Println("============================================================")
		//}
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
	fmt.Println("============================================================")
	fmt.Printf("TS packet [%d]\n", TSIndex)
	fmt.Println("============================================================")
	fmt.Printf("%s\n", p)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("Payload (%d) \n", len(p.Payload))
	fmt.Println("------------------------------------------------------------")
	displayPayload(p.Payload)
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
