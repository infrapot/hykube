/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/

package fuzzer

import (
	fuzz "github.com/google/gofuzz"
	"github.com/infrapot/hykube/pkg/apis/hykube"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

// Funcs returns the fuzzer functions for the apps api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *hykube.ProviderSpec, c fuzz.Continue) {
			c.FuzzNoCustom(s) // fuzz self without calling this function again

			if len(s.ProviderReference) != 0 {
				s.ProviderReference = ""
			}
			if len(s.ProviderReference) != 0 {
				s.ReferenceType = hykube.ProviderReferenceType
			} else {
				s.ReferenceType = ""
			}
		},
	}
}
