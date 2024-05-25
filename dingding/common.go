package dingding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type DingContent struct {
	Content string `json:"content,omitempty"`
}

type DingAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type DingLink struct {
	Text       string `json:"text,omitempty"`
	Title      string `json:"title,omitempty"`
	PicUrl     string `json:"picUrl,omitempty"`
	MessageUrl string `json:"messageUrl,omitempty"`
}

type DingMarkdown struct {
	Text  string `json:"text,omitempty"`
	Title string `json:"title,omitempty"`
}

type DingBtns struct {
	Title     string `json:"title,omitempty"`
	ActionURL string `json:"actionURL,omitempty"`
}

type DingActionCard struct {
	Text           string     `json:"text,omitempty"`
	Title          string     `json:"title,omitempty"`
	HideAvatar     string     `json:"hideAvatar,omitempty"`
	BtnOrientation string     `json:"btnOrientation,omitempty"`
	SingleTitle    string     `json:"singleTitle,omitempty"`
	SingleURL      string     `json:"singleURL,omitempty"`
	Btns           []DingBtns `json:"btns,omitempty"`
}

type DingFeedCard struct {
	Links []DingLink `json:"links,omitempty"`
}

// 因为你给他具体struct是会有默认值的，这样omitempty就不起作用了，改成指针，默认就会是个nil，此时omiempty起作用
type DingMessage struct {
	MsgType    string          `json:"msgtype"`
	Text       *DingContent    `json:"text,omitempty"`
	At         *DingAt         `json:"at,omitempty"`
	Link       *DingLink       `json:"link,omitempty"`
	Markdown   *DingMarkdown   `json:"markdown,omitempty"`
	ActionCard *DingActionCard `json:"actionCard,omitempty"`
	FeedCard   *DingFeedCard   `json:"feedCard,omitempty"`

	DDURL      string `json:"ddurl,omitempty"`
	IsPrintLog bool   `json:"isPrintLog,omitempty"`
	Message    string `json:"message,omitempty"`
	LogFile    string `json:"logFile,omitempty"`

	IsSend bool `json:"-"`
}

func NewTextDingMsg(content, ddUrl string, atAll bool) *DingMessage {
	return &DingMessage{MsgType: "text", Text: &DingContent{Content: content}, DDURL: ddUrl, At: &DingAt{IsAtAll: atAll}}
}

func NewMDDingMsg(title, text, ddUrl string, at *DingAt) *DingMessage {
	return &DingMessage{MsgType: "markdown", Markdown: &DingMarkdown{
		Text:  text,
		Title: title,
	}, DDURL: ddUrl, At: at}
}

func (m *DingMessage) FormatStr(name, value string) {
	m.Message += fmt.Sprintf("	%s: %s   \n  ", name, value)
	// m.Message +=  fmt.Sprintf(" %-13s: %s \n ", name, value)
}

func (m *DingMessage) FormatFloat(name string, value float64) {
	// m.Message += fmt.Sprintf("%-13s: %.2f%% \n", name, value)
	m.Message += fmt.Sprintf("	%s: %.2f%%   \n  ", name, value)
}

func (m *DingMessage) Round(f float64) float64 {
	n10 := math.Pow10(1)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func (m *DingMessage) FormatEnter() {
	m.Message += fmt.Sprintf("   \n  ")
}

func (m *DingMessage) SetSendFlag() {
	m.IsSend = true
	m.Message += "✘✘☟☟☟☟✘✘"
	m.FormatEnter()
}

func (m *DingMessage) SetAtAll(atAll bool) {
	if m.At == nil {
		m.At = &DingAt{
			IsAtAll: atAll,
		}
	} else {
		m.At.IsAtAll = atAll
	}
}

func (m *DingMessage) SetTextContent(content string) {
	m.Text.Content = content
}

func (m *DingMessage) SetMakeDownText(text string) {
	m.Markdown.Text = text
}

func (m *DingMessage) SetMakeDownTitle(title string) {
	m.Markdown.Title = title
}

func (m *DingMessage) SetMakeDown(title, text string) {
	m.Markdown.Text = text
	m.Markdown.Title = title
}

func (m *DingMessage) CheckDingMessage() error {
	if m.MsgType == "text" {
		// 尝试去查看
		if m.Message != "" {
			m.Text.Content += m.Message
			m.Message = ""
		}
	} else if m.MsgType == "markdown" {
		if m.Message != "" {
			m.Markdown.Text += m.Message
			m.Message = ""
		}
	}

	if m.DDURL == "" {
		return fmt.Errorf("%s", "message DDURL empty !")
	}
	return nil
}

func (m *DingMessage) Clear() {
	m.Message = ""
	m.IsSend = false
	if m.Text != nil {
		m.Text.Content = ""
	}

	if m.Markdown != nil {
		m.Markdown.Text = ""
	}
}

func (m *DingMessage) SendMsg() error {
	if err := m.CheckDingMessage(); err != nil {
		return err
	}
	marshal, err := json.Marshal(*m)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", m.DDURL, bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// {"errmsg":"ok","errcode":0}
	type RequestData struct {
		ErrMsg  string `json:"errmsg"`
		ErrCode int    `json:"errcode"`
	}
	data := RequestData{}
	if err := json.Unmarshal(result, &data); err != nil {
		return err
	}
	if data.ErrMsg == "ok" {
		return nil
	}
	return fmt.Errorf("response: +%v", data)
}

func GetSignUrl(ddUrl, secret string) string {
	if secret == "" {
		return ddUrl
	}

	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	sign := hmacSha256(stringToSign, secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", ddUrl, timestamp, sign)
	return url
}

func hmacSha256(stringToSign, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func SendTexDingMsg(content, ddUrl string, atAll bool) error {
	ding := NewTextDingMsg(content, ddUrl, atAll)
	return ding.SendMsg()
}

func SendMDMsg(title, text, ddUrl string, at *DingAt) error {
	ding := NewMDDingMsg(title, text, ddUrl, at)
	return ding.SendMsg()
}

func (m *DingMessage) GetSetMsg() string {
	marshal, _ := json.Marshal(*m)
	return string(marshal)
}
