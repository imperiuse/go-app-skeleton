package apihelper

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/google/uuid"

	"github.com/imperiuse/go-app-skeleton/internal/servers/api/controller/apierror"

	"github.com/gin-gonic/gin"
)

var ErrCouldNotConvertToEnum = errors.New("could not convert to specific enum type")

// ConvertToEnumValueOrDefault - convert given s to T if s value in enumDict, or get default T value.
func ConvertToEnumValueOrDefault[T ~string](s string, enumDict []T, defaultIfNotFoundAny T) T {
	v, err := ConvertToEnumValue(s, enumDict)
	if err != nil {
		return defaultIfNotFoundAny
	}

	return v
}

// ConvertToEnumValue - convert given s to T if s value in enumDict, or return error.
func ConvertToEnumValue[T ~string](s string, enumDict []T) (T, error) {
	s = strings.ToLower(s)
	for _, v := range enumDict {
		if strings.EqualFold(s, string(v)) {
			return v, nil
		}
	}

	return *new(T), ErrCouldNotConvertToEnum
}

// ConvertToEnumsSliceWithDefault - convert slice string, to slice T, if values s, in enumDict, if result [] empty
// return []T{defaultValueT}.
func ConvertToEnumsSliceWithDefault[T ~string](unknownStrings []string, enumDict []T, defaultIfNotFoundAny T) []T {
	if len(unknownStrings) == 0 {
		return []T{defaultIfNotFoundAny}
	}

	result := make([]T, 0, len(unknownStrings))
	for _, s := range unknownStrings {
		for _, v := range enumDict {
			if strings.EqualFold(s, string(v)) {
				result = append(result, v)
			}
		}
	}

	if len(result) == 0 {
		return []T{defaultIfNotFoundAny}
	}

	return result
}

// ConvertToEnumsSlice - convert slice string, to slice T, if values s, in enumDict, if result [] empty return empty []T.
func ConvertToEnumsSlice[T ~string](unknownStrings []string, enumDict []T) []T {
	if len(unknownStrings) == 0 {
		return []T{}
	}

	result := make([]T, 0, len(unknownStrings))
	for _, s := range unknownStrings {
		for _, v := range enumDict {
			if strings.EqualFold(s, string(v)) {
				result = append(result, v)
			}
		}
	}

	if len(result) == 0 {
		return []T{}
	}

	return result
}

// GetIntFromStr - safety convert string to int.
func GetIntFromStr(s string, defaultValue int, minValue int, maxValue int) int {
	if parsedResult, err := strconv.Atoi(s); err == nil && parsedResult >= minValue && parsedResult <= maxValue {
		return parsedResult
	}

	return defaultValue
}

// GetTestGinContext - get mock gin context.
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode) // NB! important for right fill gin ctx!
	// more info here -> https://blog.canopas.com/golang-unit-tests-with-test-gin-context-80e1ac04adcd.

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

// //nolint: gosec, - it's naming only!
const authTokenHeader = "X-AUTH-TOKEN"

func SetAPITokenToGinCtx[T any](c *gin.Context, token T) {
	c.Set(authTokenHeader, token)
}

func GetAPIToneFromGinCtx[T any](c *gin.Context) (T, error) {
	tokenS, exist := c.Get(authTokenHeader)
	if !exist {
		return *new(T), errors.New("{authTokenHeader} not exists, check auth middleware")
	}

	token, ok := tokenS.(T)
	if !ok {
		return *new(T), errors.New("could not convert {authTokenHeader} string to AuthToken")
	}

	return token, nil
}

type (
	QueryParamName = string

	AllowForSearchField string

	sortField string
	sortOrder string
)

const (
	IndexNameParam QueryParamName = "index"     // redefine index name.
	TenantIDParam  QueryParamName = "tenant_id" // redefine index name.

	fromDateParam      QueryParamName = "from_date" // from created_at date time.Format(time.RFC3339)
	toDateParam        QueryParamName = "to_date"   // to created_at date time.Format(time.RFC3339)
	searchPatternParam QueryParamName = "search"    // search pattern in complex endpoint for multi-field searching.
	cursorsParam       QueryParamName = "cursor"    // comma separated array of string values

	pageParam    QueryParamName = "page"     // >= 0.
	limitParam   QueryParamName = "limit"    // >0 and <= 10 000.
	sortByParam  QueryParamName = "sort_by"  // name_of_field order by.
	orderByParam QueryParamName = "order_by" // asc, desc.
)

