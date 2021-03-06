package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ProtossGenius/smnric/analysis/smn_rpc_itf"
	"github.com/ProtossGenius/smnric/proto_tool/itf2proto"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	i := flag.String("i", "./rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./datas/proto/", "proto output dir.")
	flag.Parse()

	err := os.MkdirAll(*o, os.ModePerm)
	checkerr(err)
	itfs, err := smn_rpc_itf.GetItfListFromDir(*i)
	checkerr(err)

	for _, list := range itfs {
		err = itf2proto.WriteProto(*o, list)
		if err != nil {
			fmt.Println("[ERROR] when write : ", err.Error())
		}
	}
}
