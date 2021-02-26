package proto_compile

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_exec"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/smnric/analysis/proto_msg_map"
)

func protoHead(pkg, module string) string {
	return fmt.Sprintf(`syntax = "proto3";
option java_package = "pb";
option java_outer_classname="%s";
option go_package = "%s/pb/%s;%s";
package %s;

`, pkg, module, pkg, pkg, pkg)
}

//生成字典协议
func Dict(in, module string) error {
	list, _, err := proto_msg_map.Dict(in)
	if err != nil {
		return err
	}

	file, err := smn_file.CreateNewFile(in + "/smn_dict.proto")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(protoHead("smn_dict", module))

	if err != nil {
		return err
	}

	_, err = file.WriteString("enum EDict{\n")
	if err != nil {
		return err
	}

	for _, val := range list {
		_, err = file.WriteString(fmt.Sprintf("\t%s = %d;\n", val.Name, val.Id))
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString("}\n")

	return err
}

func getPkg(path string) string {
	data, err := smn_file.FileReadAll(path)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "package") {
			pkg := strings.Split(line[7:], ";")[0]
			pkg = strings.TrimSpace(pkg)

			return pkg
		}
	}

	return ""
}

type compileFunc func(in, out, exportPath, ignoreDir, comp string) error

var CompileMap = map[string]compileFunc{}

func CppCompile(in, out, goMoudle, ignoreDir, comp string) error {
	var retErr error

	_, err := smn_file.DeepTraversalDir(in, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() && info.Name() == ignoreDir {
			return smn_file.FILE_DO_FUNC_RESULT_NO_DEAL
		}

		if strings.HasSuffix(info.Name(), ".proto") {
			if err := smn_exec.EasyDirExec(".", "protoc", fmt.Sprintf(comp, out), "-I", in, path); err != nil {
				retErr = err
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
		}

		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	if err != nil {
		return err
	}

	return retErr
}

func Compile(protoDir, codeOutPath, goMod, lang string) error {
	if !smn_file.IsFileExist(codeOutPath) {
		err := os.MkdirAll(codeOutPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := Dict(protoDir, goMod); err != nil {
		return err
	}

	comp := "--" + lang + "_out=%s" // "--go_out=%s"
	extPath := strings.ReplaceAll(goMod, "\\", "/") + "/" + strings.ReplaceAll(codeOutPath, "./", "")
	ignoreDir := strings.Split(extPath, "/")[0]

	if f, ok := CompileMap[lang]; ok {
		return f(protoDir, codeOutPath, extPath, ignoreDir, comp)
	}

	return CppCompile(protoDir, codeOutPath, extPath, ignoreDir, comp)
}
