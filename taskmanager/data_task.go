package taskmanager

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataReadTask,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"due_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parent_task_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subtasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"assignees": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"comments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataReadTask(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	idInt := d.Get("id").(int)
	idStr := strconv.Itoa(idInt)

	var task map[string]interface{}
	if err := client.Get("api/tasks/"+idStr, &task); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idStr)

	result, ok := task["task"].(map[string]interface{})
	if !ok {
		return diag.Errorf("unable to get team")
	}

	d.Set("title", result["title"])
	d.Set("description", result["description"])
	d.Set("status", result["status"])
	d.Set("priority", result["priority"])
	d.Set("creator_id", result["creator_id"])
	d.Set("team_id", result["team_id"])
	d.Set("assignees", result["assignees"])
	d.Set("parent_task_id", result["parent_task_id"])
	d.Set("subtasks", result["subtasks"])
	d.Set("labels", result["labels"])
	d.Set("comments", result["comments"])
	d.Set("attachments", result["attachments"])

	return nil
}
