package taskmanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAttachment,
		ReadContext:   resourceReadAttachment,
		UpdateContext: resourceUpdateAttachment,
		DeleteContext: resourceDeleteAttachment,
		Schema: map[string]*schema.Schema{
			"file_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"uploader_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceCreateAttachment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	taskID := d.Get("task_id").(int)
	filePath := d.Get("url").(string)

	absPath, err := filepath.Abs(filepath.Clean(filePath))
	if err != nil {
		return diag.FromErr(err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return diag.FromErr(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(absPath))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create form part: %w", err))
	}
	if _, err := io.Copy(part, file); err != nil {
		return diag.FromErr(fmt.Errorf("failed to copy file: %w", err))
	}
	writer.Close()

	endpoint := fmt.Sprintf("api/tasks/%d/attachments", taskID)
	req, err := http.NewRequest("POST", client.baseURL+endpoint, body)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", "Bearer "+client.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return diag.FromErr(fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody)))
	}

	log.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$ before unmarshal")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	log.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$ after unmarshal")

	if attachment, ok := result["attachment"].(map[string]interface{}); ok {
		if id, ok := attachment["ID"].(float64); ok {
			d.SetId(fmt.Sprintf("%d", int(id)))
		} else {
			return diag.FromErr(fmt.Errorf("invalid response: missing 'id'"))
		}
	} else {
		return diag.FromErr(fmt.Errorf("invalid response structure"))
	}

	var url interface{} = filePath
	d.Set("url", url)

	return resourceReadAttachment(ctx, d, m)
}

func resourceUpdateAttachment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceDeleteAttachment(ctx, d, m)

	resourceCreateAttachment(ctx, d, m)
	return nil
}

func resourceReadAttachment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	var attachment map[string]interface{}

	if err := client.Get("api/attachments/"+d.Id(), &attachment); err != nil {
		return diag.FromErr(err)
	}

	result := attachment["attachment"].(map[string]interface{})

	d.Set("file_name", result["file_name"])
	d.Set("task_id", result["task_id"])
	if uploaderID, ok := result["uploader_id"].(float64); ok {
		d.Set("uploader_id", int(uploaderID))
	}

	return nil
}

func resourceDeleteAttachment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*TaskManagerClient)

	if err := client.Delete("api/attachments/" + d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
