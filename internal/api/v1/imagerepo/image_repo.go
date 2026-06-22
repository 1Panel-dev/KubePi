package imagerepo

import (
	goContext "context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/1Panel-dev/KubePi/internal/api/v1/commons"
	v1Session "github.com/1Panel-dev/KubePi/internal/api/v1/session"
	v1ImageRepo "github.com/1Panel-dev/KubePi/internal/model/v1/imagerepo"
	v1Role "github.com/1Panel-dev/KubePi/internal/model/v1/role"
	"github.com/1Panel-dev/KubePi/internal/server"
	"github.com/1Panel-dev/KubePi/internal/service/v1/clusterrepo"
	"github.com/1Panel-dev/KubePi/internal/service/v1/common"
	"github.com/1Panel-dev/KubePi/internal/service/v1/imagerepo"
	pkgV1 "github.com/1Panel-dev/KubePi/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	imageRepoService   imagerepo.Service
	clusterRepoService clusterrepo.Service
}

const imageRepoResource = "imagerepos"

func NewHandler() *Handler {
	return &Handler{
		imageRepoService:   imagerepo.NewService(),
		clusterRepoService: clusterrepo.NewService(),
	}
}

func sanitizeImageRepo(repo v1ImageRepo.ImageRepo) v1ImageRepo.ImageRepo {
	repo.Credential.Password = ""
	return repo
}

func sanitizeImageRepoList(repos []v1ImageRepo.ImageRepo) []v1ImageRepo.ImageRepo {
	for i := range repos {
		repos[i] = sanitizeImageRepo(repos[i])
	}
	return repos
}

func validateImageRepoEndpoint(endpoint string) error {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return fmt.Errorf("image repo endpoint is empty")
	}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid image repo endpoint: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("image repo endpoint only supports http or https")
	}
	if parsed.User != nil {
		return fmt.Errorf("image repo endpoint must not contain user info")
	}
	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("image repo endpoint host is empty")
	}
	if ip := net.ParseIP(host); ip != nil {
		if isUnsafeImageRepoIP(ip) {
			return fmt.Errorf("image repo endpoint host is not allowed")
		}
		return nil
	}

	ctx, cancel := goContext.WithTimeout(goContext.Background(), 2*time.Second)
	defer cancel()
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return fmt.Errorf("resolve image repo endpoint failed: %w", err)
	}
	for i := range addrs {
		if isUnsafeImageRepoIP(addrs[i].IP) {
			return fmt.Errorf("image repo endpoint resolved to an unsafe address")
		}
	}
	return nil
}

func isUnsafeImageRepoIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	if ipv4 := ip.To4(); ipv4 != nil {
		return ipv4[0] == 169 && ipv4[1] == 254 && ipv4[2] == 169 && ipv4[3] == 254
	}
	return false
}

func canManageImageRepos(ctx *context.Context) bool {
	profile, ok := ctx.Values().Get("profile").(v1Session.UserProfile)
	if !ok {
		return false
	}
	if profile.IsAdministrator {
		return true
	}
	roles, ok := ctx.Values().Get("roles").([]v1Role.Role)
	if !ok {
		return false
	}
	for i := range roles {
		for j := range roles[i].Rules {
			rule := roles[i].Rules[j]
			if !containsRoleValue(rule.Resource, imageRepoResource) && !containsRoleValue(rule.Resource, "*") {
				continue
			}
			if containsRoleValue(rule.Verbs, "create") || containsRoleValue(rule.Verbs, "update") || containsRoleValue(rule.Verbs, "*") {
				return true
			}
		}
	}
	return false
}

func containsRoleValue(values []string, value string) bool {
	for i := range values {
		if values[i] == value {
			return true
		}
	}
	return false
}

func (h *Handler) SearchRepos() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		var conditions commons.SearchConditions
		if err := ctx.ReadJSON(&conditions); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		repos, total, err := h.imageRepoService.Search(pageNum, pageSize, conditions.Conditions, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		repos = sanitizeImageRepoList(repos)
		ctx.Values().Set("data", pkgV1.Page{Items: repos, Total: total})
	}
}

// Create Repo
// @Tags repos
// @Summary Create repo
// @Description Create repo
// @Accept  json
// @Produce  json
// @Param request body RepoConfig true "request"
// @Success 200 {object} RepoConfig
// @Security ApiKeyAuth
// @Router /imagerepos [post]
func (h *Handler) CreateRepo() iris.Handler {
	return func(ctx *context.Context) {
		var req RepoConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := validateImageRepoEndpoint(req.ImageRepo.EndPoint); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.imageRepoService.Create(&req.ImageRepo, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		req.ImageRepo = sanitizeImageRepo(req.ImageRepo)
		ctx.Values().Set("data", req)
	}
}

// List Internal Repos
// @Tags repos
// @Summary List Internal repos
// @Description List Internal repos
// @Accept  json
// @Produce  json
// @Param request body RepoConfig true "request"
// @Success 200 {object} []string
// @Security ApiKeyAuth
// @Router /imagerepos/repositories/search [post]
func (h *Handler) ListInternalRepos() iris.Handler {
	return func(ctx *context.Context) {
		if !canManageImageRepos(ctx) {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.Values().Set("message", "only image repo managers can search remote repositories")
			return
		}
		var req RepoConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := validateImageRepoEndpoint(req.ImageRepo.EndPoint); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}

		names, err := h.imageRepoService.ListInternalRepos(req.ImageRepo, req.Page, req.Limit, req.Search)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", names)
	}
}

