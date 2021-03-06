package smn_pglang

import (
	"github.com/ProtossGenius/SureMoonNet/basis/smn_data"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"os"
	"strings"
)

func LoadSystem(folderPath, sysName string) (sMap map[string]interface{}, err error) {
	sMap = make(map[string]interface{})
	smn_file.DeepTraversalDir(folderPath, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if !info.IsDir() {
			ts := &SystemDef{}
			bytes, err := smn_file.FileReadAll(path)
			if iserr(err) {
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
			smn_data.GetDataFromStr(string(bytes), &ts)
			ts.Name = strings.Split(info.Name(), ".")[0]
			sMap[ts.Name] = ts
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	return
}
