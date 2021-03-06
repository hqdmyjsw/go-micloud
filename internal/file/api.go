package file

import (
	"errors"
	"github.com/tidwall/gjson"
	"go-micloud/internal/user"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	BaseUri      = "https://i.mi.com"
	GetFiles     = BaseUri + "/drive/user/files/%s?jsonpCallback=callback"
	CreateFile   = BaseUri + "/drive/user/files/create"
	UploadFile   = BaseUri + "/drive/user/files"
	DeleteFiles  = BaseUri + "/drive/v2/user/records/filemanager"
	GetFolders   = BaseUri + "/drive/user/folders/%s/children"
	CreateFolder = BaseUri + "/drive/v2/user/folders/create"
)

type Api struct {
	User *user.User
}

func NewApi(user *user.User) *Api {
	api := Api{
		User: user,
	}
	return &api
}

func (api *Api) get(url string) ([]byte, error) {
	result, err := api.User.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if result.StatusCode == http.StatusFound {
		return api.get(result.Header.Get("Location"))
	}
	if result.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("登录授权失败")
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	if gjson.Get(string(body), "R").Int() == 401 {
		return nil, errors.New("登录授权失败")
	}
	return body, nil
}

func (api *Api) postForm(url string, values url.Values) (*[]byte, error) {
	result, err := api.User.HttpClient.PostForm(url, values)
	if err != nil {
		return nil, err
	}
	if result.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("登录授权失败")
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	if gjson.Get(string(body), "R").Int() == 401 {
		return nil, errors.New("登录授权失败")
	}
	return &body, nil
}
