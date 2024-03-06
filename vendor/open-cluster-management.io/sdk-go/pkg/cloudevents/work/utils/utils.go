package utils

import (
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/types"
	workv1 "open-cluster-management.io/api/work/v1"
	"open-cluster-management.io/sdk-go/pkg/cloudevents/work/common"
)

// Patch applies the patch to a work with the patch type.
func Patch(patchType types.PatchType, work *workv1.ManifestWork, patchData []byte) (*workv1.ManifestWork, error) {
	workData, err := json.Marshal(work)
	if err != nil {
		return nil, err
	}

	var patchedData []byte
	switch patchType {
	case types.JSONPatchType:
		var patchObj jsonpatch.Patch
		patchObj, err = jsonpatch.DecodePatch(patchData)
		if err != nil {
			return nil, err
		}
		patchedData, err = patchObj.Apply(workData)
		if err != nil {
			return nil, err
		}

	case types.MergePatchType:
		patchedData, err = jsonpatch.MergePatch(workData, patchData)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported patch type: %s", patchType)
	}

	patchedWork := &workv1.ManifestWork{}
	if err := json.Unmarshal(patchedData, patchedWork); err != nil {
		return nil, err
	}

	return patchedWork, nil
}

// UID returns a v5 UUID based on sourceID, work name and namespace to make sure it is consistent
func UID(sourceID, namespace, name string) string {
	id := fmt.Sprintf("%s-%s-%s-%s", sourceID, common.ManifestWorkGR.String(), namespace, name)
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(id)).String()
}
