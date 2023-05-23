/*
Copyright (c***REMOVED*** 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License"***REMOVED***;
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

***REMOVED***
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkg/errors"
***REMOVED***

const versionPrefix = "openshift-v"

// shouldPatchInt changed checks if the change between the given state and plan requires sending a
// patch request to the server. If it does it returns the value to add to the patch.
func ShouldPatchInt(state, plan types.Int64***REMOVED*** (value int64, ok bool***REMOVED*** {
	if plan.Unknown || plan.Null {
		return
	}
	if state.Unknown || state.Null {
		value = plan.Value
		ok = true
		return
	}
	if plan.Value != state.Value {
		value = plan.Value
		ok = true
	}
	return
}

// shouldPatchString changed checks if the change between the given state and plan requires sending
// a patch request to the server. If it does it returns the value to add to the patch.
func ShouldPatchString(state, plan types.String***REMOVED*** (value string, ok bool***REMOVED*** {
	if plan.Unknown || plan.Null {
		return
	}
	if state.Unknown || state.Null {
		value = plan.Value
		ok = true
		return
	}
	if plan.Value != state.Value {
		value = plan.Value
		ok = true
	}
	return
}

// ShouldPatchBool changed checks if the change between the given state and plan requires sending
// a patch request to the server. If it does it return the value to add to the patch.
func ShouldPatchBool(state, plan types.Bool***REMOVED*** (value bool, ok bool***REMOVED*** {
	if plan.Unknown || plan.Null {
		return
	}
	if state.Unknown || state.Null {
		value = plan.Value
		ok = true
		return
	}
	if plan.Value != state.Value {
		value = plan.Value
		ok = true
	}
	return
}

// ShouldPatchMap changed checks if the change between the given state and plan requires sending
// a patch request to the server. If it does it return the value to add to the patch.
func ShouldPatchMap(state, plan types.Map***REMOVED*** (types.Map, bool***REMOVED*** {
	return plan, !reflect.DeepEqual(state.Elems, plan.Elems***REMOVED***
}

// TF types converter functions
func StringArrayToList(arr []string***REMOVED*** types.List {
	list := types.List{
		ElemType: types.StringType,
		Elems:    []attr.Value{},
	}

	for _, elm := range arr {
		list.Elems = append(list.Elems, types.String{Value: elm}***REMOVED***
	}

	return list
}

func StringListToArray(list types.List***REMOVED*** ([]string, error***REMOVED*** {
	arr := []string{}
	for _, elm := range list.Elems {
		stype, ok := elm.(types.String***REMOVED***
		if !ok {
			return arr, errors.New("Failed to convert TF list to string slice."***REMOVED***
***REMOVED***
		arr = append(arr, stype.Value***REMOVED***
	}
	return arr, nil
}

func IsValidDomain(candidate string***REMOVED*** bool {
	var domainRegexp = regexp.MustCompile(`^(?i***REMOVED***[a-z0-9-]+(\.[a-z0-9-]+***REMOVED***+\.?$`***REMOVED***
	return domainRegexp.MatchString(candidate***REMOVED***
}

func IsStringAttributeEmpty(param types.String***REMOVED*** bool {
	return param.Unknown || param.Null || param.Value == ""
}

func IsGreaterThanOrEqual(version1, version2 string***REMOVED*** (bool, error***REMOVED*** {
	v1, err := version.NewVersion(strings.TrimPrefix(version1, versionPrefix***REMOVED******REMOVED***
	if err != nil {
		return false, err
	}
	v2, err := version.NewVersion(strings.TrimPrefix(version2, versionPrefix***REMOVED******REMOVED***
	if err != nil {
		return false, err
	}
	return v1.GreaterThanOrEqual(v2***REMOVED***, nil
}
