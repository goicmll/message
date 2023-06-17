package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// 官方文档: https://developers.dingtalk.com/document/app/dingtalk-openapi-overview
// 包含各种url拼接的函数
type dingTalkURL struct {
	openApiURL  string
	dingTalkURL string
}

var dtu = dingTalkURL{"https://oapi.dingtalk.com", "dingtalk://dingtalkclient"}

// 拼接发送dingtalk 消息的url
func (dtu dingTalkURL) BuildRobotSendMsgURL(accessToken string) string {
	elems := make([]string, 3)
	elems[0] = dtu.openApiURL
	elems[1] = "/robot/send?access_token="
	elems[2] = accessToken
	return strings.Join(elems, "")
}

// 拼接钉钉登录的url
func (dtu dingTalkURL) BuildDingTalkLoginURL(appID, state, callbackURL string) string {
	elems := make([]string, 8)
	elems[0] = dtu.openApiURL
	elems[1] = "/connect/oauth2/sns_authorize?appid="
	elems[2] = appID
	elems[3] = "&response_type=code&scope=snsapi_auth&state="
	elems[4] = state
	elems[5] = "&redirect_uri="
	elems[6] = url.QueryEscape(callbackURL)
	elems[7] = "&container_type=work_platform"
	return strings.Join(elems, "")
}

// 拼接通过用户ID获取钉钉用户详情的URL
func (dtu dingTalkURL) BuildGetUserDetialByUserIDURL(accessToken string) string {
	elems := make([]string, 3)
	elems[0] = dtu.openApiURL
	elems[1] = "/topapi/v2/user/get?access_token="
	elems[2] = accessToken
	return strings.Join(elems, "")
}

// 拼接通过临时code获取用户信息的URL
func (dtu dingTalkURL) BuildGetUserInfoUrlByTempCodeURL(accessToken string) string {
	elems := make([]string, 3)
	elems[0] = dtu.openApiURL
	elems[1] = "/topapi/v2/user/getuserinfo?access_token="
	elems[2] = accessToken
	return strings.Join(elems, "")
}

// 拼接获取 access_token 的URL
func (dtu dingTalkURL) BuildGetAccessTokenURL(appKey, appSecret string) string {
	elems := make([]string, 5)
	elems[0] = dtu.openApiURL
	elems[1] = "/gettoken?appkey="
	elems[2] = appKey
	elems[3] = "&appsecret="
	elems[4] = appSecret
	return strings.Join(elems, "")
}

// 拼接钉钉工作台打开地址的url
func (dtu dingTalkURL) BuildOpenLinkByWorkPlatformURL(corpID, agentID, targetURL string) string {
	elems := make([]string, 7)
	elems[0] = dtu.dingTalkURL
	elems[1] = "/action/openapp?corpid="
	elems[2] = corpID
	elems[3] = "&container_type=work_platform&app_id=0_"
	elems[4] = agentID
	elems[5] = "&redirect_type=jump&redirect_url="
	elems[6] = url.QueryEscape(targetURL)
	return strings.Join(elems, "")
}

// 拼接钉钉侧边栏打开地址的url
func (dtu dingTalkURL) BuildOpenLinkBySlideURL(corpID, agentID, targetURL string) string {
	elems := make([]string, 4)
	elems[0] = dtu.dingTalkURL
	elems[1] = "/page/link?url="
	elems[2] = url.QueryEscape(targetURL)
	elems[3] = "&pc_slide=true"
	return strings.Join(elems, "")
}

// 钉钉消息相关
type DingTalkMsg struct {
	Msgtype  string              `json:"msgtype"`
	Markdown DingTalkMarkdownMsg `json:"markdown,omitempty"`
	At       DingTalkAt          `json:"at,omitempty"`
}

type DingTalkMarkdownMsg struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

type DingTalkAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type DingTalkRobot struct {
	AccessToken string
}

// 发送消息
// 发送markdown消息
// """
//
//	标题
//	# 一级标题
//	## 二级标题
//	### 三级标题
//	#### 四级标题
//	##### 五级标题
//	###### 六级标题
//
//	引用
//	> A man who stands for nothing will fall for anything.
//
//	文字加粗、斜体
//	**bold**
//	*italic*
//
//	链接
//	[this is a link](http://name.com)
//
//	图片
//	![](http://name.com/pic.jpg)
//
//	无序列表
//	- item1
//	- item2
//
//	有序列表
//	1. item1
//	2. item2
//	"""
func (dtr DingTalkRobot) SendMsg(msg DingTalkMsg) error {

	jsonBytes, _ := json.Marshal(msg)
	resp, err := http.Post(dtu.BuildRobotSendMsgURL(dtr.AccessToken), "application/json", bytes.NewBuffer(jsonBytes))
	defer closeRespBody(resp)
	if err != nil {
		return NewMessageError(err.Error())
	}
	return nil
}

