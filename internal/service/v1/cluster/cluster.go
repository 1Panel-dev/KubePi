package cluster

import (
	v1Cluster "github.com/KubeOperator/kubepi/internal/model/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterapp"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/storm"
	"github.com/KubeOperator/kubepi/pkg/util/lang"
	storm2 "github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	Create(cluster *v1Cluster.Cluster, options common.DBOptions) error
	Update(name string, cluster *v1Cluster.Cluster, options common.DBOptions) error
	Get(name string, options common.DBOptions) (*v1Cluster.Cluster, error)
	List(options common.DBOptions) ([]v1Cluster.Cluster, error)
	Delete(name string, options common.DBOptions) error
	Search(num, size int, conditions common.Conditions, options common.DBOptions) ([]v1Cluster.Cluster, int, error)
}

func NewService() Service {
	return &cluster{
		DefaultDBService:  common.DefaultDBService{},
		clusterAppService: clusterapp.NewService(),
	}
}

type cluster struct {
	common.DefaultDBService
	clusterAppService clusterapp.Service
}

func (c *cluster) Update(name string, cluster *v1Cluster.Cluster, options common.DBOptions) error {
	db := c.GetDB(options)
	r, err := c.Get(name, options)
	if err != nil {
		return err
	}
	cluster.UUID = r.UUID
	cluster.CreateAt = r.CreateAt
	cluster.UpdateAt = time.Now()
	return db.Update(cluster)
}

func (c *cluster) Create(cluster *v1Cluster.Cluster, options common.DBOptions) error {
	db := c.GetDB(options)
	cluster.UUID = uuid.New().String()
	cluster.CreateAt = time.Now()
	cluster.UpdateAt = time.Now()
	return db.Save(cluster)
}

func (c *cluster) Get(name string, options common.DBOptions) (*v1Cluster.Cluster, error) {
	db := c.GetDB(options)
	var cluster v1Cluster.Cluster
	if err := db.One("Name", name, &cluster); err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *cluster) List(options common.DBOptions) ([]v1Cluster.Cluster, error) {
	db := c.GetDB(options)
	var clusters []v1Cluster.Cluster
	if err := db.All(&clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

func (c *cluster) Search(num, size int, conditions common.Conditions, options common.DBOptions) ([]v1Cluster.Cluster, int, error) {
	db := c.GetDB(options)

	query := db.Select()

	query = func() storm2.Query {
		var ms []q.Matcher

		for k := range conditions {
			if k == "quick" {
				ms = append(ms, storm.Like("Name", conditions[k].Value))
			} else if k == "labels" {
				switch conditions[k].Operator {
				case "like":
					ms = append(ms, storm.ArrayValueLike("Labels", conditions[k].Value))
				case "not like":
					ms = append(ms, q.Not(storm.ArrayValueLike("Labels", conditions[k].Value)))
				case "eq":
					ms = append(ms, storm.ArrayValueEq("Labels", conditions[k].Value))
				case "ne":
					ms = append(ms, q.Not(storm.ArrayValueEq("Labels", conditions[k].Value)))
				}
			} else {
				field := lang.FirstToUpper(conditions[k].Field)
				switch conditions[k].Operator {
				case "eq":
					ms = append(ms, q.Eq(field, conditions[k].Value))
				case "ne":
					ms = append(ms, q.Not(q.Eq(field, conditions[k].Value)))
				case "like":
					ms = append(ms, storm.Like(field, conditions[k].Value))
				case "not like":
					ms = append(ms, q.Not(storm.Like(field, conditions[k].Value)))
				}
			}
		}
		if len(conditions) > 0 {
			return db.Select(q.And(ms...)).Limit(size).Skip((num - 1) * size).OrderBy("CreateAt").Reverse()
		} else {
			return db.Select().Limit(size).Skip((num - 1) * size).OrderBy("CreateAt").Reverse()
		}
	}()
	count, err := db.Select().Count(&v1Cluster.Cluster{})
	if err != nil {
		return nil, 0, err
	}
	clusters := make([]v1Cluster.Cluster, 0)
	if err := query.Find(&clusters); err != nil {
		return clusters, 0, err
	}
	return clusters, count, nil
}

func (c *cluster) Delete(name string, options common.DBOptions) error {
	db := c.GetDB(options)
	cluster, err := c.Get(name, options)
	if err != nil {
		return err
	}
	if err := c.clusterAppService.DeleteByCluster(name, options); err != nil && err != storm2.ErrNotFound {
		return err
	}
	return db.DeleteStruct(cluster)
}
