/*
Copyright 2026 Thomas Boerger <thomas@webhippie.de>.

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

package controller

import (
	"context"

	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// RefIndexKey builds the field-index value used to correlate the
// namespace/name of a referenced Secret or ConfigMap with the CRs that
// reference it. defaultNamespace is applied by callers when the reference
// itself does not specify an explicit namespace.
func RefIndexKey(namespace, name string) string {
	return namespace + "/" + name
}

// RegisterRefIndex registers a field index named indexField on obj's kind,
// populated by keysFunc, so that instances referencing a given Secret or
// ConfigMap (identified via RefIndexKey) can be looked up efficiently. It
// must be called once per CR kind during SetupWithManager, before the
// manager starts.
func RegisterRefIndex(mgr ctrl.Manager, obj client.Object, indexField string, keysFunc client.IndexerFunc) error {
	return mgr.GetFieldIndexer().IndexField(context.Background(), obj, indexField, keysFunc)
}

// RefEventHandler returns an event handler that, given a changed Secret or
// ConfigMap, lists newList's kind via the indexField index and enqueues a
// reconcile request for every matching CR. newList must return a fresh, empty
// list of the target CR kind on every call.
func RefEventHandler(cl client.Client, newList func() client.ObjectList, indexField string) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
		key := RefIndexKey(obj.GetNamespace(), obj.GetName())

		list := newList()
		if err := cl.List(ctx, list, client.MatchingFields{indexField: key}); err != nil {
			return nil
		}

		items, err := apimeta.ExtractList(list)
		if err != nil {
			return nil
		}

		requests := make([]reconcile.Request, 0, len(items))
		for _, item := range items {
			o, ok := item.(client.Object)
			if !ok {
				continue
			}
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{Namespace: o.GetNamespace(), Name: o.GetName()},
			})
		}

		return requests
	})
}
