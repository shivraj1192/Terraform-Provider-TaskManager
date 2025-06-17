package taskmanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateTeam,
		ReadContext:   resourceReadTeam,
		UpdateContext: resourceUpdateTeam,
		DeleteContext: resourceDeleteTeam,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"members": {
				Type:     schema.TypeList,
				Optional: true,
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

func resourceCreateTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	team := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		team["description"] = desc.(string)
	}

	var created map[string]interface{}
	if err := client.Post("api/teams", team, &created); err != nil {
		return diag.FromErr(err)
	}

	result, ok := created["team"].(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("unable to get team"))
	}
	id := fmt.Sprintf("%v", result["ID"])
	d.SetId(id)

	if members, ok := d.GetOk("members"); ok {
		if err := addMembers(client, members.([]interface{}), d); err != nil {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Members not added",
				Detail:   err.Error(),
			}}
		}
	}

	return waitThenRead(ctx, d, m)
}

func resourceUpdateTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	team := map[string]interface{}{
		"name": d.Get("name").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		team["description"] = desc.(string)
	}

	var updated map[string]interface{}
	if err := client.Put("api/teams/"+d.Id(), team, &updated); err != nil {
		return diag.FromErr(err)
	}

	if members, ok := d.GetOk("members"); ok {
		if err := addMembers(client, members.([]interface{}), d); err != nil {
			return diag.Diagnostics{diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Members not updated",
				Detail:   err.Error(),
			}}
		}
	}

	return waitThenRead(ctx, d, m)
}

func resourceReadTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	var outer map[string]interface{}
	if err := client.Get("api/teams/"+d.Id(), &outer); err != nil {
		return diag.FromErr(err)
	}

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

func resourceDeleteTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	if err := client.Delete("api/teams/" + d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func addMembers(client *TaskManagerClient, members []interface{}, d *schema.ResourceData) error {
	team := map[string]interface{}{
		"userids": members,
	}

	var updated map[string]interface{}
	if err := client.Put("api/teams/"+d.Id()+"/add-members", team, &updated); err != nil {
		return err
	}
	return nil
}

func waitThenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	time.Sleep(1 * time.Second) // Add a 1-second delay (adjust if needed)
	return resourceReadTeam(ctx, d, m)
}
