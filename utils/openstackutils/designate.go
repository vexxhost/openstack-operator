// Copyright 2020 VEXXHOST, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openstackutils

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"opendev.org/vexxhost/openstack-operator/utils/tlsutils"
)

const (
	_openAPIVersion = "3"
)

// DesignateClientBuilder is an implementation of the designateClientInterface
type DesignateClientBuilder struct {
	ServiceClient *gophercloud.ServiceClient
	isAuth        bool
}

// CloudYAML is for parsing the clouds.yaml
type CloudYAML struct {
	Clouds map[string]struct {
		Auth struct {
			Auth_url            string `yaml:"auth_url"`
			Project_name        string `yaml:"project_name"`
			Project_id          string `yaml:"project_id"`
			Username            string `yaml:"username"`
			Password            string `yaml:"password"`
			User_domain_name    string `yaml:"user_domain_name"`
			Project_domain_name string `yaml:"project_domain_name"`
		} `yaml:"auth"`
		Region_name string `yaml:"region_name"`
		Interface   string `yaml:"interface"`
	} `yaml:"clouds"`
}

// DesignateClient is a constructor for the DesignateBuilder
func DesignateClient(existing *DesignateClientBuilder, rc []byte, cloudName string) error {
	if existing.GetAuthStatus() {
		log.Infof("Already authenticated")
		return nil
	}
	if err := setAuthSettings(rc, cloudName); err != nil {
		log.Infof("Authentication failed - %s", err.Error())
		return err
	}
	serviceClient, err := createDesignateServiceClient()
	if err != nil {
		log.Infof("createDesignateServiceClient failed - %s", err.Error())
		return err
	}

	existing.ServiceClient = serviceClient
	existing.SetAuthSuccess()
	log.Infof("Authentication success!")
	return nil
}

// CreateZone creates a zone
func (c *DesignateClientBuilder) CreateZone(dn string, ttl int, email string) (string, error) {
	// zone create

	createOpts := zones.CreateOpts{
		Name:        dn,
		Email:       email,
		Type:        "PRIMARY",
		TTL:         ttl,
		Description: "This is a zone.",
	}

	res := zones.Create(c.ServiceClient, createOpts)
	if res.Err != nil {
		log.Errorf("Create Zone failed - %s", res.Err.Error())
		c.SetAuthFailed()
		return "", res.Err
	}

	log.Infof("Gained zone infor successfully")
	zoneInfo, err := res.Extract()
	if err != nil {
		c.SetAuthFailed()
		log.Errorf("Extract zone infor failed")
		return "", err
	}
	return zoneInfo.ID, err

}

// UpdateZone updates zone TTL and Email.
func (c *DesignateClientBuilder) UpdateZone(zoneID string, TTL int, Email string) error {
	updateOpts := zones.UpdateOpts{
		TTL:   TTL,
		Email: Email,
	}
	if err := zones.Update(c.ServiceClient, zoneID, updateOpts).Err; err != nil {
		log.Errorf("Update zone failed")
		c.SetAuthFailed()
		return err
	}
	return nil
}

// DeleteZone deletes a zone
func (c *DesignateClientBuilder) DeleteZone(Domain string) error {
	zoneList, err := c.ListZone()
	if err != nil {
		return err
	}
	for _, zone := range zoneList {
		if zone.Name == Domain {
			return c.deleteZoneByID(zone.ID)
		}
	}
	log.Infof("No such zone exists to delete.")
	return nil
}

// ListZone gets the zone list
func (c *DesignateClientBuilder) ListZone() ([]zones.Zone, error) {
	listOpts := zones.ListOpts{}
	allPages, err := zones.List(c.ServiceClient, listOpts).AllPages()
	if err != nil {
		log.Errorf("List zone list failed")
		c.SetAuthFailed()
		return []zones.Zone{}, err
	}

	allZones, err := zones.ExtractZones(allPages)
	if err != nil {
		log.Errorf("Extract zone infor from the zone list failed")
		c.SetAuthFailed()
		return []zones.Zone{}, err
	}
	return allZones, nil
}

// CreateOrUpdateZone sync the zone list
func (c *DesignateClientBuilder) CreateOrUpdateZone(Domain string, TTL int, Email string) error {
	zoneList, err := c.ListZone()
	if err != nil {
		return err
	}
	for _, zone := range zoneList {
		if Domain == zone.Name {
			// Update Zone
			log.Infof("Designate: Zone %s already exists", zone.Name)
			return c.UpdateZone(zone.ID, TTL, Email)
		}
	}

	// Create zone
	_, err = c.CreateZone(Domain, TTL, Email)

	return err
}

