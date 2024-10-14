package service

import (
	"encoding/json"
	"mstuca_schedule/internal/botErrors"
	"mstuca_schedule/internal/models"
	"net/http"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type ScheduleGetter interface {
	GetSchedule(filter *models.Filter) (*models.Schedule, error)
	GetGroupID(groupName string) ([]*models.Group, error)
}

type scheduleGetter struct {
	client *http.Client
	json   jsoniter.API
}

func New() ScheduleGetter {
	return &scheduleGetter{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
		json: jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			UseNumber:              true,
		}.Froze(),
	}
}

func (sg *scheduleGetter) GetSchedule(filter *models.Filter) (*models.Schedule, error) {
	return &models.Schedule{}, nil
}

func (sg *scheduleGetter) GetGroupID(groupName string) ([]*models.Group, error) {
	request, err := http.NewRequest(http.MethodGet, "https://ruz.mstuca.ru/api/search", nil)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("term", groupName)
	values.Add("type", "group")

	request.URL.RawQuery = values.Encode()

	response, err := sg.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, err
	}

	groups := make([]*models.Group, 0)

	if err := json.NewDecoder(response.Body).Decode(&groups); err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, botErrors.ErrNoGroupsFound
	}

	return groups, nil
}