// 钉钉openapi相关
type DingTalkClient struct {
	AppKey    string
	AppSecret string
	AgentID   string
	CorpID    string
}

type dtBaseResp struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type dingTalkAccessTokenResp struct {
	dtBaseResp
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type dingTalkUserInfo struct {
	AssociatedUnionID string `json:"associated_unionid"`
	UnionID           string `json:"unionid"`
	DeviceID          string `json:"device_id"`
	SysLevel          int    `json:"sys_level"`
	Name              string `json:"name"`
	Sys               bool   `json:"sys"`
	UserID            string `json:"userid"`
}
type dingTalkUserInfoResp struct {
	dtBaseResp
	Result dingTalkUserInfo `json:"result"`
}

type dingTalkUserDetail struct {
	UnionID string `json:"unionid"`
	UserID  string `json:"userid"`
	Email   string `json:"email"`
	Mobile  string `json:"mobile"`
	Active  bool   `json:"active"`
	Remark  string `json:"remark"`
	Name    string `json:"name"`
}

type dingTalkUserDetailResp struct {
	dtBaseResp
	Result dingTalkUserDetail `json:"result"`
}

// open api access_token 缓存
type DingTalkAccessTokenCache struct {
	C *cache.Cache
}

var c = DingTalkAccessTokenCache{cache.New(7000*time.Second, 7200*time.Second)}

// 获取token, 返回 "" 说明token已失效
func (d DingTalkAccessTokenCache) Get(appKey string) string {
	if v, ok := d.C.Get(appKey); ok {
		if v1, yes := v.(string); yes {
			return v1
		} else {
			return ""
		}
	} else {
		return ""
	}
}

// 添加或更新access token
func (d DingTalkAccessTokenCache) Update(appKey, token string) {
	d.C.Set(appKey, token, cache.DefaultExpiration)
}

// 获取本地的token缓存, 如果没有则请求更新, 并更新缓存, 返回新token
func (dtc *DingTalkClient) GetAccessTokenFromCache() (string, error) {
	if ac := c.Get(dtc.AppKey); ac != "" {
		return ac, nil
	}
	ac, err := dtc.GetAccessToken()
	if err != nil {
		return "", NewMessageError(err.Error())
	}
	return ac.AccessToken, nil
}

// 获取钉钉access_token
func (dtc *DingTalkClient) GetAccessToken() (*dingTalkAccessTokenResp, error) {
	if dtc.AppKey == "" || dtc.AppSecret == "" {
		return nil, NewMessageError("不可用的AppKey或AppSecret!")
	}
	req, err := http.NewRequest("GET", dtu.BuildGetAccessTokenURL(dtc.AppKey, dtc.AppSecret), http.NoBody)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	resp, err := (&http.Client{Timeout: 3 * time.Second}).Do(req)
	defer closeRespBody(resp)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	r := dingTalkAccessTokenResp{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, NewMessageError("json解析钉钉接口响应错误")
	}
	if r.ErrCode == 0 {
		c.Update(dtc.AppKey, r.AccessToken)
		return &r, nil
	} else {
		return nil, NewMessageError(r.ErrMsg)
	}
}

// 通过临时code获取用户信息, 包括userid unionid
func (dtc *DingTalkClient) GetUserInfoByTempCode(tempCode string) (*dingTalkUserInfoResp, error) {
	ac, err := dtc.GetAccessTokenFromCache()
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	reqBody := strings.NewReader(fmt.Sprintf(`{"code": "%s"}`, tempCode))
	req, err := http.NewRequest("POST", dtu.BuildGetUserInfoUrlByTempCodeURL(ac), reqBody)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	resp, err := (&http.Client{Timeout: 3 * time.Second}).Do(req)
	defer closeRespBody(resp)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	r := dingTalkUserInfoResp{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, NewMessageError("json解析钉钉接口响应错误")
	}
	if r.ErrCode == 0 {
		return &r, nil
	} else {
		return nil, NewMessageError(r.ErrMsg)
	}
}

// 通过 tempCode 获取用户详情
func (dtc *DingTalkClient) GetUserDetailByTemCode(tempCode string) (*dingTalkUserDetailResp, error) {
	ac, err := dtc.GetAccessTokenFromCache()
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	userInfo, err := dtc.GetUserInfoByTempCode(tempCode)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	reqBody := strings.NewReader(fmt.Sprintf(`{"userid": "%s", "language": "zh_CN"}`, userInfo.Result.UserID))
	req, err := http.NewRequest("POST", dtu.BuildGetUserDetialByUserIDURL(ac), reqBody)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	resp, err := (&http.Client{Timeout: 3 * time.Second}).Do(req)
	defer closeRespBody(resp)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewMessageError(err.Error())
	}
	r := dingTalkUserDetailResp{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, NewMessageError("json解析钉钉接口响应错误")
	}
	if r.ErrCode == 0 {
		return &r, nil
	} else {
		return nil, NewMessageError(r.ErrMsg)
	}
}
