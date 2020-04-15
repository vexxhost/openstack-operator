// Copyright 2020 VEXXHOST, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8sutils

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// CreateOrUpdate wraps the function provided by controller-runtime to include
// some additional logging and common functionality across all resources.
func CreateOrUpdate(ctx context.Context, c client.Client, obj runtime.Object, f controllerutil.MutateFn) (controllerutil.OperationResult, error) {

	return controllerutil.CreateOrUpdate(ctx, c, obj, func() error {
		original := obj.DeepCopyObject()

		err := f()
		if err != nil {
			return err
		}

		generateObjectDiff(original, obj)
		return nil
	})
}

func generateObjectDiff(original runtime.Object, modified runtime.Object) {
	diff := cmp.Diff(original, modified)

	if len(diff) != 0 {
		fmt.Println(diff)
	}
}
