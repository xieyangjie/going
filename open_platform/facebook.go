package open_platform

import (
	"github.com/smartwalle/going/http"
	"errors"
	"github.com/smartwalle/going/tools"
)

const (
	kFBCheckAccessTokenURL	= "https://graph.facebook.com/debug_token"
	kFBGeUserInfoURL		= "https://graph.facebook.com/v2.3/"
)

type Facebook struct {
}

func (this *Facebook) checkError(result map[string]interface{}) (error) {
	var rError map[string]interface{}
	if v, ok := result["error"]; ok {
		if d, ok := v.(map[string]interface{}); ok {
			rError = d
		}
	}

	var err error
	if _, ok := rError["code"]; ok {
		value := rError["message"]
		if str, ok := value.(string); ok {
			err = errors.New(str)
		}
	}

	return  err
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

	if err = this.checkError(result); err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if v, ok := result["data"]; ok {
		if d, ok := v.(map[string]interface{}); ok {
			data = d
		}
	}

	var userId = tools.GetString(data["user_id"])
	var expireIn = tools.GetInt64(data["expires_at"])

	token := &AccessToken{}
	token.UserId = userId
	token.ExpireIn = expireIn

	return token, nil
}

func (this *Facebook) GetUserInfo(userId string, accessToken string) (map[string]interface{}, error) {
	var param = map[string]string {
		"access_token": accessToken,
	}

	var url = kFBGeUserInfoURL + userId

	var result, err = http.DoGet(url, param)

	if err != nil {
		return nil, err
	}

	if err = this.checkError(result); err != nil {
		return nil, err
	}

	return result, nil
}