package main
import (
	"github.com/smartwalle/going/open_platform"
//	"fmt"
//	"github.com/smartwalle/going/http"
	"fmt"
)

func main() {

//	var wb = open_platform.Weibo{}
//	var r, er = wb.CheckAccessToken("2.00kLXhIG9Cf3tD0015e73a690SCzU5")
//	fmt.Println(r, er)

	var fb = open_platform.Facebook{}
	var s, err = fb.CheckAccessToken("CAAU7ErrNaQcBADckp3MKVj6iyLVAWUusD4xdLPbjJ97Aq6W8NuiEBMSXZBYfsLV8UN5jBsNXWZA2yfwfZCuuPwZBWiaNk8xlE3ZAUhUSRvZB6ZBUdt2lrvBf3P98N326c18iQXdJaZAoyAIUmBGHUXG4Fic4uElQZAorxRUBQXhFbZADnAGsZAmygtKxWiIZCfiFy3tLQAeGZC2YA3uuAmBefVMtR")
	fmt.Println(s, err)

	s, err = fb.GetUserInfo("564882523638640", "CAAU7ErrNaQcBADckp3MKVj6iyLVAWUusD4xdLPbjJ97Aq6W8NuiEBMSXZBYfsLV8UN5jBsNXWZA2yfwfZCuuPwZBWiaNk8xlE3ZAUhUSRvZB6ZBUdt2lrvBf3P98N326c18iQXdJaZAoyAIUmBGHUXG4Fic4uElQZAorxRUBQXhFbZADnAGsZAmygtKxWiIZCfiFy3tLQAeGZC2YA3uuAmBefVMtR")
	fmt.Println(s, err)
}
