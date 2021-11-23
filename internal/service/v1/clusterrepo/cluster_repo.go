package clusterrepo

import (
	"encoding/base64"
	V1ClusterRepo "github.com/KubeOperator/kubepi/internal/model/v1/clusterrepo"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/imagerepo"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
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
	return &service{
		clusterService:  cluster.NewService(),
		imgarepoService: imagerepo.NewService(),
	}
}

type service struct {
	common.DefaultDBService
	clusterService  cluster.Service
	imgarepoService imagerepo.Service
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

	clu, err := s.clusterService.Get(clusterRepo.Cluster, options)
	if err != nil {
		return err
	}
	client := kubernetes.NewKubernetes(clu)
	if err := client.Ping(); err != nil {
		return err
	}
	repo, err := s.imgarepoService.GetByName(clusterRepo.Repo, options)
	if err != nil {
		return err
	}
	password := base64.StdEncoding.EncodeToString([]byte(repo.Credential.Password))
	err = client.CreateImageRepo(kubernetes.CustomImageRepo{
		Type:     repo.Type,
		Name:     repo.Name,
		Username: repo.Credential.Username,
		Password: password,
		Url:      repo.EndPoint,
	})
	if err != nil {
		return err
	}
	db := s.GetDB(options)
	clusterRepo.UUID = uuid.New().String()
	clusterRepo.CreateAt = time.Now()
	clusterRepo.UpdateAt = time.Now()
	err = db.Save(clusterRepo)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(cluster, repo string, options common.DBOptions) error {
	clu, err := s.clusterService.Get(cluster, options)
	if err != nil {
		return err
	}
	client := kubernetes.NewKubernetes(clu)
	if err := client.Ping(); err != nil {
		return err
	}
	if err := client.DeleteImageRepo(repo); err != nil {
		return err
	}
	db := s.GetDB(options)
	query := db.Select(q.And(q.Eq("Cluster", cluster), q.Eq("Repo", repo)))
	var clusterRepo V1ClusterRepo.ClusterRepo
	if err := query.First(&clusterRepo); err != nil {
		return err
	}
	return db.DeleteStruct(&clusterRepo)
}
