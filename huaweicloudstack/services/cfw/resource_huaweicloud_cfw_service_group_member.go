// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CFW
// ---------------------------------------------------------------

package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func ResourceServiceGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceGroupMemberCreate,
		ReadContext:   resourceServiceGroupMemberRead,
		DeleteContext: resourceServiceGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceServiceGroupMemberImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the service group.`,
			},
			"protocol": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the protocol type.`,
			},
			"source_port": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the source port.source_port`,
			},
			"dest_port": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the destination port.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the service member name`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the service member description.`,
			},
		},
	}
}

func resourceServiceGroupMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createServiceGroupMember: Create a CFW service group member.
	var (
		createServiceGroupMemberHttpUrl = "v1/{project_id}/service-items"
		createServiceGroupMemberProduct = "cfw"
	)
	createServiceGroupMemberClient, err := cfg.NewServiceClient(createServiceGroupMemberProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	createServiceGroupMemberPath := createServiceGroupMemberClient.Endpoint + createServiceGroupMemberHttpUrl
	createServiceGroupMemberPath = strings.ReplaceAll(createServiceGroupMemberPath, "{project_id}",
		createServiceGroupMemberClient.ProjectID)

	createServiceGroupMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createServiceGroupMemberOpt.JSONBody = utils.RemoveNil(buildCreateServiceGroupMemberBodyParams(d))
	createServiceGroupMemberResp, err := createServiceGroupMemberClient.Request("POST", createServiceGroupMemberPath,
		&createServiceGroupMemberOpt)
	if err != nil {
		return diag.Errorf("error creating ServiceGroupMember: %s", err)
	}

	createServiceGroupMemberRespBody, err := utils.FlattenResponse(createServiceGroupMemberResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("data.items[0].id", createServiceGroupMemberRespBody)
	if err != nil {
		return diag.Errorf("error creating ServiceGroupMember: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceServiceGroupMemberRead(ctx, d, meta)
}

func buildCreateServiceGroupMemberBodyParams(d *schema.ResourceData) map[string]interface{} {
	item := map[string]interface{}{
		"protocol":    utils.ValueIngoreEmpty(d.Get("protocol")),
		"source_port": utils.ValueIngoreEmpty(d.Get("source_port")),
		"dest_port":   utils.ValueIngoreEmpty(d.Get("dest_port")),
		"name":        utils.ValueIngoreEmpty(d.Get("name")),
		"description": utils.ValueIngoreEmpty(d.Get("description")),
	}

	bodyParams := map[string]interface{}{
		"set_id":        utils.ValueIngoreEmpty(d.Get("group_id")),
		"service_items": []map[string]interface{}{item},
	}
	return bodyParams
}

func resourceServiceGroupMemberRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getServiceGroupMember: Query the CFW service group member detail
	var (
		getServiceGroupMemberHttpUrl = "v1/{project_id}/service-items"
		getServiceGroupMemberProduct = "cfw"
	)
	getServiceGroupMemberClient, err := cfg.NewServiceClient(getServiceGroupMemberProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	getServiceGroupMemberPath := getServiceGroupMemberClient.Endpoint + getServiceGroupMemberHttpUrl
	getServiceGroupMemberPath = strings.ReplaceAll(getServiceGroupMemberPath, "{project_id}",
		getServiceGroupMemberClient.ProjectID)

	getServiceGroupMemberqueryParams := buildGetServiceGroupMemberQueryParams(d)
	getServiceGroupMemberPath += getServiceGroupMemberqueryParams

	getServiceGroupMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getServiceGroupMemberResp, err := getServiceGroupMemberClient.Request("GET", getServiceGroupMemberPath,
		&getServiceGroupMemberOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceGroupMember")
	}

	getServiceGroupMemberRespBody, err := utils.FlattenResponse(getServiceGroupMemberResp)
	if err != nil {
		return diag.FromErr(err)
	}

	members, err := jmespath.Search("data.records", getServiceGroupMemberRespBody)
	if err != nil {
		diag.Errorf("error parsing data.records from response= %#v", getServiceGroupMemberRespBody)
	}

	member, err := FilterServiceGroupMembers(members.([]interface{}), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceGroupMember")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("group_id", utils.PathSearch("data.set_id", getServiceGroupMemberRespBody, nil)),
		d.Set("protocol", utils.PathSearch("protocol", member, nil)),
		d.Set("source_port", utils.PathSearch("source_port", member, nil)),
		d.Set("dest_port", utils.PathSearch("dest_port", member, nil)),
		d.Set("name", utils.PathSearch("name", member, nil)),
		d.Set("description", utils.PathSearch("description", member, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func FilterServiceGroupMembers(members []interface{}, id string) (interface{}, error) {
	if len(members) != 0 {
		for _, v := range members {
			member := v.(map[string]interface{})
			if member["item_id"] == id {
				return v, nil
			}
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func buildGetServiceGroupMemberQueryParams(d *schema.ResourceData) string {
	res := "?offset=0&limit=1024"
	res = fmt.Sprintf("%s&set_id=%v", res, d.Get("group_id"))

	return res
}

func resourceServiceGroupMemberDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteServiceGroupMember: Delete an existing CFW service group member
	var (
		deleteServiceGroupMemberHttpUrl = "v1/{project_id}/service-items/{id}"
		deleteServiceGroupMemberProduct = "cfw"
	)
	deleteServiceGroupMemberClient, err := cfg.NewServiceClient(deleteServiceGroupMemberProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	deleteServiceGroupMemberPath := deleteServiceGroupMemberClient.Endpoint + deleteServiceGroupMemberHttpUrl
	deleteServiceGroupMemberPath = strings.ReplaceAll(deleteServiceGroupMemberPath, "{project_id}",
		deleteServiceGroupMemberClient.ProjectID)
	deleteServiceGroupMemberPath = strings.ReplaceAll(deleteServiceGroupMemberPath, "{id}", d.Id())

	deleteServiceGroupMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteServiceGroupMemberClient.Request("DELETE", deleteServiceGroupMemberPath,
		&deleteServiceGroupMemberOpt)
	if err != nil {
		return diag.Errorf("error deleting ServiceGroupMember: %s", err)
	}

	return nil
}

func resourceServiceGroupMemberImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <group_id>/<member_id>")
	}

	d.Set("group_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
