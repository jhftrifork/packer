package openstack

// TODO
import (
	"fmt"

	packer "github.com/mitchellh/packer/packer"
	config "github.com/mitchellh/packer/helper/config"
	interpolate "github.com/mitchellh/packer/template/interpolate"

	gophercloud "github.com/rackspace/gophercloud"
	openstack "github.com/rackspace/gophercloud/openstack"
	imageservice "github.com/rackspace/gophercloud/openstack/imageservice/v2"
)

type Config struct {
	IdentityEndpoint string `mapstructure:"identity_endpoint"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	TenantId string `mapstructure:"tenant_id"`

	ImageCreateOpts imageservice.CreateOpts `mapstructure:"image"`
	
	ctx interpolate.Context // wtf is this?
}

// implements packer.PostProcessor
type OpenStackPostProcessor struct {
	config Config
	authOptions gophercloud.AuthOptions  // constructed from details in `config`
}

// Configure is responsible for setting up configuration, storing the
// state for later, and returning and errors, such as validation
// errors. Configuration is taken from `raws` and placed in
// `p`. Success is indicated by the return value `nil`.
func (p *OpenStackPostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate: true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
		AllowUnusedKeys: true,
	}, raws...)
	if err != nil {
		return err
	}

	required := map[string]*string {
		"identity_endpoint": &p.config.IdentityEndpoint,
		"username": &p.config.Username,
		"password": &p.config.Password,
		"tenant_id": &p.config.TenantId,
	}

	errs := checkAllStringsNotEmpty(required)
	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	p.authOptions = gophercloud.AuthOptions {
		IdentityEndpoint: p.config.IdentityEndpoint,
		Username: p.config.Username,
		Password: p.config.Password,
		TenantID: p.config.TenantId,
	}

	// We don't instantiate a ProviderClient here because doing so
	// has side-effects: making requests to the OpenStack
	// instance. Configure should be pure.

	return nil
}

/// If there exists an `s` for which `m[s] == ""`, returns an error
/// reporting all such `s`.  Otherwise, returns `nil`.
func checkAllStringsNotEmpty(m map[string]*string) *packer.MultiError {
	var errs *packer.MultiError
	for key, ptr := range m {
		if *ptr == "" {
			errs = packer.MultiErrorAppend(
				errs,
				fmt.Errorf("%s must be set", key))
		}
	}
	return errs
}

// PostProcess takes a previously created Artifact and produces another
// Artifact. If an error occurs, it should return that error. If `keep`
// is to true, then the previous artifact is forcibly kept.
func (p *OpenStackPostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (a packer.Artifact, keep bool, err error) {

	providerClient, err := openstack.AuthenticatedClient(p.authOptions)

	if err != nil {
		return nil, true, err
	}
	
	serviceClient := openstack.NewIdentityV3(providerClient)

	createResult := imageservice.Create(serviceClient, p.config.ImageCreateOpts)

	_, createErr := createResult.Extract()
	if createErr != nil {
		return nil, true, createErr
	}

	//image := *imagePtr
	
	// TODO

	return
}
