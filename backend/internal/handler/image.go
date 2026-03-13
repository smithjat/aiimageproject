package handler

import (
	"aiimageproject/pkg/imageengine"
	"aiimageproject/pkg/response"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	engineClient     *imageengine.Client
	engineClientOnce sync.Once
)

func getEngineClient() *imageengine.Client {
	engineClientOnce.Do(func() {
		engineClient = imageengine.NewClient()
	})
	return engineClient
}

type Text2ImageRequest struct {
	Prompt         string  `json:"prompt" binding:"required"`
	ModelID        string  `json:"model_id"`
	NegativePrompt string  `json:"negative_prompt"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	CFGScale       float64 `json:"cfg_scale"`
	Steps          int     `json:"steps"`
	Seed           int64   `json:"seed"`
	Style          string  `json:"style"`
	Quality        string  `json:"quality"`
}

type ImageEditRequest struct {
	Image          string  `json:"image" binding:"required"`
	Prompt         string  `json:"prompt" binding:"required"`
	ModelID        string  `json:"model_id"`
	NegativePrompt string  `json:"negative_prompt"`
	Mask           string  `json:"mask"`
	Strength       float64 `json:"strength"`
	CFGScale       float64 `json:"cfg_scale"`
	Steps          int     `json:"steps"`
	Seed           int64   `json:"seed"`
}

func Text2Image(c *gin.Context) {
	var req Text2ImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	engineReq := &imageengine.Text2ImageRequest{
		Prompt:         req.Prompt,
		ModelID:        req.ModelID,
		NegativePrompt: req.NegativePrompt,
		Width:          req.Width,
		Height:         req.Height,
		CFGScale:       req.CFGScale,
		Steps:          req.Steps,
		Seed:           req.Seed,
		Style:          req.Style,
		Quality:        req.Quality,
	}

	result, err := getEngineClient().Text2Image(engineReq)
	if err != nil {
		response.Fail(c, 500, "生图失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

func ImageEdit(c *gin.Context) {
	var req ImageEditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	engineReq := &imageengine.ImageEditRequest{
		Image:          req.Image,
		Prompt:         req.Prompt,
		ModelID:        req.ModelID,
		NegativePrompt: req.NegativePrompt,
		Mask:           req.Mask,
		Strength:       req.Strength,
		CFGScale:       req.CFGScale,
		Steps:          req.Steps,
		Seed:           req.Seed,
	}

	result, err := getEngineClient().ImageEdit(engineReq)
	if err != nil {
		response.Fail(c, 500, "图片编辑失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

func GetModels(c *gin.Context) {
	models, err := getEngineClient().GetModels()
	if err != nil {
		response.Fail(c, 500, "获取模型列表失败: "+err.Error())
		return
	}

	response.Success(c, models)
}

func GetModelsByCapability(c *gin.Context) {
	capability := c.Param("capability")
	if capability != "generation" && capability != "editing" {
		response.Fail(c, 400, "无效的能力类型")
		return
	}

	models, err := getEngineClient().GetModelsByCapability(capability)
	if err != nil {
		response.Fail(c, 500, "获取模型列表失败: "+err.Error())
		return
	}

	response.Success(c, models)
}

func GetCapabilities(c *gin.Context) {
	response.Success(c, []string{"generation", "editing"})
}

func HealthCheck(c *gin.Context) {
	ok, err := getEngineClient().HealthCheck()
	if err != nil {
		response.Fail(c, 503, "图片引擎不可用: "+err.Error())
		return
	}

	if !ok {
		response.Fail(c, 503, "图片引擎状态异常")
		return
	}

	response.Success(c, gin.H{
		"status":        "ok",
		"engine_status": "connected",
	})
}
