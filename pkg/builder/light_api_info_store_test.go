package builder_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	nn_collector "code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/collector"

	"github.com/SpringerPE/noisy-neighbor-reporters/pkg/builder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTPAppInfoStore", func() {
	It("issues GET requests to Cloud Controller for AppInfo", func() {
		client := &fakeHTTPClient{responses: happyPath()}
		store := builder.NewCFLightApiAppInfoStore("http://api.addr.com", client)

		actual, err := store.Lookup([]string{"a", "b"})

		Expect(err).ToNot(HaveOccurred())
		expected := map[nn_collector.AppGUID]nn_collector.AppInfo{
			"a": nn_collector.AppInfo{
				Name:  "app1",
				Space: "space1",
				Org:   "org1",
			},
			"b": nn_collector.AppInfo{
				Name:  "app2",
				Space: "space2",
				Org:   "org2",
			},
		}
		Expect(actual).To(Equal(expected))
		Expect(client.requests).To(HaveLen(1))

		req := client.requests[0]
		Expect(req.URL.Host).To(Equal("api.addr.com"))
		Expect(req.URL.Path).To(Equal("/v2/apps"))

	})

	It("returns an empty map when no GUIDInstances are passed in", func() {
		client := &fakeHTTPClient{responses: happyPath()}
		store := builder.NewCFLightApiAppInfoStore("http://api.addr.com", client)

		data, _ := store.Lookup(nil)

		Expect(data).To(HaveLen(0))
		Expect(client.doCalled).To(BeFalse())
	})

	It("returns an error when getting apps fails", func() {
		client := &fakeHTTPClient{responses: appRequestFailed()}
		store := builder.NewCFLightApiAppInfoStore("http://api.addr.com", client)

		_, err := store.Lookup([]string{"a", "b"})
		Expect(err).To(HaveOccurred())
	})

	It("returns an error when the app request is not a 200 status", func() {
		client := &fakeHTTPClient{responses: appRequestNotOK()}
		store := builder.NewCFLightApiAppInfoStore("http://api.addr.com", client)

		_, err := store.Lookup([]string{"a", "b"})
		Expect(err).To(HaveOccurred())
	})

	It("returns an error when app json unmarshalling fails", func() {
		client := &fakeHTTPClient{responses: appRequestInvalidJSON()}
		store := builder.NewCFLightApiAppInfoStore("http://api.addr.com", client)

		_, err := store.Lookup([]string{"a", "b"})
		Expect(err).To(HaveOccurred())
	})
})

type fakeHTTPClient struct {
	doCalled  bool
	responses map[string]response
	requests  []*http.Request
}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	f.doCalled = true
	f.requests = append(f.requests, req)
	resp, ok := f.responses[req.URL.Path]
	if !ok {
		return nil, nil
	}
	return resp.http, resp.err
}

type response struct {
	http *http.Response
	err  error
}

func happyPath() map[string]response {
	appsResp := response{
		http: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(appsResponse())),
		},
		err: nil,
	}
	return map[string]response{
		"/v2/apps": appsResp,
	}
}

func appRequestFailed() map[string]response {
	appsResp := response{
		http: nil,
		err:  errors.New("request failed"),
	}
	return map[string]response{
		"/v2/apps": appsResp,
	}
}

func appRequestNotOK() map[string]response {
	appsResp := response{
		http: &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader("")),
		},
		err: nil,
	}

	return map[string]response{
		"/v2/apps": appsResp,
	}
}

func appRequestInvalidJSON() map[string]response {
	appsResp := response{
		http: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("{")),
		},
		err: nil,
	}

	return map[string]response{
		"/v2/apps": appsResp,
	}
}

