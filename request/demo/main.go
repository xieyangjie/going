package main

import (
	"github.com/smartwalle/going/request"
	"fmt"
)

func main() {
	var r = request.NewRequest("GET", "http://www.baidu.com")

	rep := r.Exec()
	fmt.Println(rep.String())
}
