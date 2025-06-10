package taskmanager

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataReadAttachment,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"uploader_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func DataReadAttachment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	idInt := d.Get("id").(int)
	idStr := strconv.Itoa(idInt)

	var attachment map[string]interface{}

	if err := client.Get("api/attachments/"+idStr, &attachment); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idStr)

	result := attachment["attachment"].(map[string]interface{})

	d.Set("file_name", result["file_name"])
	d.Set("task_id", result["task_id"])
	d.Set("uploader_id", result["uploader_id"])

	return nil
}