const (
	none      sortField = ""
	CreatedAt sortField = "created_at"
)

const (
	Asc  sortOrder = "asc"
	Desc sortOrder = "desc"
)

const (
	defaultPageValue       = 0
	defaultLimitParamValue = 10
	maxPageParamValue      = math.MaxInt
	MaxLimitParamValue     = 10000
)

var (
	allSortByFields      = [...]sortField{CreatedAt} // MUST BE NOT EMPTY! @see ParseSortByAndOrderByParams func.
	allOrderByDirections = [...]sortOrder{Asc, Desc} // MUST BE NOT EMPTY! @see ParseSortByAndOrderByParams func.
)

const (
	IDField                 AllowForSearchField = "id"
	TenantIDField           AllowForSearchField = "tenant_id"
	NestedSubjectEmailField AllowForSearchField = "subject.data.email"
)

func ParseSortByAndOrderByParams(
	c *gin.Context,
) (sortBy []string, orderBy []string) {
	listOfSortBy := ConvertToEnumsSliceWithDefault(
		c.QueryArray(sortByParam),
		allSortByFields[:],
		allSortByFields[0],
	)

	listOfOrderBy := ConvertToEnumsSliceWithDefault(
		c.QueryArray(orderByParam),
		allOrderByDirections[:],
		allOrderByDirections[0],
	)

	if len(listOfSortBy) > 0 && strings.EqualFold(string(listOfSortBy[0]), string(none)) ||
		len(listOfSortBy) != len(listOfOrderBy) {
		return []string{string(allSortByFields[0])}, []string{string(allOrderByDirections[0])}
	}

	return FastUnsafeConvertToStringSlice(listOfSortBy), FastUnsafeConvertToStringSlice(listOfOrderBy)
}

func ParsePageAndLimitPaginationOptionsByParams(c *gin.Context) (page int, limit int) {
	page = GetIntFromStr(c.Query(pageParam), defaultPageValue, defaultPageValue, maxPageParamValue)
	limit = GetIntFromStr(c.Query(limitParam), defaultLimitParamValue, 1, MaxLimitParamValue)

	return page, limit
}

func FastUnsafeConvertToStringSlice[T ~string](ss []T) []string {
	// It's more cheap rather than :
	// type T string
	// ss := []T{"hello", "world"}
	//
	// s := make([]string, len(ss))
	// for i, v := range ss {
	// 	s[i] = string(v)
	// }
	// return s

	return *(*[]string)(unsafe.Pointer(&ss))
}

// ParseQueryParamsForComplexSearch - ParseQueryParamsForComplexSearch.
// //nolint: gocritic,  - here is ok.
func ParseQueryParamsForComplexSearch(c *gin.Context, endpointPath string) (
	tenantID string,
	index string,
	fromTime time.Time,
	toTime time.Time,
	searchPattern string,
	successParse bool,
) {
	index, tenantID, successParse = ParseQueryIndexAndTenantParamsOnly(c, endpointPath)
	if !successParse {
		return index, tenantID, fromTime, toTime, searchPattern, successParse
	}

	fromTime, toTime, successParse = ParseQueryDateParamsOnly(c, endpointPath)
	if !successParse {
		return index, tenantID, fromTime, toTime, searchPattern, successParse
	}

	searchPattern, successParse = ParseQuerySearchPatternParamOnly(c, endpointPath)

	return index, tenantID, fromTime, toTime, searchPattern, successParse
}

