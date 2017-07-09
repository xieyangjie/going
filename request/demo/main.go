package main

import (
	"github.com/smartwalle/social/request"
	"fmt"
)

func main() {
	var r = request.NewRequest("GET", "https://order.smoktech.com/api/retail/products")
	r.AddParam("currency_type_id", "8")

	rep := r.Exec()

	fmt.Println(rep.Bytes())
	fmt.Println(rep.String())
}
