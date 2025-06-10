package taskmanager

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataReadTeam,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataReadTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	idInt := d.Get("id").(int)
	idStr := strconv.Itoa(idInt)

	var outer map[string]interface{}
	if err := client.Get("api/teams/"+idStr, &outer); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idStr)

	result, ok := outer["team"].(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("unable to get team"))
	}

	d.Set("name", result["name"])
	d.Set("description", result["description"])

	if ownerID, ok := result["owner_id"].(float64); ok {
		d.Set("owner_id", int(ownerID))
	}

	var memberIDs []int
	if membersRaw, ok := result["members"].([]interface{}); ok {
		for _, member := range membersRaw {
			if memberMap, ok := member.(map[string]interface{}); ok {
				if idFloat, ok := memberMap["ID"].(float64); ok {
					memberIDs = append(memberIDs, int(idFloat))
				}
			}
		}
	}
	d.Set("members", memberIDs)

	var taskIDs []int
	if tasksRaw, ok := result["tasks"].([]interface{}); ok {
		for _, task := range tasksRaw {
			if taskMap, ok := task.(map[string]interface{}); ok {
				if idFloat, ok := taskMap["id"].(float64); ok {
					taskIDs = append(taskIDs, int(idFloat))
				}
			}
		}
	}
	d.Set("tasks", taskIDs)

	return nil
}
