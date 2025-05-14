package autoload

import (
	"fmt"

	"github.com/QYUbit/lib/go/envfile"
)

func init() {
	if err := envfile.Load(""); err != nil {
		fmt.Print(fmt.Errorf("failed autoload: %w", err))
	}
}
