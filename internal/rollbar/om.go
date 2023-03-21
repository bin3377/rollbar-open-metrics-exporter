package rollbar

import "time"

// Status represents the enabled or disabled status of an entity.
type Status string

// Possible values for status
const (
	StatusEnabled  = Status("enabled")
	StatusDisabled = Status("disabled")
)

type Project struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	AccountID    int    `json:"account_id"`
	DateCreated  int    `json:"date_created"`
	DateModified int    `json:"date_modified"`
	Status       Status `json:"status"`
}

type Scope string

// Possible values for project access token scope
const (
	ScopeWrite          = Scope("write")
	ScopeRead           = Scope("read")
	ScopePostServerItem = Scope("post_server_item")
	ScopePostClientItem = Scope("post_client_item")
)

type ProjectAccessToken struct {
	Name                    string  `json:"name"`
	ProjectID               int     `json:"project_id"`
	AccessToken             string  `json:"access_token"`
	Scopes                  []Scope `json:"scopes"`
	Status                  Status  `json:"status"`
	RateLimitWindowSize     int     `json:"rate_limit_window_size"`
	RateLimitWindowCount    int     `json:"rate_limit_window_count"`
	DateCreated             int     `json:"date_created"`
	DateModified            int     `json:"date_modified"`
	CurRateLimitWindowCount int     `json:"cur_rate_limit_window_count"`
	CurRateLimitWindowStart int     `json:"cur_rate_limit_window_start"`
}

type CreateProjectAccessTokenParams struct {
	Name                 string  `json:"name"`
	Scopes               []Scope `json:"scopes"`
	Status               Status  `json:"status"`
	RateLimitWindowSize  int     `json:"rate_limit_window_size,omitempty"`
	RateLimitWindowCount int     `json:"rate_limit_window_count,omitempty"`
}

type Environment struct {
	ID          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	Environment string `json:"environment"`
	Visible     int    `json:"visible"`
}

// Field represents the allowed Field name options of occurrences API.
type Field string

const (
	FieldProjectId          = Field("project_id")
	FieldItemId             = Field("item_id")
	FieldEnvironment        = Field("environment")
	FieldBrowserFamily      = Field("browser_family")
	FieldBrowserVersion     = Field("browser_version")
	FieldOsFamily           = Field("os_family")
	FieldOsVersion          = Field("os_version")
	FieldDeviceBrand        = Field("device_brand")
	FieldDeviceModel        = Field("device_model")
	FieldIpAddress          = Field("ip_address")
	FieldItemStatus         = Field("item_status")
	FieldItemLevel          = Field("item_level")
	FieldItemGroupItemId    = Field("item_group_item_id")
	FieldItemTitle          = Field("item_title")
	FieldItemCounter        = Field("item_counter")
	FieldPersonUsername     = Field("person_username")
	FieldPersonEmail        = Field("person_email")
	FieldPersonId           = Field("person_id")
	FieldCodeVersion        = Field("code_version")
	FieldCount              = Field("count")
	FieldOccurrenceId       = Field("occurrence_id")
	FieldUuid               = Field("uuid")
	FieldContext            = Field("context")
	FieldPlatform           = Field("platform")
	FieldFramework          = Field("framework")
	FieldPlatformCanonical  = Field("platform_canonical")
	FieldFrameworkCanonical = Field("framework_canonical")
	FieldLanguage           = Field("language")
	FieldLanguageName       = Field("language_name")
	FieldNotifierName       = Field("notifier_name")
	FieldNotifierVersion    = Field("notifier_version")
	FieldOccurrenceCount    = Field("occurrence_count")
	FieldMessageBody        = Field("message_body")
	FieldTimestamp          = Field("timestamp")
	FieldFingerprint        = Field("fingerprint")
	FieldServerHost         = Field("server_host")
	FieldServerRoot         = Field("server_root")
	FieldServerPid          = Field("server_pid")
	FieldServerCpu          = Field("server_cpu")
	FieldScmBranch          = Field("scm_branch")
	FieldRequestUrl         = Field("request_url")
	FieldRequestMethod      = Field("request_method")
	FieldRequestQueryString = Field("request_query_string")
	FieldRequestBody        = Field("request_body")
)

