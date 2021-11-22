package clusterrepo

import (
	V1ClusterRepo "github.com/KubeOperator/kubepi/internal/model/v1/clusterrepo"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	List(cluster string, options common.DBOptions) (result []V1ClusterRepo.ClusterRepo, err error)
	Create(clusterRepo *V1ClusterRepo.ClusterRepo, options common.DBOptions) error
	Delete(cluster, repo string, options common.DBOptions) error
}

func NewService() Service {
	return &service{}
}

type service struct {
	common.DefaultDBService
}

func (s *service) List(cluster string, options common.DBOptions) (result []V1ClusterRepo.ClusterRepo, err error) {
	db := s.GetDB(options)
	query := db.Select(q.Eq("Cluster", cluster))
	if err = query.Find(&result); err != nil {
		return
	}
	return
}

func (s *service) Create(clusterRepo *V1ClusterRepo.ClusterRepo, options common.DBOptions) error {
	db := s.GetDB(options)
	clusterRepo.UUID = uuid.New().String()
	clusterRepo.CreateAt = time.Now()
	clusterRepo.UpdateAt = time.Now()
	return db.Save(clusterRepo)
}

func (s *service) Delete(cluster, repo string, options common.DBOptions) error {
	db := s.GetDB(options)
	query := db.Select(q.And(q.Eq("Cluster", cluster), q.Eq("Repo", repo)))
	var clusterRepo V1ClusterRepo.ClusterRepo
	if err := query.First(&clusterRepo); err != nil {
		return err
	}
	return db.DeleteStruct(&clusterRepo)
}
