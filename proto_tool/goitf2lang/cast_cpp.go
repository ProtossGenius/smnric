package goitf2lang

import (
	"fmt"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/smnric/smn_pglang"
)

func ToCppType(goType string, clt ...bool) string {
	if strings.HasPrefix(goType, "int") || strings.HasPrefix(goType, "uint") {
		if goType == "int" || goType == "uint" {
			return goType + "64" + "_t"
		}
		return goType + "_t"
	}
	switch goType {
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "string":
		return "std::string"
	case "net.Conn":
		if len(clt) > 0 && clt[0] {
			return "std::function<int(std::shared_ptr<smnet::Conn>)>"
		}

		return "std::shared_ptr<smnet::Conn> "
	}
	if strings.Contains(goType, "*") {
		goType = strings.Replace(goType, "*", "", -1)
		goType = strings.Replace(goType, ".", "::", -1)
	}

	return goType
}

func CppBuiltInType(t string) bool {
	switch t {
	case "int", "unsigned int", "int32_t", "uint32_t", "long", "unsigned long", "long long", "unsigned long long",
		"int8_t", "uint8_t", "int16_t", "uint16_t", "int64_t", "uint64_t", "double", "float", "char", "unsigned char",
		"short", "unsigned short", "std::size_t", "bool":
		return true
	}
	return false
}

func ToCppVarDef(vd *smn_pglang.VarDef, clt ...bool) *smn_pglang.VarDef {
	res := &smn_pglang.VarDef{Var: vd.Var}
	vdType := vd.Type
	if vd.ArrSize != 0 {
		arrV := 0
		for strings.Contains(vdType, "[]") {
			vdType = strings.Replace(vdType, "[]", "", 1)
			arrV++
		}
		res.Type = fmt.Sprintf("%s%s%s", strings.Repeat("std::vector<", arrV), ToCppType(vdType, clt...), strings.Repeat(">", arrV))
	} else {
		res.Type = ToCppType(vdType, clt...)
	}
	return res
}

func ToCppParam(param []*smn_pglang.VarDef, clt ...bool) string {
	if len(param) == 0 {
		return "void"
	}
	list := make([]string, len(param))
	for i, p := range param {
		vp := ToCppVarDef(p, clt...)
		if vp.Var == "" {
			vp.Var = fmt.Sprintf("sm_p%d", i)
		}
		if CppBuiltInType(vp.Type) || p.Type == "net.Conn" {
			list[i] = fmt.Sprintf("%s %s", vp.Type, vp.Var)
		} else {
			list[i] = fmt.Sprintf("const %s& %s", vp.Type, vp.Var)
		}
	}

	return strings.Join(list, ", ")
}

func TooCppRet(rets []*smn_pglang.VarDef, pkg, itfName, fName string) string {
	if len(rets) == 0 {
		return "void"
	}
	if len(rets) == 1 {
		p := ToCppVarDef(rets[0])
		return p.Type
	}
	return fmt.Sprintf("rip_%s::%s_%s_Ret", pkg, itfName, fName)

}

func hasPkg(typ string) (pkg string) {
	if !strings.Contains(typ, ".") {
		return ""
	}

	pkg = strings.Split(typ, ".")[0]
	pkg = strings.ReplaceAll(pkg, "*", "")
	pkg = strings.ReplaceAll(pkg, "[]", "")
	pkg = strings.TrimSpace(pkg)

	return pkg
}

//CppNeedInc cpp func need include.
func CppNeedInc(itf *smn_pglang.ItfDef, openRet bool, needPrmPkg bool, init ...string) map[string]bool {
	result := make(map[string]bool)

	for _, str := range init {
		result[str] = true
	}

	if needPrmPkg {
		result[fmt.Sprintf("#include \"pb/rip_%s.pb.h\"", itf.Package)] = true
	}

	for _, f := range itf.Functions {
		if len(f.Returns) > 1 {
			result[fmt.Sprintf("#include \"pb/rip_%s.pb.h\"", itf.Package)] = true
		}

		for idx, ret := range f.Returns {
			if !openRet && idx > 0 {
				break
			}

			if pkg := hasPkg(ret.Type); pkg != "" {
				result[fmt.Sprintf(`#include "pb/%s.pb.h"`, pkg)] = true
			}
		}

		for _, prm := range f.Params {
			if strings.Contains(prm.Type, "net.Conn") {
				continue
			}

			if pkg := hasPkg(prm.Type); pkg != "" {
				result[fmt.Sprintf(`#include "pb/%s.pb.h"`, pkg)] = true
			}
		}
	}
	return result
}

// WriteCppItf write cpp interface from go-interface-define.
func WriteCppItf(dir, pkg string, itf *smn_pglang.ItfDef) {
	checkerr(smn_file.MakeSureDirExist(dir))
	dir += pkg + "."
	f, err := smn_file.CreateNewFile(dir + itf.Name + ".h")
	checkerr(err)
	defer f.Close()
	writef := func(s string, a ...interface{}) {
		_, err := f.WriteString(fmt.Sprintf(s, a...))
		checkerr(err)
	}
	writef(`#pragma once
#include<vector>
#include"smncpp/socket_itf.h"
`)

	for inc := range CppNeedInc(itf, false, false) {
		writef(inc)
	}

	writef(`
namespace %s{

`, pkg)
	defer writef("}//namespace %s", pkg)

	writef("class %s {\npublic:\n", itf.Name)
	defer writef("};\n")
	for _, f := range itf.Functions {
		writef("\tvirtual %s %s(%s) = 0;\n", TooCppRet(f.Returns, pkg, itf.Name, f.Name), f.Name, ToCppParam(f.Params))
	}
	writef("\n\tvirtual ~%s(){}", itf.Name)
}

func WriteCppPkg(out, pkg string, list []*smn_pglang.ItfDef) {
	for _, itf := range list {
		WriteCppItf(out, pkg, itf)
	}
}
