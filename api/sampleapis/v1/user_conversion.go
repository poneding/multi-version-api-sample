/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// +kubebuilder:docs-gen:collapse=Apache License

package v1

/*
For imports, we'll need the controller-runtime
[`conversion`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/conversion?tab=doc)
package, plus the API version for our hub type (v1), and finally some of the
standard packages.
*/
import (
	v2 "github.com/poneding/multi-version-api-sample/api/sampleapis/v2"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *User) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v2.User)

	dst.ObjectMeta = src.ObjectMeta

	parts := strings.SplitN(src.Spec.FullName, " ", 2)
	dst.Spec.FirstName = parts[0]
	if len(parts) > 1 {
		dst.Spec.LastName = parts[1]
	}

	return nil
}

func (dst *User) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v2.User)

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.FullName = src.Spec.FirstName + " " + src.Spec.LastName
	return nil
}
