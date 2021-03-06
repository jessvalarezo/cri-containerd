/*
Copyright 2017 The Kubernetes Authors.

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

package server

import (
	"testing"

	snapshot "github.com/containerd/containerd/snapshots"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"

	snapshotstore "github.com/kubernetes-incubator/cri-containerd/pkg/store/snapshot"
)

func TestImageFsInfo(t *testing.T) {
	c := newTestCRIContainerdService()
	snapshots := []snapshotstore.Snapshot{
		{
			Key:       "key1",
			Kind:      snapshot.KindActive,
			Size:      10,
			Inodes:    100,
			Timestamp: 234567,
		},
		{
			Key:       "key2",
			Kind:      snapshot.KindCommitted,
			Size:      20,
			Inodes:    200,
			Timestamp: 123456,
		},
		{
			Key:       "key3",
			Kind:      snapshot.KindView,
			Size:      0,
			Inodes:    0,
			Timestamp: 345678,
		},
	}
	expected := &runtime.FilesystemUsage{
		Timestamp:  123456,
		StorageId:  &runtime.StorageIdentifier{Uuid: testImageFSUUID},
		UsedBytes:  &runtime.UInt64Value{Value: 30},
		InodesUsed: &runtime.UInt64Value{Value: 300},
	}
	for _, sn := range snapshots {
		c.snapshotStore.Add(sn)
	}
	resp, err := c.ImageFsInfo(context.Background(), &runtime.ImageFsInfoRequest{})
	require.NoError(t, err)
	stats := resp.GetImageFilesystems()
	assert.Len(t, stats, 1)
	assert.Equal(t, expected, stats[0])
}
