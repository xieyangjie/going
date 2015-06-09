package main
import (
	"github.com/smartwalle/going/open_platform"
//	"fmt"
//	"github.com/smartwalle/going/http"
	"fmt"
)

func main() {

//	var r map[string]interface{}
//	var e error
//
//	fb := http.NewClient()
//	fb.SetMethod(http.K_HTTP_REQ_METHOD_GET)
//	fb.SetURLString("https://graph.facebook.com/v2.3/me")
//	fb.SetParam("access_token", "CAAU7ErrNaQcBAM7QArqQijX5Ts6PyE93iikKebMmosZCuOfZAGkiH3rZCU9pDaiW4BZCGoJIB4PUCs6fbcw3Xc3KWoxJZBUpVZAv3ZBlAgPbSqipA7Gx8lHpaTPkEU97uxeHpptfZBwv9GkZCYZCdUD6gPGVADcr1VHYRrVH1BMVZA5GBPwxZBdWOiPFBRP4u1gKZBJbkCGqGe4tkg17tQ61IsZBvv")
//	fb.SetParam("input_token", "CAAU7ErrNaQcBAM7QArqQijX5Ts6PyE93iikKebMmosZCuOfZAGkiH3rZCU9pDaiW4BZCGoJIB4PUCs6fbcw3Xc3KWoxJZBUpVZAv3ZBlAgPbSqipA7Gx8lHpaTPkEU97uxeHpptfZBwv9GkZCYZCdUD6gPGVADcr1VHYRrVH1BMVZA5GBPwxZBdWOiPFBRP4u1gKZBJbkCGqGe4tkg17tQ61IsZBvv")
//
//	r, e = fb.DoJsonRequest()
//	fmt.Println(r, e)

//	c := http.NewClient()
//	c.Method = "GET"
//	c.URLString = "https://graph.facebook.com/debug_token"
//	c.SetParam("input_token", "CAAU7ErrNaQcBAL7j2GZBcIi7cZBaCUOYNmxaLrAFLMdhImvh1nzrrnjALf5q9CPGB8ye9QAVMxa6JkbtClI9AuDCtOEVbZAkLej0TrhYLejtfcFuDAfgDlyDC4xgjOCMOo4wc7u27qe3HlhcozZC2WlZBJuMYytXXlG4kUlldZCQqKSiN8xEUJKi7ojwTJ0CvZBGY4YVVuURNVFfRZA9Igwv")
//	c.SetParam("access_token", "CAAU7ErrNaQcBAL7j2GZBcIi7cZBaCUOYNmxaLrAFLMdhImvh1nzrrnjALf5q9CPGB8ye9QAVMxa6JkbtClI9AuDCtOEVbZAkLej0TrhYLejtfcFuDAfgDlyDC4xgjOCMOo4wc7u27qe3HlhcozZC2WlZBJuMYytXXlG4kUlldZCQqKSiN8xEUJKi7ojwTJ0CvZBGY4YVVuURNVFfRZA9Igwv")
//
//	c.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36")
//	c.SetHeader("Accept", "*/*")
//	c.SetHeader("Accept-Encoding", "gzip, deflate, sdch")
//	c.SetHeader("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6")
//	c.SetHeader("CSP", "active")
//	c.SetHeader("Cookie", "datr=cJoHVQkiRg10CSnzAUDCwTmH; locale=en_US; lu=RhAL4fOp0NEM5lssfGMHhD9g; act=1433748238533%2F3; p=-2; presence=EM433753942EuserFA21B03506930917A2EstateFDsb2F0Et2F_5b_5dElm2FnullEuct2F14337237B00EtrFA2loadA2EtwF4030441192EatF1433753942545G433753942546CEchFDp_5f1B03506930917F4CC; c_user=100003506930917; fr=0UNYa7oWeL8XTsxRS.AWXCxYg0XTiRUcRXRw72YU1TnAA.BVcT3k.R2.FVx.0.AWXMiUHT; xs=170%3AehYKHM-3U8Qqjg%3A2%3A1433724299%3A15440; csm=2; s=Aa5772Zepk2wYyAS.BVdQmk")
//	v, e := c.DoJsonRequest()
//	fmt.Println(v, e)

//	v, e := http.DoPost("https://api.weibo.com/oauth2/get_token_info", map[string]string{"access_token": "2.00kLXhIG9Cf3tD0015e73a690SCzU5"})
//	fmt.Println(v, e)

	var fb = open_platform.Facebook{}
	var s, err = fb.CheckAccessToken("CAAU7ErrNaQcBAGJOO0Vp6EJ1VuQ6cwM7ZCWwsAZBnwTmlUxij7nRMVZBq0WKWzVbtZAmmr2wWCG7MU2yVJuOZBD0kLlOUSiuEUZBZB5hfAXOJBcyuJZBZClWUlefov5of7JKOOjHJTsahrFD46kiRc1mYwdgIyRt4krpgQH6VRzm0S3UfYRVN3kc1LivTPwdZC502SLR6ZCgy0CgTnZCwQZCdZCJOP")
	fmt.Println(s, err)

//	var wb = open_platform.Weibo{}
//	var s, err = wb.CheckAccessToken("2.00kLXhIG9Cf3tD0015e73a690SCzU5")
//	fmt.Println(s, err)
//	var s, err = wb.GetUserInfo("3749025033", "2.00kLXhIG9Cf3tD0015e73a690SCzU5")
//	fmt.Println(s, err)
}
