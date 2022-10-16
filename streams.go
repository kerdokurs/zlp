package zlp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Stream struct {
	Id          int    `json:"stream_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	InviteOnly  bool   `json:"invite_only"`
	// CreatedAt   time.Time `json:"date_created"`
}

func (b *Bot) GetSubscribedStreams() ([]Stream, error) {
	body, err := b.getResponseData("GET", "users/me/subscriptions", nil)

	type response struct {
		Result        string   `json:"result"`
		Msg           string   `json:"msg"`
		Subscriptions []Stream `json:"subscriptions"`
	}

	var res response
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res.Subscriptions, nil
}

func (b *Bot) GetStreams() ([]Stream, error) {
	body, err := b.getResponseData("GET", "streams", nil)
	if err != nil {
		return nil, err
	}

	type response struct {
		Msg     string   `json:"msg"`
		Streams []Stream `json:"streams"`
		Result  string   `json:"result"`
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return res.Streams, nil
}

func (b *Bot) GetStream(id int) (*Stream, error) {
	req, err := b.constructRequest("GET", fmt.Sprintf("streams/%d", id), nil)
	if err != nil {
		return nil, err
	}
	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type response struct {
		Response
		Msg    string `json:"msg"`
		Stream Stream `json:"stream"`
		Result string `json:"result"`
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	stream := res.Stream

	return &stream, nil
}

func (b *Bot) CreateStream(name, description string, private bool, principals []string) (*Stream, error) {
	type subscription struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	sub := subscription{
		Name:        name,
		Description: description,
	}

	sw := bytes.NewBufferString("")
	if err := json.NewEncoder(sw).Encode(&sub); err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("subscriptions", fmt.Sprintf("[%s]", sw.String()))
	values.Set("invite_only", fmt.Sprintf("%t", private))
	values.Set("principals", fmt.Sprintf("%v", principals))

	_, err := b.getResponseData("POST", "users/me/subscriptions", &values)
	if err != nil {
		return nil, err
	}

	// type response struct {
	// 	Response

	// 	Subscribed        map[string][]string `json:"subscribed"`
	// 	AlreadySubscribed map[string][]string `json:"already_subscribed"`
	// 	Unauthorized      bool                `json:"unauthorized"`
	// }

	return nil, nil
}
