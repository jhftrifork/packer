package openstack

import (
	"testing"
	"reflect"
	
	gophercloud "github.com/rackspace/gophercloud"
)

func AssertEquals(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %s but got %s", expected, actual)
	}
}

func TestPostProcessor_Configure_Good(t *testing.T) {
	var p OpenStackPostProcessor
	err := p.Configure(map[string]interface{}{
		"identity_endpoint": "http://192.168.10.5:9292/v2",
		"username":          "adminuser",
		"password":          "adminpass",
		"tenant_id":         "example_tenant",
		"irrelevant_key":    "should just be ignored",
	})

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	AssertEquals(
		t,
		gophercloud.AuthOptions {
			IdentityEndpoint: "http://192.168.10.5:9292/v2",
			Username: "adminuser",
			Password: "adminpass",
			TenantID: "example_tenant",
		},
		p.authOptions,
	)
}

func TestPostProcessor_Configure_Bad(t *testing.T) {
	badConfigs := []map[string]interface{}{
		map[string]interface{}{},
		map[string]interface{}{
			"username":          "adminuser",
			"password":          "adminpass",
			"tenant_id":         "example_tenant",
			"irrelevant_key":    "should just be ignored",
		},
	}

	for _, badConfig := range badConfigs {
		var p OpenStackPostProcessor
		err := p.Configure(badConfig)
		if err == nil {
			t.Fatalf("Expected error but got nil when configuring with: %s", badConfig)
		}
	}
}
