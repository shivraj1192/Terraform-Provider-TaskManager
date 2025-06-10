package taskmanager

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateComment,
		UpdateContext: resourceUpdateComment,
		ReadContext:   resourceReadComment,
		DeleteContext: resourceDeleteComment,
		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"task_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"parent_comment_id": {
				Type:     schema.TypeInt,
				Optional: true,
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

func resourceCreateComment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	log.Println("[INFO] getting details")

	comment := map[string]interface{}{
		"content": d.Get("content").(string),
	}

	taskId := d.Get("task_id").(int)
	taskIdStr := strconv.Itoa(taskId)

	if parentCommentId, ok := d.GetOk("parent_comment_id"); ok && parentCommentId.(int) > 0 {
		comment["parent_comment_id"] = parentCommentId.(int)
	}

	log.Println("[INFO] Creationg Comment")

	var created map[string]interface{}
	if err := client.Post("api/tasks/"+taskIdStr+"/comments", comment, &created); err != nil {
		return diag.FromErr(err)
	}

	result := created["comment"].(map[string]interface{})

	id := fmt.Sprintf("%v", result["ID"])
	d.SetId(id)

	return resourceReadComment(ctx, d, m)
}

func resourceUpdateComment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	comment := map[string]interface{}{
		"content": d.Get("content").(string),
	}
	if parentCommentId, ok := d.GetOk("parent_comment_id"); ok && parentCommentId.(int) > 0 {
		comment["parent_comment_id"] = parentCommentId.(int)
	}

	var updated map[string]interface{}
	if err := client.Put("api/comments/"+d.Id(), comment, &updated); err != nil {
		return diag.FromErr(err)
	}

	return resourceReadComment(ctx, d, m)
}

func resourceReadComment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	var comment map[string]interface{}

	log.Println("[INFO] Getting Comment")

	if err := client.Get("api/comments/"+d.Id(), &comment); err != nil {
		return diag.FromErr(err)
	}

	result := comment["comment"].(map[string]interface{})

	d.Set("content", result["content"])
	d.Set("user_id", result["user_id"])
	d.Set("task_id", result["task_id"])
	d.Set("parent_comment_id", result["parent_comment_id"])
	d.Set("subcomments", result["subcomments"])

	return nil
}

func resourceDeleteComment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	if err := client.Delete("api/comments/" + d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
