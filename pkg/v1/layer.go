// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/google/go-containerregistry/pkg/v1/types"
)

// Layer is an interface for accessing the properties of a particular layer of a v1.Image
type Layer interface {
	// Digest returns the Hash of the compressed layer.
	Digest() (Hash, error)

	// DiffID returns the Hash of the uncompressed layer.
	DiffID() (Hash, error)

	// Compressed returns an io.ReadCloser for the compressed layer contents.
	Compressed() (io.ReadCloser, error)

	// Uncompressed returns an io.ReadCloser for the uncompressed layer contents.
	Uncompressed() (io.ReadCloser, error)

	// Size returns the compressed size of the Layer.
	Size() (int64, error)

	// MediaType returns the media type of the Layer.
	MediaType() (types.MediaType, error)
}

// StaticLayer returns a Layer which references a static payload.
// `Compressed` and `Uncompressed` are equivalent, returning a Reader which returns the raw contents.
// `Digest` and `DiffID` are similarly equivalent, returning the SHA256 Hash of the raw contents.
func StaticLayer(contents []byte, mediaType types.MediaType) (Layer, error) {
	h, _, err := SHA256(bytes.NewReader(contents))
	if err != nil {
		return nil, err
	}
	return &staticLayer{
		b:  contents,
		h:  h,
		mt: mediaType,
	}, nil
}

type staticLayer struct {
	b  []byte
	h  Hash
	mt types.MediaType
}

func (l *staticLayer) Digest() (Hash, error) {
	return l.h, nil
}

// DiffID returns the Hash of the uncompressed layer.
func (l *staticLayer) DiffID() (Hash, error) {
	return l.h, nil
}

// Compressed returns an io.ReadCloser for the compressed layer contents.
func (l *staticLayer) Compressed() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(l.b)), nil
}

// Uncompressed returns an io.ReadCloser for the uncompressed layer contents.
func (l *staticLayer) Uncompressed() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(l.b)), nil
}

// Size returns the compressed size of the Layer.
func (l *staticLayer) Size() (int64, error) {
	return int64(len(l.b)), nil
}

// MediaType returns the media type of the Layer.
func (l *staticLayer) MediaType() (types.MediaType, error) {
	return l.mt, nil
}