func appsResponse() string {
	return `
[
	{
	  "name": "app1",
	  "production": false,
	  "space_guid": "abf1b2dc-5c27-49cb-b7a5-985021a50b57",
	  "stack_guid": "dedf751a-1eac-4e81-967e-9eafc19444c6",
	  "buildpack": null,
	  "detected_buildpack": null,
	  "detected_buildpack_guid": null,
	  "environment_json": {},
	  "memory": 256,
	  "instances": [],
	  "disk_quota": 1024,
	  "state": "STOPPED",
	  "version": "68489fc0-93c5-464b-b588-284c2675f68a",
	  "command": null,
	  "console": false,
	  "debug": null,
	  "staging_task_id": null,
	  "package_state": "PENDING",
	  "health_check_type": "port",
	  "health_check_timeout": null,
	  "health_check_http_endpoint": null,
	  "staging_failed_reason": null,
	  "staging_failed_description": null,
	  "diego": false,
	  "docker_image": null,
	  "docker_credentials": {
	    "username": null,
	    "password": null
	  },
	  "package_updated_at": null,
	  "detected_start_command": "",
	  "enable_ssh": true,
	  "ports": null,
	  "space_url": "/v2/spaces/abf1b2dc-5c27-49cb-b7a5-985021a50b57",
	  "stack_url": "/v2/stacks/dedf751a-1eac-4e81-967e-9eafc19444c6",
	  "stack": "cflinuxfs2",
	  "routes_url": "/v2/apps/63e0e006-b690-4938-9a11-542a58ff6b75/routes",
	  "routes": [],
	  "events_url": "/v2/apps/63e0e006-b690-4938-9a11-542a58ff6b75/events",
	  "service_bindings_url": "/v2/apps/63e0e006-b690-4938-9a11-542a58ff6b75/service_bindings",
	  "route_mappings_url": "/v2/apps/63e0e006-b690-4938-9a11-542a58ff6b75/route_mappings",
	  "created_at": "2015-08-26T12:51:31Z",
	  "updated_at": "2015-08-26T12:51:31Z",
	  "guid": "a",
	  "meta": {
	    "error": false
	  },
	  "space": "space1",
	  "org": "org1",
	  "running": false
	},
	{
	  "name": "app2",
	  "production": false,
	  "space_guid": "abf1b2dc-5c27-49cb-b7a5-985021a50b57",
	  "stack_guid": "dedf751a-1eac-4e81-967e-9eafc19444c6",
	  "buildpack": null,
	  "detected_buildpack": null,
	  "detected_buildpack_guid": null,
	  "environment_json": {},
	  "memory": 256,
	  "instances": [],
	  "disk_quota": 1024,
	  "state": "STOPPED",
	  "version": "68489fc0-93c5-464b-b588-284c2675f68a",
	  "command": null,
	  "console": false,
	  "debug": null,
	  "staging_task_id": null,
	  "package_state": "PENDING",
	  "health_check_type": "port",
	  "health_check_timeout": null,
	  "health_check_http_endpoint": null,
	  "staging_failed_reason": null,
	  "staging_failed_description": null,
	  "diego": false,
	  "docker_image": null,
	  "docker_credentials": {
	    "username": null,
	    "password": null
	  },
	  "package_updated_at": null,
	  "detected_start_command": "",
	  "enable_ssh": true,
	  "ports": null,
	  "space_url": "/v2/spaces/abf1b2dc-5c27-49cb-b7a5-985021a50b57",
	  "stack_url": "/v2/stacks/dedf751a-1eac-4e81-967e-9eafc19444c6",
	  "stack": "cflinuxfs2",
	  "routes_url": "/v2/apps/63e0e006-b690-4938-9a11-542a58ff6b75/routes",
	  "routes": [],
	  "events_url": "/v2/apps/63e0e006-b690-4938-9a11-542a58ff6b75/events",
	  "service_bindings_url": "/v2/apps/63e0e006-b690-4938-9a11-542a58ff6b75/service_bindings",
	  "route_mappings_url": "/v2/apps/63e0e006-b690-4938-9a11-542a58ff6b75/route_mappings",
	  "created_at": "2015-08-26T12:51:31Z",
	  "updated_at": "2015-08-26T12:51:31Z",
	  "guid": "b",
	  "meta": {
	    "error": false
	  },
	  "space": "space2",
	  "org": "org2",
	  "running": false
	}
]
`
}
