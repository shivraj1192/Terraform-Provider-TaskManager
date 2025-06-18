package taskmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider(version string) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("TOKEN", nil),
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BASE_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"taskmanager_user":       resourceUser(),
			"taskmanager_team":       resourceTeam(),
			"taskmanager_task":       resourceTask(),
			"taskmanager_comment":    resourceComment(),
			"taskmanager_attachment": resourceAttachment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"taskmanager_user":       dataUser(),
			"taskmanager_team":       dataTeam(),
			"taskmanager_task":       dataTask(),
			"taskmanager_comment":    dataComment(),
			"taskmanager_attachment": dataAttachment(),
		},
		ConfigureContextFunc: configureProviderClient,
	}
}

func configureProviderClient(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	baseURL := d.Get("base_url").(string)
	token := d.Get("token").(string)

	if baseURL == "" || token == "" {
		return nil, diag.Errorf("base_url and token must be set")
	}

	client := NewClient(baseURL, token)

	return client, nil
}
