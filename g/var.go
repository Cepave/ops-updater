package g

import (
	"github.com/toolkits/file"
)

var SelfDir string

func InitGlobalVariables() {
	SelfDir = file.SelfDir()
}
