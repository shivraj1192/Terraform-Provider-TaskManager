package taskmanager

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataReadUser,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"uname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"tasks_created": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"tasks_assigned": {
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
			"notifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataReadUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	idInt := d.Get("id").(int)
	idStr := strconv.Itoa(idInt)

	user := make(map[string]interface{})
	if err := client.Get("api/users/"+idStr, &user); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idStr)

	result := user["user"].(map[string]interface{})

	d.Set("uname", result["uname"])
	d.Set("name", result["name"])
	d.Set("email", result["email"])
	d.Set("password", result["password"])
	d.Set("role", result["role"])

	var teamIds []int
	if teamsRaw, ok := result["teams"].([]interface{}); ok {
		for _, team := range teamsRaw {
			if teamMap, ok := team.(map[string]interface{}); ok {
				if teamId, ok := teamMap["ID"].(float64); ok {
					teamIds = append(teamIds, int(teamId))
				}
			}
		}
	}
	d.Set("teams", teamIds)

	var teamCreatedIds []int
	if teamsRaw, ok := result["tasks_created"].([]interface{}); ok {
		for _, team := range teamsRaw {
			if teamMap, ok := team.(map[string]interface{}); ok {
				if teamId, ok := teamMap["ID"].(float64); ok {
					teamCreatedIds = append(teamCreatedIds, int(teamId))
				}
			}
		}
	}
	d.Set("tasks_created", teamCreatedIds)

	var teamAssignedIds []int
	if teamsRaw, ok := result["tasks_assigned"].([]interface{}); ok {
		for _, team := range teamsRaw {
			if teamMap, ok := team.(map[string]interface{}); ok {
				if teamId, ok := teamMap["ID"].(float64); ok {
					teamAssignedIds = append(teamAssignedIds, int(teamId))
				}
			}
		}
	}
	d.Set("tasks_assigned", teamAssignedIds)

	var commentIds []int
	if commentsRaw, ok := result["comments"].([]interface{}); ok {
		for _, comment := range commentsRaw {
			if commentMap, ok := comment.(map[string]interface{}); ok {
				if commentId, ok := commentMap["ID"].(float64); ok {
					commentIds = append(commentIds, int(commentId))
				}
			}
		}
	}
	d.Set("comments", commentIds)

	var attachmentIds []int
	if attachmentRaw, ok := result["attachments"].([]interface{}); ok {
		for _, attachment := range attachmentRaw {
			if attachmentMap, ok := attachment.(map[string]interface{}); ok {
				if attachmentId, ok := attachmentMap["ID"].(float64); ok {
					attachmentIds = append(attachmentIds, int(attachmentId))
				}
			}
		}
	}
	d.Set("attachments", attachmentIds)

	var notificationIds []int
	if notificationsRaw, ok := result["notifications"].([]interface{}); ok {
		for _, notification := range notificationsRaw {
			if notificationMap, ok := notification.(map[string]interface{}); ok {
				if notificationId, ok := notificationMap["ID"].(float64); ok {
					notificationIds = append(notificationIds, int(notificationId))
				}
			}
		}
	}
	d.Set("notifications", notificationIds)

	return nil
}
