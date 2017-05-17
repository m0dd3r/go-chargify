package chargify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

const (
	libraryVersion = "0"
	userAgent      = "go-chargify/" + libraryVersion
)

var (
	baseURLTemplate *template.Template
)

func init() {
	var err error
	baseURLTemplate, err = template.New("baseUrl").Parse("https://{{.Subdomain}}.chargify.com/")
	if err != nil {
		panic(err)
	}
}

type Client struct {
	client        *http.Client
	BaseURL       *url.URL
	ApiKey        string
	UserAgent     string
	common        service
	Subscriptions *SubscriptionsService
	Products      *ProductsService
}

type service struct {
	client *Client
}

// NewClient returns a new Chargify API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(subdomain, api_key string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	var buf bytes.Buffer
	err := baseURLTemplate.Execute(&buf, struct{ Subdomain string }{subdomain})
	if err != nil {
		panic(err)
	}
	baseURL, _ := url.Parse(buf.String())

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent, ApiKey: api_key}
	c.common.client = c
	c.Subscriptions = (*SubscriptionsService)(&c.common)
	c.Products = (*ProductsService)(&c.common)
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.ApiKey, "X")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Response is a Chargify API response. This wraps the standard http.Response
// returned from chargify and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response

	// These fields provide the page values for paginating through a set of
	// results. Any or all of these may be set to the zero value for
	// responses that are not part of a paginated set, or for which there
	// are no additional pages.

	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	//response.populatePageValues()
	return response
}

type FormattedTime struct {
	*time.Time
}

const timeFormat = `"2006-01-02T15:04:05-07:00"`

func (f FormattedTime) MarshalJSON() ([]byte, error) {
	return []byte(f.Format(timeFormat)), nil
}

func (f *FormattedTime) UnmarshalJSON(input []byte) (err error) {
	t, err := time.Parse(timeFormat, string(input))
	*f = FormattedTime{&t}
	return
}

func NewFormattedTime(input string) *FormattedTime {
	t, _ := time.Parse(timeFormat, input)
	return &FormattedTime{&t}
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it. If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	ctx, req = withContext(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

func withContext(ctx context.Context, req *http.Request) (context.Context, *http.Request) {
	return ctx, req.WithContext(ctx)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range.
// API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
//
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

/*
An ErrorResponse reports one or more errors caused by an API request.
Chargify API docs: https://developer.github.com/v3/#client-errors
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Errors   []string       `json:"errors"` // error messages
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Errors)
}
