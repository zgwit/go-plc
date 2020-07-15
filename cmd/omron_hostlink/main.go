package main

import (
	"fmt"
	"github.com/zgwit/go-plc/adapter/omron"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/link"
)

func main()  {

	l := link.NewSerial()
	l.Name = "COM3"
	l.Parity = 'E'
	l.StopBits = 2

	//defer l.Close()
	if e := l.Open(); e!=nil {
		fmt.Println(e)
		return
	}

	a := omron.NewHostLink(l)

	//b, e := a.Read("D100", 4)
	b, e := a.ReadBit(omron.HR, 0, 0, 8)
	fmt.Println("ReadBit DM 10", helper.BooleansToBytes(b), e)

	w, e := a.ReadWord(omron.HR, 0, 2)
	fmt.Println("ReadWord DM 10", w, e)

}