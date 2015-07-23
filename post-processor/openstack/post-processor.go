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

type ImageCreateConfig struct {
	Name *string `mapstructure:"name"`
	Id *string `mapstructure:"id"`
	Visibility *imageservice.ImageVisibility `mapstructure:"visibility"`
	Tags []string `mapstructure:"tags"`
	ContainerFormat *string `mapstructure:"container_format"`
	DiskFormat *string `mapstructure:"disk_format"`
	MinDiskGigabytes *int `mapstructure:"min_disk_gigabytes"`
	MinRamMegabytes *int `mapstructure:"min_ram_megabytes"`
	Protected *bool `mapstructure:"protected"`
	Properties map[string]string `mapstructure:"properties"`
}

func toImageServiceCreateOpts(c ImageCreateConfig) imageservice.CreateOpts {
	return imageservice.CreateOpts{
		Name: c.Name,
		Id: c.Id,
		Visibility: c.Visibility,
		Tags: c.Tags,
		ContainerFormat: c.ContainerFormat,
		DiskFormat: c.DiskFormat,
		MinDiskGigabytes: c.MinDiskGigabytes,
		MinRamMegabytes: c.MinRamMegabytes,
		Protected: c.Protected,
		Properties: c.Properties,
	}
}

type Config struct {
	IdentityEndpoint *string `mapstructure:"identity_endpoint"`
	Username *string `mapstructure:"username"`
	Password *string `mapstructure:"password"`
	TenantId *string `mapstructure:"tenant_id"`

	ImageCreateConfig *ImageCreateConfig `mapstructure:"image"`
	
	ctx interpolate.Context // wtf is this?
}

// implements packer.PostProcessor
type OpenStackPostProcessor struct {
	config Config
	authOptions gophercloud.AuthOptions      // constructed from details in `config`
	imageCreateOpts imageservice.CreateOpts  // constructed from details in `config`
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

	var errs *packer.MultiError
	if p.config.IdentityEndpoint == nil {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("%s must be set", "identity_endpoint"))
	}
	
	if p.config.Username == nil {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("%s must be set", "username"))
	}
	
	if p.config.Password == nil {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("%s must be set", "password"))
	}
	
	if p.config.TenantId == nil {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("%s must be set", "tenant_id"))
	}
	
	if p.config.ImageCreateConfig == nil {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("%s must be set", "image"))
	}
	
	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	p.authOptions = gophercloud.AuthOptions {
		IdentityEndpoint: *p.config.IdentityEndpoint,
		Username: *p.config.Username,
		Password: *p.config.Password,
		TenantID: *p.config.TenantId,
	}

	p.imageCreateOpts = toImageServiceCreateOpts(*p.config.ImageCreateConfig)

	// We don't instantiate a ProviderClient here because doing so
	// has side-effects: making requests to the OpenStack
	// instance. Configure should be pure.

	return nil
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

	createResult := imageservice.Create(serviceClient, p.imageCreateOpts)

	_, createErr := createResult.Extract()
	if createErr != nil {
		return nil, true, createErr
	}

	//image := *imagePtr
	
	// TODO

	return
}
