package local

import (
	"fmt"
	"os"

	pst "github.com/archeryue/go-postern/postern"
)

func main() {
	path := os.Args[1]
	config, err := pst.LoadConfig(path)
	if err != nil {
		fmt.Println("Config Error: " + path)
		return
	}

	cipher := pst.NewCipher(config.Key)

	fmt.Println("vim-go")
}
