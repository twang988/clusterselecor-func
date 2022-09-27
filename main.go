package main

import (
	"os"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Run(rl *fn.ResourceList) (bool, error) {

	var pd PackageDeployment
	for _, obj := range rl.Items {
		if obj.IsGVK("automation.nephio.org", "v1alpha1", "PackageDeployment") {
			obj.As(&pd)
		}

	}
	if pd.Spec.Selector.MatchLabels == nil {
		return true, nil
	}
	//Select Cluster by label
	var err error
	items, err := clusterFilter(rl.Items, pd.Spec.Selector.MatchLabels)
	if err != nil {
		return false, err
	}
	rl.Items = items

	return true, nil
}
func clusterFilter(objs []*fn.KubeObject, labels map[string]string) ([]*fn.KubeObject, error) {
	var newobjs []*fn.KubeObject
	for _, obj := range objs {
		if obj.IsGVK("infra.nephio.org", "v1alpha1", "Cluster") &&
			obj.HasLabels(labels) {
			newobjs = append(newobjs, obj)
		}
	}

	return newobjs, nil
}

// PackageDeployment is the Schema for the packagedeployments API
type PackageDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PackageDeploymentSpec `json:"spec,omitempty"`
}

// PackageDeploymentSpec defines the desired state of PackageDeployment
type PackageDeploymentSpec struct {
	// Label selector for Clusters on which to deploy the package
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// PackageRef identifies the package revision to deploy
	PackageRef PackageRevisionReference `json:"packageRef"`

	// Namespace identifies the namespace in which to deploy the package
	// The namespace will be added to the resource list of the package
	// If not present, the package will be installed in the default namespace
	Namespace *string `json:"namespace,omitempty"`

	//internalFunctions
	InternalFunctions InternalFunction `json:"internalFunctions,omitempty"`
}

// Internal funcion for clusterselector.
// can Define more for further general needs
type InternalFunction struct {
	ClusterSelector           string `json:"clusterSelector,omitempty"`
	AnotherFuncForOtherThings string `json:"anotherFuncForOtherThings,omitempty"`
}

// PackageRevisionReference is used to reference a particular package revision.
type PackageRevisionReference struct {
	// Namespace is the namespace for both the repository and package revision
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Repository is the name of the repository containing the package
	RepositoryName string `json:"repository"`

	// PackageName is the name of the package for the revision
	PackageName string `json:"packageName"`

	// Revision is the specific version number of the revision of the package
	Revision string `json:"revision"`
}

func main() {
	if err := fn.AsMain(fn.ResourceListProcessorFunc(Run)); err != nil {
		os.Exit(1)
	}
}