// AggregateFunction represents the allowed functions in Aggregate
type AggregateFunction string

const (
	AggregateFunctionCountAll      = AggregateFunction("count_all")
	AggregateFunctionCountDistinct = AggregateFunction("count_distinct")
	AggregateFunctionMax           = AggregateFunction("max")
	AggregateFunctionMin           = AggregateFunction("min")
)

type Aggregate struct {
	Field    Field             `json:"fileld"`
	Function AggregateFunction `json:"function"`
	Alias    string            `json:"alias"`
}

// FilterOperator represents the allowed operators in Filter
type FilterOperator string

const (
	FilterOperatorEq         = FilterOperator("eq")
	FilterOperatorNe         = FilterOperator("ne")
	FilterOperatorGt         = FilterOperator("gt")
	FilterOperatorGte        = FilterOperator("gte")
	FilterOperatorLt         = FilterOperator("lt")
	FilterOperatorLte        = FilterOperator("lte")
	FilterOperatorNotLike    = FilterOperator("not_like")
	FilterOperatorBetween    = FilterOperator("between")
	FilterOperatorNotBetween = FilterOperator("not_between")
)

type Filter struct {
	Field    Field          `json:"fileld"`
	Values   []string       `json:"values"`
	Operator FilterOperator `json:"operator"`
}

type Granularity string

const (
	GranularitySecond = Granularity("second")
	GranularityMinute = Granularity("minute")
	GranularityHour   = Granularity("hour")
	GranularityDay    = Granularity("day")
	GranularityWeek   = Granularity("week")
	GranularityMonth  = Granularity("month")
	GranularityYear   = Granularity("year")
)

type Order string

const (
	OrderAsc  = Order("asc")
	OrderDesc = Order("desc")
)

type Sort struct {
	Order Order `json:"order"`
	Field Field `json:"field"`
}

type OccurrenceMetricsParams struct {
	StartTime   int64        `json:"start_time"`
	EndTime     int64        `json:"end_time"`
	Filters     []Filter     `json:"filters,omitempty"`
	GroupBy     []Field      `json:"group_by"`
	Aggregates  []Aggregate  `json:"aggregates,omitempty"`
	Sort        *Sort        `json:"sort,omitempty"`
	Granularity *Granularity `json:"granularity,omitempty"`
}

type FieldValue struct {
	Field Field `json:"field"`
	Value any   `json:"value"`
}

type MetricsRows [][]FieldValue

type TimePoint struct {
	Timestamp   int64       `json:"timestamp"`
	MetricsRows MetricsRows `json:"metrics_rows"`
}

type OccurenceMetricsResult struct {
	LastOccurrenceTimestamp int64       `json:"last_occurrence_timestamp"`
	QueryExecution          float64     `json:"query_execution"`
	Timepoints              []TimePoint `json:"timepoints"`
}

type ItemOccurrence struct {
	Time            time.Time
	ItemID          int
	Environment     string
	ItemTitle       string
	ItemStatus      string
	ItemLevel       string
	OccurrenceCount int64
}

type Item struct {
	ID                       int    `json:"id"`
	ProjectID                int    `json:"project_id"`
	CounterID                int    `json:"counter"`
	Environment              string `json:"environment"`
	Platform                 string `json:"platform"`
	Framework                string `json:"framework"`
	Hash                     string `json:"hash"`
	Title                    string `json:"title"`
	Status                   string `json:"status"`
	Level                    string `json:"level"`
	FirstOccurrenceId        int    `json:"first_occurrence_id"`
	FirstOccurrenceTimestamp int    `json:"first_occurrence_timestamp"`
	LastOccurrenceId         int    `json:"last_occurrence_id"`
	LastOccurrenceTimestamp  int    `json:"last_occurrence_timestamp"`
	TotalOccurrences         int64  `json:"total_occurrences"`
}
