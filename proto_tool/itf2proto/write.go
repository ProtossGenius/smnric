package itf2proto

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_muti_write_cache"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str"
	"github.com/ProtossGenius/smnric/smn_pglang"
)

// WriteProto write itf to proto.
func WriteProto(outDir, module string, list []*smn_pglang.ItfDef) (err error) {
	pkg := list[0].Package
	oPath := outDir + "/rip_" + pkg + ".proto"

	var file *os.File

	if file, err = smn_file.CreateNewFile(oPath); err != nil {
		return err
	}

	defer file.Close()

	w := smn_muti_write_cache.NewFileMutiWriteCache()
	writeLine := func(f string, a ...interface{}) {
		_, _ = w.WriteTail(fmt.Sprintf(f, a...) + "\n")
	}

	w.Append(smn_muti_write_cache.NewStrCache(fmt.Sprintf(`syntax = "proto3";
	option java_package = "pb";
	option java_outer_classname="rip_%s";
	option go_package = "%s/pb/rip_%s;rip_%s";
	package rip_%s;
`, pkg, module, pkg, pkg, pkg)))
	impts := smn_muti_write_cache.NewStrCache()
	impMap := make(map[string]bool)
	checkImport := func(typ string) {
		if typ == "net.Conn" {
			return
		}
		nameList := strings.Split(typ, ".")
		if len(nameList) == 1 {
			return
		}
		if !impMap[nameList[0]] {
			impts.WriteTail(fmt.Sprintf("import \"%s.proto\";", nameList[0]) + "\n")
			impMap[nameList[0]] = true
		}
	}
	w.Append(impts)
	for _, itf := range list {
		for _, mtd := range itf.Functions {
			if mtd.Name[0] < 'A' || mtd.Name[0] > 'Z' {
				log.Printf("warning! mtd.Name %s's first letter not upper\n", mtd.Name)
				continue
			}
			writeLine("message %s_%s_Prm {", itf.Name, mtd.Name)
			for i, prm := range mtd.Params {
				isArray, typ := smn_str.ProtoUseDeal(prm.Type)
				if typ == "net.Conn" {
					continue
				}
				checkImport(typ)
				rpt := ""
				if isArray {
					rpt = "repeated "
				}
				writeLine("\t%s%s %s = %d;", rpt, typ, prm.Var, i+1)
			}
			writeLine("}")
			writeLine("message %s_%s_Ret {", itf.Name, mtd.Name)
			for i, res := range mtd.Returns {
				isArray, typ := smn_str.ProtoUseDeal(res.Type)
				checkImport(typ)
				rpt := ""
				if isArray {
					rpt = "repeated "
				}
				writeLine("\t%s%s %s = %d;", rpt, typ, res.Var, i+1)
			}
			writeLine("}")
		}
	}
	w.Output(file)
	return nil
}