func ParseQueryIndexAndTenantParamsOnly(c *gin.Context, endpointPath string) (
	tenantID string,
	index string,
	successParse bool,
) {
	tenantID = strings.TrimSpace(c.Query(TenantIDParam))
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, apierror.APIError{
			Type: "reports-service/issues/bad_query_param/tenant_id",
			Title: "Query param is empty." +
				"The reason might be that `" + TenantIDParam + "` is empty." +
				"Please check this query param.",
			Status: http.StatusBadRequest,
			Detail: "Check json body. Reason query param:" + TenantIDParam + " " +
				"@See more in reports-service/app/internal/servers/api/controller/... -> ",
			Instance: endpointPath,
		})

		return tenantID, index, false
	}

	index = strings.TrimSpace(c.Query(IndexNameParam))

	return index, tenantID, true
}

func ParseQuerySearchPatternParamOnly(c *gin.Context, _ string) (
	searchPattern string,
	successParse bool,
) {
	return c.Query(searchPatternParam), true
}

func ParseQueryDateParamsOnly(c *gin.Context, endpointPath string) (
	fromTime time.Time,
	toTime time.Time,
	successParse bool,
) {
	fromDate := c.Query(fromDateParam)
	toDate := c.Query(toDateParam)

	var err error
	if fromDate != "" {
		fromTime, err = time.Parse(time.RFC3339, fromDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, apierror.APIError{
				Type: "reports-service/issues/bad_query_param/from",
				Title: "Query param `from_data` could not parse as time.Time obj. " +
					"The reason might be that `" + fromDateParam +
					"` query param contain time string represent not equivalent of format time.RFC3339" +
					"Please check that query param.",
				Status: http.StatusBadRequest,
				Detail: "Check json body. Reason query param:" + fromDateParam + " " +
					"@See more in reports-service/app/internal/servers/api/controller/... -> ",
				Instance: endpointPath,
			})

			return fromTime, toTime, false
		}
	}

	if toDate != "" {
		toTime, err = time.Parse(time.RFC3339, toDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, apierror.APIError{
				Type: "reports-service/issues/bad_query_param/to",
				Title: "Query param `to_data` could not parse as time.Time obj. " +
					"The reason might be that `" + toDateParam +
					"` query param contain time string represent not equivalent of format time.RFC3339" +
					"Please check that query param.",
				Status: http.StatusBadRequest,
				Detail: "Check json body. Reason query param:" + toDateParam + " " +
					"@See more in reports-service/app/internal/servers/api/controller/... -> ",
				Instance: endpointPath,
			})

			return fromTime, toTime, false
		}
	}

	return fromTime, toTime, true
}

// ParseCursorsOnly - This function extracts multiple cursor values from the query parameters.
func ParseCursorsOnly(c *gin.Context) []string {
	if cursors := c.QueryArray(cursorsParam); len(cursors) != 0 {
		return cursors
	}

	return nil
}

type (
	keyType        string
	queryParamType string
)

const (
	TenantID         = queryParamType(tenantID)         // for Auth middleware.
	RequestingUserID = queryParamType(requestingUserID) // for Auth middleware.
)

const (
	tenantID         keyType = "tenantID"
	requestingUserID keyType = "requestingUserID"
)

var (
	emptyUUID = uuid.UUID{}

	errEmptyUUID = errors.New("empty UUID param")
)

// SetTenantUUIDForRequest - set tenant uuid from jwt token.
func SetTenantUUIDForRequest(c *gin.Context, tID uuid.UUID) {
	StoreInGinCtxKV(c, tenantID, tID.String())
}

// SetRequestingUserUUIDForRequest - store requesting user id from JWT token.
func SetRequestingUserUUIDForRequest(c *gin.Context, userID uuid.UUID) {
	StoreInGinCtxKV(c, requestingUserID, userID.String())
}

// StoreInGinCtxKV - store CUSTOM key-value pairs into gin context. For DI purposes.
func StoreInGinCtxKV(c *gin.Context, key keyType, value any) {
	c.Set(string(key), value)
}

// GetUUIDFromQueryString - get UUID from query string.
func GetUUIDFromQueryString(c *gin.Context, queryParamName queryParamType) (uuid.UUID, error) {
	if v := c.Query(string(queryParamName)); v != "" {
		return uuid.Parse(v)
	}

	return emptyUUID, errEmptyUUID
}
