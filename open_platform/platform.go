package open_platform


type IPlatform interface {

	//验证 access token 是否有效
	CheckAccessToken(accessToken string) (*AccessToken, error)

	//获取用户信息
	GetUserInfo(userId string, accessToken string) (map[string]interface{}, error)
}



