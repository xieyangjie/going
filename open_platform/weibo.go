package open_platform

import (
	"github.com/smartwalle/going/http"
	"errors"
	"github.com/smartwalle/going/tools"
)

const (
	kWBCheckAccessTokenURL 	= "https://api.weibo.com/oauth2/get_token_info"
	kWBGetUserInfoURL		= "https://api.weibo.com/2/users/show.json"
)

type Weibo struct {
}

func (this *Weibo) CheckError(result map[string]interface{}) (error) {
	var err error
	if value, ok := result["error"]; ok {
		if str, ok := value.(string); ok {
			err = errors.New(str)
		}
		return err
	}
	return err
}

func (this *Weibo) CheckAccessToken(accessToken string) (*AccessToken, error) {
	var param = map[string]string{
		"access_token": accessToken,
	}

	var result, err = http.DoPost(kWBCheckAccessTokenURL, param)

	if err != nil {
		return nil, err
	}

	if err := this.CheckError(result); err != nil {
		return nil, err
	}

	var userId = tools.GetString(result["uid"])

	var expireIn float64 = tools.GetFloat64(result["create_at"])
	var createdAt float64 = tools.GetFloat64(result["expire_in"])

	token := &AccessToken{}
	token.UserId = userId
	token.ExpireIn = int64(expireIn + createdAt)

	return token, nil
}

func (this *Weibo) GetUserInfo(userId string, accessToken string) (map[string]interface{}, error) {
	var param = map[string]string{
		"access_token": accessToken,
		"uid": userId,
	}

	var result, err = http.DoGet(kWBGetUserInfoURL, param)

	if err != nil {
		return nil, err
	}

	if err := this.CheckError(result); err != nil {
		return nil, err
	}

	return result, nil
}