// deleteZoneByID deletes the zone using zoneID, consuming the designate API directly
func (c *DesignateClientBuilder) deleteZoneByID(zoneID string) error {
	if err := zones.Delete(c.ServiceClient, zoneID).Err; err != nil {
		c.SetAuthFailed()
		return err
	}
	return nil
}

// SetAuthSuccess means the current client already authenticated
func (c *DesignateClientBuilder) SetAuthSuccess() {
	c.isAuth = true
}

// SetAuthFailed means the current client needs to authenticate
func (c *DesignateClientBuilder) SetAuthFailed() {
	c.isAuth = false
}

// GetAuthStatus returns the authentication status
func (c *DesignateClientBuilder) GetAuthStatus() bool {
	return c.isAuth
}

// authenticate in OpenStack and obtain Designate service endpoint
func createDesignateServiceClient() (*gophercloud.ServiceClient, error) {
	opts, err := getAuthSettings()
	if err != nil {
		return nil, err
	}
	log.Infof("Using OpenStack Keystone at %s", opts.IdentityEndpoint)
	authProvider, err := openstack.NewClient(opts.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := tlsutils.CreateTLSConfig("OPENSTACK")
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
	}
	authProvider.HTTPClient.Transport = transport

	if err = openstack.Authenticate(authProvider, opts); err != nil {
		return nil, err
	}

	eo := gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	}

	client, err := openstack.NewDNSV2(authProvider, eo)
	if err != nil {
		return nil, err
	}
	log.Infof("Found OpenStack Designate service at %s", client.Endpoint)
	return client, nil
}

// returns OpenStack Keystone authentication settings by obtaining values from standard environment variables.
// also fixes incompatibilities between gophercloud implementation and *-stackrc files that can be downloaded
// from OpenStack dashboard in latest versions
func getAuthSettings() (gophercloud.AuthOptions, error) {
	remapEnv(map[string]string{
		"OS_TENANT_NAME": "OS_PROJECT_NAME",
		"OS_TENANT_ID":   "OS_PROJECT_ID",
		"OS_DOMAIN_NAME": "OS_USER_DOMAIN_NAME",
		"OS_DOMAIN_ID":   "OS_USER_DOMAIN_ID",
	})

	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return gophercloud.AuthOptions{}, err
	}
	opts.AllowReauth = true
	if !strings.HasSuffix(opts.IdentityEndpoint, "/") {
		opts.IdentityEndpoint += "/"
	}
	if !strings.HasSuffix(opts.IdentityEndpoint, "/v2.0/") && !strings.HasSuffix(opts.IdentityEndpoint, "/v3/") {
		opts.IdentityEndpoint += "v2.0/"
	}
	return opts, nil
}

// copies environment variables to new names without overwriting existing values
func remapEnv(mapping map[string]string) {
	for k, v := range mapping {
		currentVal := os.Getenv(k)
		newVal := os.Getenv(v)
		if currentVal == "" && newVal != "" {
			os.Setenv(k, newVal)
		}
	}
}

func setAuthSettings(rc []byte, cloudName string) error {
	var cloudYaml CloudYAML
	parseCloudYAML(rc, &cloudYaml)
	credential, ok := cloudYaml.Clouds[cloudName]
	if !ok {
		return fmt.Errorf("rc secret does not involve the current cloud credential ")
	}
	os.Setenv("OS_AUTH_URL", credential.Auth.Auth_url)
	os.Setenv("OS_PROJECT_ID", credential.Auth.Project_id)
	os.Setenv("OS_PROJECT_NAME", credential.Auth.Project_name)
	os.Setenv("OS_USER_DOMAIN_NAME", credential.Auth.User_domain_name)
	os.Setenv("OS_USERNAME", credential.Auth.Username)
	os.Setenv("OS_PASSWORD", credential.Auth.Password)
	os.Setenv("OS_REGION_NAME", credential.Region_name)
	os.Setenv("OS_INTERFACE", credential.Interface)
	os.Setenv("OS_IDENTITY_API_VERSION", _openAPIVersion)
	return nil
}

func parseCloudYAML(y []byte, cloudYaml *CloudYAML) {
	err := yaml.Unmarshal([]byte(y), cloudYaml)
	if err != nil {
		panic(err)
	}
}
