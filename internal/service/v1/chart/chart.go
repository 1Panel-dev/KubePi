package chart

import (
	v1Chart "github.com/KubeOperator/kubepi/internal/model/v1/chart"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	costomStorm "github.com/KubeOperator/kubepi/pkg/storm"
	"github.com/KubeOperator/kubepi/pkg/util/helm"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"helm.sh/helm/v3/pkg/repo"
	"time"
)

type Service interface {
	common.DBService
	Create(chart *v1Chart.Chart, options common.DBOptions) error
	Search(num, size int, cluster, pattern string, options common.DBOptions) ([]v1Chart.Chart, int, error)
	Delete(name string, options common.DBOptions) error
	Update(name string, char *v1Chart.Chart, options common.DBOptions) error
	GetByName(name string, options common.DBOptions) (*v1Chart.Chart, error)
	SearchRepo(cluster string) ([]*repo.Entry, error)
	AddRepo(cluster string, create *v1Chart.RepoCreate) error
}

func NewService() Service {
	return &service{
		clusterService: cluster.NewService(),
	}
}

type service struct {
	common.DefaultDBService
	clusterService cluster.Service
}

func (c *service) SearchRepo(cluster string) ([]*repo.Entry, error) {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return nil, err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return nil, err
	}
	repos, err := helmClient.ListRepo()
	if err != nil {
		return nil, err
	}
	return repos, err
}

func (c *service) AddRepo(cluster string, create *v1Chart.RepoCreate) error {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return err
	}
	err = helmClient.AddRepo(create.Name, create.Url, create.UserName, create.Password)
	if err != nil {
		return err
	}
	return nil
}

func (c *service) Search(num, size int, cluster, pattern string, options common.DBOptions) ([]v1Chart.Chart, int, error) {

	db := c.GetDB(options)
	query := func() storm.Query {
		if pattern != "" {
			return db.Select(q.And(q.Eq("Cluster", cluster), q.Or(
				costomStorm.Like("Name", pattern),
			))).OrderBy("CreateAt")
		}
		return db.Select().OrderBy("CreateAt")
	}()
	if num != 0 && size != 0 {
		query.Limit(size).Skip((num - 1) * size)
	}
	count, err := query.Count(&v1Chart.Chart{})
	if err != nil {
		return nil, 0, err
	}
	charts := make([]v1Chart.Chart, 0)
	if err := query.Find(&charts); err != nil {
		return nil, 0, err
	}

	return charts, count, nil
}

func (c *service) Create(chart *v1Chart.Chart, options common.DBOptions) error {

	cluster, err := c.clusterService.Get(chart.Cluster, common.DBOptions{})
	if err != nil {
		return err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        cluster.Spec.Connect.Forward.ApiServer,
		BearerToken: cluster.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return err
	}
	_, err = helmClient.List()
	if err != nil {
		return err
	}
	db := c.GetDB(options)
	chart.UUID = uuid.New().String()
	chart.CreateAt = time.Now()
	chart.UpdateAt = time.Now()
	return db.Save(chart)
}

func (c *service) Update(name string, char *v1Chart.Chart, options common.DBOptions) error {
	var chart v1Chart.Chart
	db := c.GetDB(options)
	if err := db.One("Name", name, &chart); err != nil {
		return err
	}
	char.UUID = chart.UUID
	char.CreateAt = chart.CreateAt
	char.UpdateAt = time.Now()
	return db.Update(char)
}

func (c *service) GetByName(name string, options common.DBOptions) (*v1Chart.Chart, error) {
	db := c.GetDB(options)
	var chart v1Chart.Chart
	query := db.Select(q.Eq("Name", name))

	if err := query.First(&chart); err != nil {
		return nil, err
	}
	return &chart, nil
}

func (c *service) Delete(name string, options common.DBOptions) error {
	db := c.GetDB(options)
	var chart v1Chart.Chart
	if err := db.One("Name", name, &chart); err != nil {
		return err
	}
	return db.DeleteStruct(&chart)
}
