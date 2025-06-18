package taskmanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/bcrypt"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateUser,
		DeleteContext: resourceDeleteUser,
		UpdateContext: resourceUpdateUser,
		ReadContext:   resourceReadUser,
		Schema: map[string]*schema.Schema{
			"uname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// If password has not changed or user didn't provide a new one, ignore it
					if err := bcrypt.CompareHashAndPassword([]byte(old), []byte(new)); err != nil {
						return false
					} else {
						return true
					}
				},
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	user := map[string]interface{}{
		"uname": d.Get("uname").(string),
		"name":  d.Get("name").(string),
		"email": d.Get("email").(string),
	}

	password, ok := d.GetOk("password")
	if !ok {
		return diag.FromErr(fmt.Errorf("password field is missing"))
	}
	user["password"] = password.(string)

	if role, ok := d.GetOk("role"); ok {
		user["role"] = role.(string)
	}

	var created map[string]interface{}
	if err := client.Post("api/register", user, &created); err != nil {
		return diag.FromErr(err)
	}
	fmt.Printf("API Response on create: %+v\n", created)

	resultRaw, ok := created["user"]
	if !ok {
		return diag.Errorf("user field missing in API response: %+v", created)
	}

	result, ok := resultRaw.(map[string]interface{})
	if !ok {
		return diag.Errorf("user field has unexpected type")
	}

	id := fmt.Sprintf("%v", result["ID"])
	d.SetId(id)
	return resourceReadUser(ctx, d, m)
}

func resourceUpdateUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	user := map[string]interface{}{
		"uname": d.Get("uname").(string),
		"name":  d.Get("name").(string),
		"email": d.Get("email").(string),
	}
	updated := make(map[string]interface{})
	if err := client.Put("api/users/"+d.Id(), user, &updated); err != nil {
		return diag.FromErr(err)
	}

	return resourceReadUser(ctx, d, m)

}

func resourceReadUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	user := make(map[string]interface{})
	if err := client.Get("api/users/"+d.Id(), &user); err != nil {
		return diag.FromErr(err)
	}

	result := user["user"].(map[string]interface{})

	d.Set("uname", result["uname"])
	d.Set("name", result["name"])
	d.Set("email", result["email"])
	d.Set("password", result["password"])
	d.Set("role", result["role"])

	var teamIds []int
	if teamsRaw, ok := result["teams"].(*schema.Set); ok {
		for _, team := range teamsRaw.List() {
			if teamMap, ok := team.(map[string]interface{}); ok {
				if teamId, ok := teamMap["ID"].(float64); ok {
					teamIds = append(teamIds, int(teamId))
				}
			}
		}
	}
	d.Set("teams", teamIds)

	var teamCreatedIds []int
	if teamsRaw, ok := result["tasks_created"].(*schema.Set); ok {
		for _, team := range teamsRaw.List() {
			if teamMap, ok := team.(map[string]interface{}); ok {
				if teamId, ok := teamMap["ID"].(float64); ok {
					teamCreatedIds = append(teamCreatedIds, int(teamId))
				}
			}
		}
	}
	d.Set("tasks_created", teamCreatedIds)

	var teamAssignedIds []int
	if teamsRaw, ok := result["tasks_assigned"].(*schema.Set); ok {
		for _, team := range teamsRaw.List() {
			if teamMap, ok := team.(map[string]interface{}); ok {
				if teamId, ok := teamMap["ID"].(float64); ok {
					teamAssignedIds = append(teamAssignedIds, int(teamId))
				}
			}
		}
	}
	d.Set("tasks_assigned", teamAssignedIds)

	var commentIds []int
	if commentsRaw, ok := result["comments"].(*schema.Set); ok {
		for _, comment := range commentsRaw.List() {
			if commentMap, ok := comment.(map[string]interface{}); ok {
				if commentId, ok := commentMap["ID"].(float64); ok {
					commentIds = append(commentIds, int(commentId))
				}
			}
		}
	}
	d.Set("comments", commentIds)

	var attachmentIds []int
	if attachmentRaw, ok := result["attachments"].(*schema.Set); ok {
		for _, attachment := range attachmentRaw.List() {
			if attachmentMap, ok := attachment.(map[string]interface{}); ok {
				if attachmentId, ok := attachmentMap["ID"].(float64); ok {
					attachmentIds = append(attachmentIds, int(attachmentId))
				}
			}
		}
	}
	d.Set("attachments", attachmentIds)

	var notificationIds []int
	if notificationsRaw, ok := result["notifications"].(*schema.Set); ok {
		for _, notification := range notificationsRaw.List() {
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

func resourceDeleteUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	if err := client.Delete("api/users/" + d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
