package open_platform

import (
	"github.com/smartwalle/going/http"
	"errors"
	"fmt"
)

const (
	kFBCheckAccessTokenURL = "https://graph.facebook.com/debug_token"
)

type Facebook struct {
}

func (this *Facebook) CheckAccessToken(accessToken string) (*AccessToken, error) {
	var param = map[string]string{
		"access_token": accessToken,
		"input_token": accessToken,
	}

	var result, err = http.DoGet(kFBCheckAccessTokenURL, param)

	if err != nil {
		return nil, err
	}

	var rError map[string]interface{}
	if v, ok := result["error"]; ok {
		if d, ok := v.(map[string]interface{}); ok {
			rError = d
		}
	}
	if _, ok := rError["code"]; ok {
		value := rError["message"]
		if str, ok := value.(string); ok {
			err = errors.New(str)
		}
		return nil, err
	}

	var data map[string]interface{}
	if v, ok := result["data"]; ok {
		if d, ok := v.(map[string]interface{}); ok {
			data = d
		}
	}

	fmt.Println(result)

	var userId = ""
	if v, ok := data["user_id"]; ok {
		if str, ok := v.(string); ok {
			userId = str
		}
	}

	var expireIn float64
	if str, ok := data["expires_at"].(float64); ok {
		expireIn = str
	}

	token := &AccessToken{}
	token.UserId = userId
	token.ExpireIn = int64(expireIn)

	return token, nil
}