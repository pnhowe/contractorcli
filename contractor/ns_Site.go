/*Package contractor(version: "0.1") - Automatically generated by cinp-codegen from /api/v1/Site/ at 2020-02-28T18:22:28.464196
 */
package contractor

import (
	"time"
	cinp "github.com/cinp/go"
	"reflect"
)

//SiteSite - Model Site(/api/v1/Site/Site)
/*
Site(name, zone, description, parent, config_values, updated, created)
 */
type SiteSite struct {
	cinp.BaseObject
	cinp *cinp.CInP
	Name string `json:"name"`
	Zone string `json:"zone"`
	Description string `json:"description"`
	Parent string `json:"parent"`
	ConfigValues map[string]interface{} `json:"config_values"`
	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

// AsMap returns a map[string]interface{} that is required for create and update
func (object *SiteSite) AsMap(isCreate bool) *map[string]interface{} {
	if isCreate {
		return &map[string]interface{}{ 
			"name": object.Name,
			"zone": object.Zone,
			"description": object.Description,
			"parent": object.Parent,
			"config_values": object.ConfigValues,
		}
	}
	return &map[string]interface{}{ 
		"zone": object.Zone,
		"description": object.Description,
		"parent": object.Parent,
		"config_values": object.ConfigValues,
	}
}

// SiteSiteNew - Make a new object of Model Site
func (service *Contractor) SiteSiteNew() *SiteSite {
	return &SiteSite{cinp: service.cinp}
}

// SiteSiteNewWithID - Make a new object of Model Site
func (service *Contractor) SiteSiteNewWithID(id string) *SiteSite {
	result := SiteSite{cinp: service.cinp}
	result.SetID("/api/v1/Site/Site:"+id+":")
	return &result
}

// SiteSiteGet - Get function for Model Site
func (service *Contractor) SiteSiteGet(id string) (*SiteSite, error) {
	object, err := service.cinp.Get("/api/v1/Site/Site:"+id+":")
	if err != nil {
		return nil, err
	}
	result := (*object).(*SiteSite)
	result.cinp = service.cinp

	return result, nil
}

// Create - Create function for Model Site
func (object *SiteSite) Create() error {
	if err := object.cinp.Create("/api/v1/Site/Site", object); err != nil {
		return err
	}

	return nil
}

// Update - Update function for Model Site
func (object *SiteSite) Update(fieldList []string) error {
	if err := object.cinp.Update(object, fieldList); err != nil {
		return err
	}

	return nil
}

// Delete - Delete function for Model Site
func (object *SiteSite) Delete() error {
	if err := object.cinp.Delete(object); err != nil {
		return err
	}

	return nil
}

// SiteSiteList - List function for Model Site
func (service *Contractor) SiteSiteList(filterName string, filterValues map[string]interface{}) <-chan *SiteSite {
	in := service.cinp.ListObjects("/api/v1/Site/Site", reflect.TypeOf(SiteSite{}), filterName, filterValues, 50)
	out := make(chan *SiteSite)
	go func() {
		defer close(out)
		for v := range in {
			out <- v.(*SiteSite)
		}
	}()
	return out
}

// CallGetConfig calls getConfig
func (object *SiteSite) CallGetConfig() (map[string]interface{}, error) {
	args := map[string]interface{}{
	}
	_, _, _, ids, _, err := object.cinp.Split(object.GetID())
	if err != nil {
		return nil, err
	}
	uri, err := object.cinp.UpdateIDs("/api/v1/Site/Site(getConfig)", ids)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}

	if err := object.cinp.Call(uri, &args, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CallGetDependencyMap calls getDependencyMap
func (object *SiteSite) CallGetDependencyMap() (map[string]interface{}, error) {
	args := map[string]interface{}{
	}
	_, _, _, ids, _, err := object.cinp.Split(object.GetID())
	if err != nil {
		return nil, err
	}
	uri, err := object.cinp.UpdateIDs("/api/v1/Site/Site(getDependencyMap)", ids)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}

	if err := object.cinp.Call(uri, &args, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func registerSite(cinp *cinp.CInP) { 
	cinp.RegisterType("/api/v1/Site/Site", reflect.TypeOf((*SiteSite)(nil)).Elem())
}