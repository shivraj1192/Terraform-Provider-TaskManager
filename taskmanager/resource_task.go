package taskmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateTask,
		ReadContext:   resourceReadTask,
		UpdateContext: resourceUpdateTask,
		DeleteContext: resourceDeleteTask,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"due_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"parent_task_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"subtasks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"assignees": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"comments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"attachments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceCreateTask(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	task := map[string]interface{}{
		"title":   d.Get("title").(string),
		"team_id": d.Get("team_id").(int),
	}

	if description, ok := d.GetOk("description"); ok {
		task["description"] = description.(string)
	}
	if status, ok := d.GetOk("status"); ok {
		task["status"] = status.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		task["priority"] = priority.(string)
	}
	if due_date, ok := d.GetOk("due_date"); ok {
		task["due_date"] = due_date.(string)
	}
	if assignees, ok := d.GetOk("assignees"); ok {
		task["assignee_ids"] = assignees.([]interface{})
	}
	if parentTaskId, ok := d.GetOk("parent_task_id"); ok {
		task["parent_task_id"] = parentTaskId.(int)
	} else {
		task["parent_task_id"] = 0
	}
	if labels, ok := d.GetOk("labels"); ok {
		task["label_ids"] = labels.([]interface{})
	}

	var created map[string]interface{}
	if err := client.Post("api/tasks", task, &created); err != nil {
		return diag.FromErr(err)
	}

	result, ok := created["task"].(map[string]interface{})
	if !ok {
		return diag.Errorf("unable to create task")
	}

	id := fmt.Sprintf("%v", result["ID"])
	d.SetId(id)

	return resourceReadTask(ctx, d, m)

}

func resourceUpdateTask(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	task := map[string]interface{}{
		"title": d.Get("title").(string),
	}
	if description, ok := d.GetOk("description"); ok {
		task["description"] = description.(string)
	}
	if status, ok := d.GetOk("status"); ok {
		task["status"] = status.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		task["priority"] = priority.(string)
	}
	if due_date, ok := d.GetOk("due_date"); ok {
		task["due_date"] = due_date.(string)
	}
	if creatorId, ok := d.GetOk("creator_id"); ok {
		task["creator_id"] = creatorId.(int)
	}

	var updated map[string]interface{}
	if err := client.Put("api/tasks/"+d.Id(), task, &updated); err != nil {
		return diag.FromErr(err)
	}

	log.Println("some details got updated")

	task = map[string]interface{}{}
	task["team_id"] = d.Get("team_id").(int)
	updated = map[string]interface{}{}
	if err := client.Put("api/tasks/"+d.Id()+"/change-team", task, &updated); err != nil {
		return diag.FromErr(err)
	}

	log.Println("teamId got updated")

	task = map[string]interface{}{}
	if assignees, ok := d.GetOk("assignees"); ok {
		task["assignees"] = assignees.([]interface{})
	}
	updated = map[string]interface{}{}
	if err := client.Put("api/tasks/"+d.Id()+"/add-assignee", task, &updated); err != nil {
		return diag.FromErr(err)
	}

	log.Println("parentId got updated")

	task = map[string]interface{}{}
	var parentTaskID int
	if parentTaskId, ok := d.GetOk("parent_task_id"); ok {
		task["parent_task_id"] = parentTaskId.(int)
		parentTaskID = parentTaskId.(int)
	}
	if parentTaskID > 0 {
		updated = map[string]interface{}{}
		if err := client.Put("api/tasks/"+d.Id()+"/parent-id", task, &updated); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Println("labels got updated")

	task = map[string]interface{}{}
	if labels, ok := d.GetOk("labels"); ok {
		task["labels"] = labels.([]interface{})
	}
	updated = map[string]interface{}{}
	if err := client.Put("api/tasks/"+d.Id()+"/add-labels", task, &updated); err != nil {
		return diag.FromErr(err)
	}

	return resourceReadTask(ctx, d, m)
}

func resourceReadTask(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	var task map[string]interface{}
	if err := client.Get("api/tasks/"+d.Id(), &task); err != nil {
		return diag.FromErr(err)
	}

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

	var assigneeIDs []int
	if assigneesRaw, ok := result["assignees"].([]interface{}); ok {
		for _, assignee := range assigneesRaw {
			if assigneeMap, ok := assignee.(map[string]interface{}); ok {
				if idFloat, ok := assigneeMap["ID"].(float64); ok {
					assigneeIDs = append(assigneeIDs, int(idFloat))
				}
			}
		}
	}
	d.Set("assignees", assigneeIDs)

	d.Set("parent_task_id", result["parent_task_id"])

	var subTaskIDs []int
	if subTaskRaw, ok := result["subtasks"].([]interface{}); ok {
		for _, subTask := range subTaskRaw {
			if subTaskMap, ok := subTask.(map[string]interface{}); ok {
				if idFloat, ok := subTaskMap["ID"].(float64); ok {
					subTaskIDs = append(subTaskIDs, int(idFloat))
				}
			}
		}
	}
	d.Set("subtasks", subTaskIDs)

	var labelIDs []int
	if labelRaw, ok := result["labels"].([]interface{}); ok {
		for _, label := range labelRaw {
			if labelMap, ok := label.(map[string]interface{}); ok {
				if idFloat, ok := labelMap["ID"].(float64); ok {
					labelIDs = append(labelIDs, int(idFloat))
				}
			}
		}
	}
	d.Set("labels", labelIDs)

	var commentIDs []int
	if commentRaw, ok := result["comments"].([]interface{}); ok {
		for _, comment := range commentRaw {
			if commentMap, ok := comment.(map[string]interface{}); ok {
				if idFloat, ok := commentMap["ID"].(float64); ok {
					commentIDs = append(commentIDs, int(idFloat))
				}
			}
		}
	}
	d.Set("comments", commentIDs)

	var attachmantIDs []int
	if attachmantRaw, ok := result["attachments"].([]interface{}); ok {
		for _, attachmant := range attachmantRaw {
			if attachmantMap, ok := attachmant.(map[string]interface{}); ok {
				if idFloat, ok := attachmantMap["ID"].(float64); ok {
					attachmantIDs = append(attachmantIDs, int(idFloat))
				}
			}
		}
	}
	d.Set("attachments", attachmantIDs)

	return nil
}

func resourceDeleteTask(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	if err := client.Delete("api/tasks/" + d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