// Update Repo
// @Tags repos
// @Summary Update repo by name
// @Description Update repo by name
// @Accept  json
// @Produce  json
// @Param request body RepoConfig true "request"
// @Param name path string true "镜像仓库名称"
// @Success 200 {object} v1ImageRepo.ImageRepo
// @Security ApiKeyAuth
// @Router /imagerepos/{name} [put]
func (h *Handler) UpdateRepo() iris.Handler {
	return func(ctx *context.Context) {
		var req RepoConfig
		imageRepoName := ctx.Params().GetString("name")
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := validateImageRepoEndpoint(req.ImageRepo.EndPoint); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.imageRepoService.UpdateRepo(imageRepoName, &req.ImageRepo, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

// Delete Repo
// @Tags repos
// @Summary Delete repo by name
// @Description Delete repo by name
// @Accept  json
// @Produce  json
// @Param name path string true "镜像仓库名称"
// @Success 200 {object} v1ImageRepo.ImageRepo
// @Security ApiKeyAuth
// @Router /imagerepos/{name} [delete]
func (h *Handler) DeleteRepo() iris.Handler {
	return func(ctx *context.Context) {
		imageRepoName := ctx.Params().GetString("name")
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		txOptions := common.DBOptions{DB: tx}
		if err := h.imageRepoService.Delete(imageRepoName, txOptions); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			_ = tx.Rollback()
			return
		}
		if err := h.clusterRepoService.DeleteByRepo(imageRepoName, txOptions); err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}
}

// Get Repo
// @Tags repos
// @Summary Get repo by name
// @Description Get repo by name
// @Accept  json
// @Produce  json
// @Param name path string true "镜像仓库名称"
// @Success 200 {object} v1ImageRepo.ImageRepo
// @Security ApiKeyAuth
// @Router /imagerepos/{name} [get]
func (h *Handler) GetRepo() iris.Handler {
	return func(ctx *context.Context) {
		imageRepoName := ctx.Params().GetString("name")
		imageRepo, err := h.imageRepoService.GetByName(imageRepoName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		imageRepo = sanitizeImageRepo(imageRepo)
		ctx.Values().Set("data", imageRepo)
	}
}

// List Repo for Cluster
// @Tags repos
// @Summary List repo for cluster
// @Description Get repo for cluster
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Success 200 {object} []v1ImageRepo.ImageRepo
// @Security ApiKeyAuth
// @Router /imagerepos/cluster/{cluster} [get]
func (h *Handler) ListRepoForCluster() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("cluster")
		var imageRepos []v1ImageRepo.ImageRepo
		imageRepos, err := h.imageRepoService.ListByCluster(cluster, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		imageRepos = sanitizeImageRepoList(imageRepos)
		ctx.Values().Set("data", imageRepos)
	}
}

// List Images for Cluster
// @Tags repos
// @Summary List images for cluster
// @Description Get images for cluster
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Param repo path string true "镜像仓库名称"
// @Success 200 {object} []string
// @Security ApiKeyAuth
// @Router /imagerepos/{cluster}/{repo} [put]
func (h *Handler) ListImages() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("cluster")
		name := ctx.Params().GetString("repo")
		imageRepos, err := h.imageRepoService.ListImages(name, cluster, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", imageRepos)
	}
}

// List Images By Repo
// @Tags repos
// @Summary List images by repo
// @Description Get images by repo
// @Accept  json
// @Produce  json
// @Param repo path string true "镜像仓库名称"
// @Success 200 {object} v1ImageRepo.RepoResponse
// @Security ApiKeyAuth
// @Router /imagerepos/{repo}/search [get]
func (h *Handler) ListImagesByRepo() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("repo")
		var req RepoConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		res, err := h.imageRepoService.ListImagesByRepo(name, req.Page, req.Limit, req.Search, req.ContinueToken, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", res)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/imagerepos")
	sp.Post("/search", handler.SearchRepos())
	sp.Post("/", handler.CreateRepo())
	sp.Delete("/:name", handler.DeleteRepo())
	sp.Post("/repositories/search", handler.ListInternalRepos())
	sp.Get("/:name", handler.GetRepo())
	sp.Put("/:name", handler.UpdateRepo())
	sp.Get("/cluster/:cluster", handler.ListRepoForCluster())
	sp.Get("/images/:cluster/:repo", handler.ListImages())
	sp.Post("/images/:repo/search", handler.ListImagesByRepo())
}
