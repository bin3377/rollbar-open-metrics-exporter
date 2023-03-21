package rollbar

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var BaseURL = "https://api.rollbar.com/api/1"

var AccountReadAccessToken = ""
var AccountWriteAccessToken = ""

type listProjectsResponse struct {
	Err    int       `json:"err"`
	Result []Project `json:"result"`
}

func ListProjects() ([]Project, error) {
	var resp listProjectsResponse
	if err := jcall(
		"GET",
		AccountReadAccessToken,
		fmt.Sprintf("%s/projects", BaseURL),
		nil,
		&resp); err != nil {
		return nil, err
	}
	if resp.Err != 0 {
		return nil, fmt.Errorf("rollbar returns error code %d", resp.Err)
	}
	return resp.Result, nil
}

type listProjectAccessTokensResponse struct {
	Err    int                  `json:"err"`
	Result []ProjectAccessToken `json:"result"`
}

func ListProjectAccessTokens(projectID int) ([]ProjectAccessToken, error) {
	var resp listProjectAccessTokensResponse
	if err := jcall(
		"GET",
		AccountReadAccessToken,
		fmt.Sprintf("%s/project/%d/access_tokens", BaseURL, projectID),
		nil,
		&resp); err != nil {
		return nil, err
	}
	if resp.Err != 0 {
		return nil, fmt.Errorf("rollbar returns error code %d", resp.Err)
	}
	return resp.Result, nil
}

type createProjectAccessTokenResponse struct {
	Err    int                `json:"err"`
	Result ProjectAccessToken `json:"result"`
}

func CreateProjectAccessToken(projectID int, params CreateProjectAccessTokenParams) (*ProjectAccessToken, error) {
	var resp createProjectAccessTokenResponse
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	if err := jcall(
		"POST",
		AccountWriteAccessToken,
		fmt.Sprintf("%s/project/%d/access_tokens", BaseURL, projectID),
		payload,
		&resp); err != nil {
		return nil, err
	}
	if resp.Err != 0 {
		return nil, fmt.Errorf("rollbar returns error code %d", resp.Err)
	}
	return &resp.Result, nil
}

var ErrReadTokenNotFound = errors.New("read token is not found")

func GetProjectReadToken(projectID int) (*ProjectAccessToken, error) {
	tokens, err := ListProjectAccessTokens(projectID)
	if err != nil {
		return nil, err
	}
	for _, token := range tokens {
		if token.Status == StatusDisabled {
			continue
		}
		for _, scope := range token.Scopes {
			if scope == ScopeRead {
				return &token, nil
			}
		}
	}
	return nil, ErrReadTokenNotFound
}

func GetOrCreateProjectReadToken(projectID int) (*ProjectAccessToken, error) {
	token, err := GetProjectReadToken(projectID)
	if err != nil {
		if err == ErrReadTokenNotFound {
			logrus.Debugf("read token of project %d is not found, creating one...", projectID)
			return CreateProjectAccessToken(projectID, CreateProjectAccessTokenParams{
				Name:   "read",
				Scopes: []Scope{ScopeRead},
				Status: StatusEnabled,
			})
		}
	}
	return token, err
}

type listEnvironmentsResult struct {
	Err    int `json:"err"`
	Result struct {
		Environments []Environment `json:"environments"`
		Page         int           `json:"page"`
		Limit        int           `json:"limit"`
	} `json:"result"`
}

func ListEnvrionments(projectToken string) ([]Environment, error) {
	page := 1
	limit := 5000
	var result []Environment
	var resp listEnvironmentsResult
	for {
		if err := jcall(
			"GET",
			projectToken,
			fmt.Sprintf("%s/environments?page=%d&limit=%d", BaseURL, page, limit),
			nil,
			&resp); err != nil {
			return nil, err
		}
		if resp.Err != 0 {
			return nil, fmt.Errorf("rollbar returns error code %d", resp.Err)
		}
		result = append(result, resp.Result.Environments...)
		if len(resp.Result.Environments) < limit {
			break
		}
	}
	return result, nil
}

type getItemByIDResponse struct {
	Err    int  `json:"err"`
	Result Item `json:"result"`
}

