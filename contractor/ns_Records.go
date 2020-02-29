/*Package contractor(version: "0.1") - Automatically generated by cinp-codegen from /api/v1/Records/ at 2020-02-28T18:22:28.464196
 */
package contractor

import (
	cinp "github.com/cinp/go"
	"reflect"
)

//RecordsRecorder - Model Recorder(/api/v1/Records/Recorder)
/*

 */
type RecordsRecorder struct {
	cinp.BaseObject
	cinp *cinp.CInP
}

// AsMap returns a map[string]interface{} that is required for create and update
func (object *RecordsRecorder) AsMap(isCreate bool) *map[string]interface{} {
	if isCreate {
		return &map[string]interface{}{ 
		}
	}
	return &map[string]interface{}{ 
	}
}

// RecordsRecorderNew - Make a new object of Model Recorder
func (service *Contractor) RecordsRecorderNew() *RecordsRecorder {
	return &RecordsRecorder{cinp: service.cinp}
}

// RecordsRecorderNewWithID - Make a new object of Model Recorder
func (service *Contractor) RecordsRecorderNewWithID(id string) *RecordsRecorder {
	result := RecordsRecorder{cinp: service.cinp}
	result.SetID("/api/v1/Records/Recorder:"+id+":")
	return &result
}

// RecordsRecorderGet - Get function for Model Recorder
func (service *Contractor) RecordsRecorderGet(id string) (*RecordsRecorder, error) {
	object, err := service.cinp.Get("/api/v1/Records/Recorder:"+id+":")
	if err != nil {
		return nil, err
	}
	result := (*object).(*RecordsRecorder)
	result.cinp = service.cinp

	return result, nil
}

// Create - Create function for Model Recorder
func (object *RecordsRecorder) Create() error {
	if err := object.cinp.Create("/api/v1/Records/Recorder", object); err != nil {
		return err
	}

	return nil
}

// Update - Update function for Model Recorder
func (object *RecordsRecorder) Update(fieldList []string) error {
	if err := object.cinp.Update(object, fieldList); err != nil {
		return err
	}

	return nil
}

// Delete - Delete function for Model Recorder
func (object *RecordsRecorder) Delete() error {
	if err := object.cinp.Delete(object); err != nil {
		return err
	}

	return nil
}

// RecordsRecorderList - List function for Model Recorder
func (service *Contractor) RecordsRecorderList(filterName string, filterValues map[string]interface{}) <-chan *RecordsRecorder {
	in := service.cinp.ListObjects("/api/v1/Records/Recorder", reflect.TypeOf(RecordsRecorder{}), filterName, filterValues, 50)
	out := make(chan *RecordsRecorder)
	go func() {
		defer close(out)
		for v := range in {
			out <- v.(*RecordsRecorder)
		}
	}()
	return out
}

// RecordsRecorderCallQuery calls queryNoneNoneNoneNone
func (service *Contractor) RecordsRecorderCallQuery(group string, query string, fields string, max_results int) ([]string, error) {
	args := map[string]interface{}{
		"group": group,
		"query": query,
		"fields": fields,
		"max_results": max_results,
	}
	uri := "/api/v1/Records/Recorder(query)"

	result := []string{}

	if err := service.cinp.Call(uri, &args, &result); err != nil {
		return []string{}, err
	}

	return result, nil
}

// RecordsRecorderCallQuery_objects calls query_objectsNoneNoneNone
func (service *Contractor) RecordsRecorderCallQuery_objects(group string, query string, max_results int) (string, error) {
	args := map[string]interface{}{
		"group": group,
		"query": query,
		"max_results": max_results,
	}
	uri := "/api/v1/Records/Recorder(query_objects)"

	result := ""

	if err := service.cinp.Call(uri, &args, &result); err != nil {
		return "", err
	}

	return result, nil
}

func registerRecords(cinp *cinp.CInP) { 
	cinp.RegisterType("/api/v1/Records/Recorder", reflect.TypeOf((*RecordsRecorder)(nil)).Elem())
}