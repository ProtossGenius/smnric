package itf2rpc

import (
	"fmt"

	"github.com/ProtossGenius/smnric/smn_pglang"
)

//FWriteRPC write RPC.
type FWriteRPC func(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error

//TargetMap from target to func.
var TargetMap = map[string]FWriteRPC{
	"go_s":  GoSvr,
	"go_c":  GoClient,
	"go_ac": GoAsynClient,
	"cpp_c": CppClient,
	"cpp_s": CppServer,
}

//Write itf to rpc.
func Write(target, path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error {
	f, ok := TargetMap[target]
	if !ok {
		return fmt.Errorf("Can't Found Target %s", target)
	}
	return f(path, module, itfFullPkg, itf)
}
