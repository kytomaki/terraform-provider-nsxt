/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceDhcpRelayService() *schema.Resource {
	return &schema.Resource{
		Create: resourceDhcpRelayServiceCreate,
		Read:   resourceDhcpRelayServiceRead,
		Update: resourceDhcpRelayServiceUpdate,
		Delete: resourceDhcpRelayServiceDelete,

		Schema: map[string]*schema.Schema{
			"revision":     getRevisionSchema(),
			"system_owned": getSystemOwnedSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to ID if not set",
				Optional:    true,
			},
			"tags": getTagsSchema(),
			"dhcp_relay_profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "dhcp relay profile referenced by the dhcp relay service",
				Required:    true,
			},
		},
	}
}

func resourceDhcpRelayServiceCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	dhcp_relay_profile_id := d.Get("dhcp_relay_profile_id").(string)
	dhcp_relay_service := manager.DhcpRelayService{
		Description:        description,
		DisplayName:        display_name,
		Tags:               tags,
		DhcpRelayProfileId: dhcp_relay_profile_id,
	}

	dhcp_relay_service, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateDhcpRelay(nsxClient.Context, dhcp_relay_service)

	if err != nil {
		return fmt.Errorf("Error during DhcpRelayService create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during DhcpRelayService create: %v", resp.StatusCode)
	}
	d.SetId(dhcp_relay_service.Id)

	return resourceDhcpRelayServiceRead(d, m)
}

func resourceDhcpRelayServiceRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining dhcp relay service id")
	}

	dhcp_relay_service, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadDhcpRelay(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("DhcpRelayService %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during DhcpRelayService read: %v", err)
	}

	d.Set("Revision", dhcp_relay_service.Revision)
	d.Set("SystemOwned", dhcp_relay_service.SystemOwned)
	d.Set("Description", dhcp_relay_service.Description)
	d.Set("DisplayName", dhcp_relay_service.DisplayName)
	setTagsInSchema(d, dhcp_relay_service.Tags)
	d.Set("DhcpRelayProfileId", dhcp_relay_service.DhcpRelayProfileId)

	return nil
}

func resourceDhcpRelayServiceUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining dhcp relay service id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	dhcp_relay_profile_id := d.Get("dhcp_relay_profile_id").(string)
	dhcp_relay_service := manager.DhcpRelayService{
		Revision:           revision,
		Description:        description,
		DisplayName:        display_name,
		Tags:               tags,
		DhcpRelayProfileId: dhcp_relay_profile_id,
	}

	dhcp_relay_service, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateDhcpRelay(nsxClient.Context, id, dhcp_relay_service)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during DhcpRelayService update: %v", err)
	}

	return resourceDhcpRelayServiceRead(d, m)
}

func resourceDhcpRelayServiceDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining dhcp relay service id")
	}

	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteDhcpRelay(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during DhcpRelayService delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("DhcpRelayService %s not found", id)
		d.SetId("")
	}
	return nil
}