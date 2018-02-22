package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	nn_collector "code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/collector"
)

// AppGUID represets an application GUID.
type AppGUID string

// HTTPClient supports HTTP requests.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// HTTPAppInfoStore provides a focused source of Cloud Controller API data.
type HTTPAppInfoStore struct {
	apiAddr string
	client  HTTPClient
}

// NewCFLightApiAppInfoStore initializes an APIStore and sends all HTTP requests to
// the API URL specified by apiAddr.
func NewCFLightApiAppInfoStore(apiAddr string, client HTTPClient) nn_collector.AppInfoStore {
	return &HTTPAppInfoStore{
		apiAddr: apiAddr,
		client:  client,
	}
}

// Lookup reads AppInfo from a remote API. Because of how the Light api works the apps are not
// currently filtered against their guids
func (s *HTTPAppInfoStore) Lookup(guids []string) (
	map[nn_collector.AppGUID]nn_collector.AppInfo, error) {

	log.Println("Looking up apps...")
	if len(guids) < 1 {
		return nil, nil
	}

	apps, err := s.lookupApps(guids)
	if err != nil {
		return nil, err
	}

	res := make(map[nn_collector.AppGUID]nn_collector.AppInfo)
	for _, app := range apps {

		if app.GUID != "" {
			res[app.GUID] = nn_collector.AppInfo{
				Name:  app.Name,
				Space: app.Space,
				Org:   app.Org,
			}
		}
	}

	return res, nil
}

func (s *HTTPAppInfoStore) lookupApps(guids []string) (CFLightResponse, error) {
	u, err := url.Parse(fmt.Sprintf("%s/v2/apps", s.apiAddr))
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, err
	}

	r, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		buf := bytes.NewBuffer(nil)
		_, _ = buf.ReadFrom(r.Body)
		err := fmt.Errorf("failed to get apps, expected 200, got %d: %s", r.StatusCode, buf.String())

		return nil, err
	}

	var apps CFLightResponse
	err = json.NewDecoder(r.Body).Decode(&apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// V3Resource represents application data returned from the Cloud Controller
// API.
type CFLightApp struct {
	GUID  nn_collector.AppGUID `json:"guid"`
	Name  string               `json:"name"`
	Space string               `json:"space"`
	Org   string               `json:"org"`
}

// V3Response represents a list of V3 API resources and associated data.
type CFLightResponse []CFLightApp

// AppInfo holds the names of an application, space, and organization.
type AppInfo struct {
	Name  string
	Space string
	Org   string
}

// String implements the Stringer interface.
func (a AppInfo) String() string {
	return fmt.Sprintf("%s.%s.%s", a.Org, a.Space, a.Name)
}