func GetItemByID(projectToken string, id int) (*Item, error) {
	var resp getItemByIDResponse
	if err := jcall(
		"GET",
		projectToken,
		fmt.Sprintf("%s/item/%d", BaseURL, id),
		nil,
		&resp); err != nil {
		return nil, err
	}
	if resp.Err != 0 {
		return nil, fmt.Errorf("rollbar returns error code %d", resp.Err)
	}
	return &resp.Result, nil
}

type listItemsWithIDsResponse struct {
	Err    int `json:"err"`
	Result struct {
		Items []Item `json:"items"`
		Page  int    `json:"page"`
		Limit int    `json:"limit"`
	} `json:"result"`
}

func ListItemsWithIDs(projectToken string, ids []int) ([]Item, error) {
	var resp listItemsWithIDsResponse
	strIDs := ""
	for _, id := range ids {
		strIDs += fmt.Sprint(id) + ","
	}
	if err := jcall(
		"GET",
		projectToken,
		fmt.Sprintf("%s/items?ids=%s", BaseURL, strIDs),
		nil,
		&resp); err != nil {
		return nil, err
	}
	if resp.Err != 0 {
		return nil, fmt.Errorf("rollbar returns error code %d", resp.Err)
	}
	return resp.Result.Items, nil
}

type getOccurencesMetricsResponse struct {
	Err    int                    `json:"err"`
	Result OccurenceMetricsResult `json:"result"`
}

func GetOccurrencesMetrics(projectToken string, params OccurrenceMetricsParams) (*OccurenceMetricsResult, error) {
	var resp getOccurencesMetricsResponse
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	if err := jcall(
		"POST",
		projectToken,
		fmt.Sprintf("%s/metrics/occurrences", BaseURL),
		payload,
		&resp); err != nil {
		return nil, err
	}
	if resp.Err != 0 {
		return nil, fmt.Errorf("rollbar returns error code %d", resp.Err)
	}
	return &resp.Result, nil
}

func NewItemOccurrencesInput(ago time.Duration, limit int) OccurrenceMetricsParams {
	end := time.Now()
	start := end.Add(-ago)
	return OccurrenceMetricsParams{
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		GroupBy: []Field{
			FieldItemId,
		},
		Limit: limit,
	}
}

func NewItemOccurrencesFullInput(ago time.Duration, limit int) OccurrenceMetricsParams {
	end := time.Now()
	start := end.Add(-ago)
	return OccurrenceMetricsParams{
		StartTime: start.Unix(),
		EndTime:   end.Unix(),
		GroupBy: []Field{
			FieldItemId,
			FieldEnvironment,
			FieldItemTitle,
			FieldItemStatus,
			FieldItemLevel,
		},
		Limit: limit,
	}
}

func GetItemOccurrences(projectToken string, ago time.Duration, limit int) ([]ItemOccurrence, error) {
	metrics, err := GetOccurrencesMetrics(projectToken, NewItemOccurrencesFullInput(ago, limit))
	if err != nil {
		return nil, err
	}
	conv := func(v any) int64 {
		i, err := v.(json.Number).Int64()
		if err != nil {
			logrus.Errorf("%v is not int64", v)
		}
		return i
	}
	result := make([]ItemOccurrence, 0)
	for _, tp := range metrics.Timepoints {
		for _, row := range tp.MetricsRows {
			logrus.Debugf("%v", row)
			single := ItemOccurrence{
				Time: time.Unix(tp.Timestamp, 0),
			}
			for _, cell := range row {
				switch cell.Field {
				case FieldItemId:
					single.ItemID = int(conv(cell.Value))
				case FieldOccurrenceCount:
					single.OccurrenceCount = conv(cell.Value)
				case FieldEnvironment:
					single.Environment = cell.Value.(string)
				case FieldItemTitle:
					single.ItemTitle = cell.Value.(string)
				case FieldItemStatus:
					single.ItemStatus = cell.Value.(string)
				case FieldItemLevel:
					single.ItemLevel = cell.Value.(string)
				}
			}
			result = append(result, single)
		}
	}
	return result, nil
}
