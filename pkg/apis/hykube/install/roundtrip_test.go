/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/

package install

import (
	"testing"

	hykubefuzzer "github.com/infrapot/hykube/pkg/apis/hykube/fuzzer"
	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForAPIGroup(t, Install, hykubefuzzer.Funcs)
	// TODO: enable protobuf generation for the sample-apiserver
	// roundtrip.RoundTripProtobufTestForAPIGroup(t, Install, hykubefuzzer.Funcs)
}
