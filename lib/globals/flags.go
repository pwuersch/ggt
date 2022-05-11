package globals

import (
	"fmt"

	"github.com/logrusorgru/aurora/v3"
)

var (
	RootDebug      bool
	CloneRemoteUrl string
	CloneDestDir   string
)

func Debug(message string) {
	if RootDebug {
		fmt.Println(aurora.Gray(12, fmt.Sprintf("[debug] %s", message)))
	}
}
