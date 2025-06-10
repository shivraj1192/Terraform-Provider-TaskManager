package taskmanager

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataComment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataReadComment,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parent_comment_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subcomments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataReadComment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	idInt := d.Get("id").(int)
	idStr := strconv.Itoa(idInt)

	var comment map[string]interface{}

	log.Println("[INFO] Getting Comment")

	if err := client.Get("api/comments/"+idStr, &comment); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idStr)

	result := comment["comment"].(map[string]interface{})

	d.Set("content", result["content"])
	d.Set("user_id", result["user_id"])
	d.Set("task_id", result["task_id"])
	d.Set("parent_comment_id", result["parent_comment_id"])
	d.Set("subcomments", result["subcomments"])

	return nil
}
