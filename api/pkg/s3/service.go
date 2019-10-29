// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3

import (
	"context"
	"io"
	"io/ioutil"
	"math"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	"github.com/opensds/multi-cloud/backend/proto"
	backendpb "github.com/opensds/multi-cloud/backend/proto"
	"github.com/opensds/multi-cloud/s3/proto"
	log "github.com/sirupsen/logrus"
)

const (
	s3Service      = "s3"
	backendService = "backend"
)

type APIService struct {
	s3Client      s3.S3Service
	backendClient backend.BackendService
}

func NewAPIService(c client.Client) *APIService {
	return &APIService{
		s3Client:      s3.NewS3Service(s3Service, c),
		backendClient: backend.NewBackendService(backendService, c),
	}
}

func IsQuery(request *restful.Request, name string) bool {
	params := request.Request.URL.Query()
	if params == nil {
		return false
	}
	if _, ok := params[name]; ok {
		return true
	}
	return false
}
func HasHeader(request *restful.Request, name string) bool {
	param := request.HeaderParameter(name)
	if param == "" {
		return false
	}
	return true
}

func ReadBody(r *restful.Request) []byte {
	var reader io.Reader = r.Request.Body
	b, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil
	}
	return b
}

/*func getBackendClient(ctx context.Context, s *APIService, bucketName string) datastore.DataStoreAdapter {
	log.Infof("bucketName is %v:\n", bucketName)
	bucket, err := s.s3Client.GetBucket(ctx, &s3.Bucket{Name: bucketName})
	if err != nil {
		return nil
	}

	log.Infof("bucketName is %v\n", bucketName)
	backendRep, backendErr := s.backendClient.ListBackend(ctx, &backendpb.ListBackendRequest{
		Offset: 0,
		Limit:  math.MaxInt32,
		Filter: map[string]string{"name": bucket.Backend}})
	log.Infof("backendErr is %v:", backendErr)
	if backendErr != nil {
		log.Errorf("get backend %s failed.", bucket.Backend)
		return nil
	}
	log.Infof("backendRep is %v:", backendRep)
	backend := backendRep.Backends[0]
	client, _ := datastore.Init(backend)
	return client
}
*/

func (s *APIService) getBucketMeta(ctx context.Context, bucketName string) *s3.Bucket {
	bucket, err := s.s3Client.GetBucket(ctx, &s3.BaseRequest{Id: bucketName})
	if err != nil {
		log.Infof("get bucket[name=%s] failed, err=%v.\n", bucketName, err)
		return nil
	}

	return bucket
}

func (s *APIService) isBackendExist(ctx context.Context, backendName string) bool {
	flag := false

	backendRep, backendErr := s.backendClient.ListBackend(ctx, &backendpb.ListBackendRequest{
		Offset: 0,
		Limit:  math.MaxInt32,
		Filter: map[string]string{"name": backendName}})
	log.Infof("backendErr is %v:", backendErr)
	if backendErr != nil {
		log.Infof("Get backend[name=%s] failed.", backendName)
	} else {
		log.Infof("backendRep=%+v\n", backendRep)
		if len(backendRep.Backends) > 0 {
			log.Infof("backend[name=%s] exist.", backendName)
			flag = true
		}
	}

	return flag
}

