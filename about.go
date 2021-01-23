package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	list := []string{
		"github.com/ProtossGenius/smnric/smn_str",
		"github.com/ProtossGenius/smnric/smn_file",
		"github.com/ProtossGenius/smnric/smn_data",
		"github.com/ProtossGenius/smnric/smn_str_rendering",
		"github.com/ProtossGenius/smnric/smn_muti_write_cache",
		"github.com/ProtossGenius/smnric/smn_net",
		"github.com/ProtossGenius/smnric/smn_stream",
		"github.com/ProtossGenius/smnric/smn_err",
		"github.com/ProtossGenius/smnric/smn_exec",
	}

	exist := map[string]string{}
	for _, str := range list {
		exist[str] = strings.ReplaceAll(str, "github.com/ProtossGenius/smnric", "github.com/ProtossGenius/SureMoonNet/basis")
	}

	_, err := smn_file.DeepTraversalDir("./", func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}

		if !strings.HasSuffix(info.Name(), ".go") {
			return smn_file.FILE_DO_FUNC_RESULT_NO_DEAL
		}
		fmt.Println(path)
		data, err := smn_file.FileReadAll(path)
		check(err)
		f, err := smn_file.CreateNewFile(path)
		check(err)
		defer f.Close()
		str := string(data)
		for key, val := range exist {
			str = strings.ReplaceAll(str, key, val)
		}
		_, err = f.WriteString(str)
		check(err)

		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})

	check(err)
}
