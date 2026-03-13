<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import InspirationLibrary from '@/components/InspirationLibrary.vue'
import { modelApi, imageApi, type AIModel } from '@/api'

const route = useRoute()
const models = ref<AIModel[]>([])
const prompt = ref('')
const negativePrompt = ref('')
const selectedModel = ref<string>('')
const isGenerating = ref(false)
const generatedImages = ref<string[]>([])
const showAdvanced = ref(false)
const uploadedImage = ref<string | null>(null)
const maskImage = ref<string | null>(null)
const strength = ref(0.75)

const loadModels = async () => {
  try {
    const res = await modelApi.getByCapability('editing')
    models.value = res.data
  } catch (error) {
    console.error('加载模型失败', error)
  }
}

const handleInspirationSelect = (inspiration: { prompt?: string; imageUrl?: string; type?: string }) => {
  if (inspiration.type === 'image' && inspiration.imageUrl) {
    uploadedImage.value = inspiration.imageUrl
    ElMessage.success('已应用灵感图片')
  }
  if (inspiration.prompt) {
    prompt.value = inspiration.prompt
  }
}

const selectModel = (modelId: string) => {
  selectedModel.value = modelId
}

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      uploadedImage.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

const handleMaskUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      maskImage.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

const generateImage = async () => {
  if (!uploadedImage.value) {
    ElMessage.warning('请先上传图片')
    return
  }
  if (!prompt.value.trim()) {
    ElMessage.warning('请输入编辑提示词')
    return
  }
  if (!selectedModel.value) {
    ElMessage.warning('请选择模型')
    return
  }

  isGenerating.value = true
  try {
    const res = await imageApi.imageEdit({
      image: uploadedImage.value,
      prompt: prompt.value,
      model_id: selectedModel.value,
      negative_prompt: negativePrompt.value,
      mask: maskImage.value || undefined,
      strength: strength.value,
    })
    generatedImages.value = res.data.images
    ElMessage.success('编辑成功')
  } catch (error) {
    ElMessage.error('编辑失败，请重试')
  } finally {
    isGenerating.value = false
  }
}

const downloadImage = async (imageUrl: string, index: number) => {
  try {
    const response = await fetch(imageUrl)
    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `ai-edited-${Date.now()}-${index + 1}.png`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('下载成功')
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

onMounted(() => {
  loadModels()
  if (route.query.image) {
    uploadedImage.value = route.query.image as string
    ElMessage.success('已加载图片，请选择模型进行编辑')
  }
})
</script>

<template>
  <div class="image-edit">
    <div class="main-content">
      <div class="page-header">
        <div class="header-left">
          <div class="breadcrumb">
            <span class="breadcrumb-item">AI创作</span>
            <span class="breadcrumb-sep">/</span>
            <span class="breadcrumb-item active">图片编辑</span>
          </div>
          <h1>图片编辑</h1>
        </div>
        <div class="header-right">
          <div class="credits-info">
            <span class="credits-label">账户余额</span>
            <span class="credits-value">
              <span class="credits-icon">💎</span>
              <span class="credits-count">128 积分</span>
            </span>
          </div>
          <button class="btn-recharge">立即充值</button>
        </div>
      </div>

      <div class="content-body">
        <div class="left-panel">
          <div class="panel-card">
            <div class="card-header">
              <span class="card-title">上传图片</span>
              <span class="card-extra">支持 JPG、PNG 格式</span>
            </div>
            <div class="upload-area-wrapper">
              <div v-if="!uploadedImage" class="upload-area">
                <input type="file" accept="image/*" @change="handleFileUpload" id="file-input" hidden />
                <label for="file-input" class="upload-label">
                  <div class="upload-icon">📤</div>
                  <p class="upload-text">点击上传图片</p>
                  <p class="upload-hint">或将图片拖拽到此处</p>
                </label>
              </div>
              <div v-else class="preview-area">
                <img :src="uploadedImage" alt="预览" />
                <div class="preview-actions">
                  <button class="btn-preview-action" @click="uploadedImage = null">重新上传</button>
                </div>
              </div>
            </div>
          </div>

          <div class="panel-card">
            <div class="card-header">
              <span class="card-title">编辑提示词</span>
              <span class="card-extra">描述您想要的编辑效果</span>
            </div>
            <div class="prompt-area">
              <textarea
                v-model="prompt"
                class="prompt-input"
                placeholder="请输入编辑效果描述，例如：将背景替换为现代办公室场景，保持人物不变..."
                rows="4"
              ></textarea>
              <div class="prompt-footer">
                <div class="prompt-tools">
                  <button class="tool-btn">🌐 翻译</button>
                  <button class="tool-btn">✨ 优化</button>
                </div>
                <span class="char-count">{{ prompt.length }}/500</span>
              </div>
            </div>

            <div class="advanced-toggle" @click="showAdvanced = !showAdvanced">
              <span :class="['toggle-icon', { expanded: showAdvanced }]">▶</span>
              <span>高级设置</span>
            </div>

            <div v-if="showAdvanced" class="advanced-settings">
              <div class="setting-item">
                <label>负面提示词</label>
                <textarea
                  v-model="negativePrompt"
                  placeholder="输入不想出现的元素..."
                  rows="2"
                ></textarea>
              </div>
              <div class="setting-item">
                <label>遮罩图片（可选）</label>
                <div class="mask-upload">
                  <input type="file" accept="image/*" @change="handleMaskUpload" id="mask-input" hidden />
                  <label for="mask-input" class="mask-label">
                    <span v-if="!maskImage">点击上传遮罩</span>
                    <span v-else>已上传遮罩 ✓</span>
                  </label>
                  <button v-if="maskImage" class="btn-clear-mask" @click="maskImage = null">清除</button>
                </div>
              </div>
              <div class="setting-item">
                <label>编辑强度: {{ strength }}</label>
                <input type="range" v-model="strength" min="0.1" max="1" step="0.05" class="strength-slider" />
              </div>
            </div>
          </div>

          <div class="panel-card">
            <div class="card-header">
              <span class="card-title">选择模型</span>
              <span class="card-link">查看全部模型 →</span>
            </div>
            <div class="model-grid">
              <div
                v-for="model in models"
                :key="model.model_id"
                :class="['model-card', { active: selectedModel === model.model_id }]"
                @click="selectModel(model.model_id)"
              >
                <div class="model-info">
                  <div class="model-name">{{ model.model_name }}</div>
                  <div class="model-desc">{{ model.description }}</div>
                </div>
                <div v-if="selectedModel === model.model_id" class="model-check">
                  <span>✓</span>
                </div>
              </div>
            </div>
          </div>

          <div class="action-bar">
            <div class="cost-info">
              <span class="cost-label">预计消耗</span>
              <span class="cost-value">3 积分</span>
            </div>
            <button
              class="btn-generate"
              :disabled="!uploadedImage || !prompt.trim() || !selectedModel || isGenerating"
              @click="generateImage"
            >
              <span v-if="isGenerating" class="loading-content">
                <span class="loading-spinner"></span>
                <span>处理中...</span>
              </span>
              <span v-else>开始编辑</span>
            </button>
          </div>
        </div>

        <div class="right-panel">
          <div class="panel-card result-panel">
            <div class="card-header">
              <span class="card-title">编辑结果</span>
              <span v-if="generatedImages.length" class="card-extra">{{ generatedImages.length }} 张图片</span>
            </div>
            
            <div v-if="generatedImages.length === 0" class="empty-result">
              <div class="empty-icon">🎨</div>
              <p class="empty-title">等待编辑</p>
              <p class="empty-hint">上传图片并输入提示词后开始编辑</p>
            </div>
            
            <div v-else class="result-grid">
              <div v-for="(img, index) in generatedImages" :key="index" class="result-item">
                <img :src="img" :alt="`编辑图片 ${index + 1}`" />
                <div class="result-overlay">
                  <div class="result-actions">
                    <button class="action-btn" @click="downloadImage(img, index)">下载</button>
                    <button class="action-btn">收藏</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <InspirationLibrary type="edit" @select="handleInspirationSelect" />
  </div>
</template>

<style scoped>
.image-edit {
  display: flex;
  min-height: 100vh;
  margin-left: 240px;
  background: #f5f7fa;
}

.main-content {
  flex: 1;
  padding: 24px;
  margin-right: 320px;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.breadcrumb-item {
  font-size: 13px;
  color: #64748b;
}

.breadcrumb-item.active {
  color: #1e293b;
  font-weight: 500;
}

.breadcrumb-sep {
  color: #cbd5e1;
}

.header-left h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1e293b;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.credits-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.credits-label {
  font-size: 12px;
  color: #64748b;
}

.credits-value {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
}

.credits-icon {
  font-size: 16px;
}

.btn-recharge {
  padding: 8px 20px;
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-recharge:hover {
  background: #1d4ed8;
}

.content-body {
  display: flex;
  gap: 20px;
  flex: 1;
}

.left-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.right-panel {
  width: 360px;
}

.panel-card {
  background: #fff;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e2e8f0;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
}

.card-extra {
  font-size: 13px;
  color: #64748b;
}

.card-link {
  font-size: 13px;
  color: #2563eb;
  cursor: pointer;
}

.card-link:hover {
  text-decoration: underline;
}

.upload-area-wrapper {
  padding: 16px 20px;
}

.upload-area {
  border: 2px dashed #e2e8f0;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.upload-area:hover {
  border-color: #2563eb;
  background: #f8fafc;
}

.upload-label {
  cursor: pointer;
  display: block;
}

.upload-icon {
  font-size: 40px;
  margin-bottom: 12px;
}

.upload-text {
  font-size: 15px;
  color: #1e293b;
  margin-bottom: 4px;
}

.upload-hint {
  font-size: 13px;
  color: #94a3b8;
}

.preview-area {
  text-align: center;
}

.preview-area img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 6px;
  margin-bottom: 12px;
}

.preview-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.btn-preview-action {
  padding: 8px 20px;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-preview-action:hover {
  background: #e2e8f0;
}

.prompt-area {
  padding: 16px 20px;
}

.prompt-input {
  width: 100%;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 12px;
  font-size: 14px;
  line-height: 1.6;
  resize: none;
  outline: none;
  transition: border-color 0.2s;
}

.prompt-input:focus {
  border-color: #2563eb;
}

.prompt-input::placeholder {
  color: #94a3b8;
}

.prompt-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.prompt-tools {
  display: flex;
  gap: 8px;
}

.tool-btn {
  padding: 6px 12px;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 13px;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s;
}

.tool-btn:hover {
  background: #e2e8f0;
  border-color: #cbd5e1;
}

.char-count {
  font-size: 12px;
  color: #94a3b8;
}

.advanced-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-top: 1px solid #e2e8f0;
  color: #64748b;
  font-size: 14px;
  cursor: pointer;
  user-select: none;
}

.advanced-toggle:hover {
  color: #2563eb;
}

.toggle-icon {
  font-size: 10px;
  transition: transform 0.2s;
}

.toggle-icon.expanded {
  transform: rotate(90deg);
}

.advanced-settings {
  padding: 16px 20px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}

.setting-item {
  margin-bottom: 16px;
}

.setting-item:last-child {
  margin-bottom: 0;
}

.setting-item label {
  display: block;
  font-size: 13px;
  color: #475569;
  margin-bottom: 8px;
}

.setting-item textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
  resize: none;
}

.setting-item textarea:focus {
  border-color: #2563eb;
}

.mask-upload {
  display: flex;
  align-items: center;
  gap: 12px;
}

.mask-label {
  display: inline-block;
  padding: 10px 16px;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 14px;
  color: #475569;
  cursor: pointer;
  transition: all 0.2s;
}

.mask-label:hover {
  background: #e2e8f0;
}

.btn-clear-mask {
  padding: 10px 16px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 6px;
  font-size: 14px;
  color: #dc2626;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-clear-mask:hover {
  background: #fee2e2;
}

.strength-slider {
  width: 100%;
  height: 6px;
  -webkit-appearance: none;
  appearance: none;
  background: #e2e8f0;
  border-radius: 3px;
  outline: none;
}

.strength-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 18px;
  height: 18px;
  background: #2563eb;
  border-radius: 50%;
  cursor: pointer;
}

.model-grid {
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.model-card {
  position: relative;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.model-card:hover {
  border-color: #2563eb;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.1);
}

.model-card.active {
  border-color: #2563eb;
  background: #eff6ff;
}

.model-info {
  padding-right: 24px;
}

.model-name {
  font-size: 13px;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 4px;
}

.model-desc {
  font-size: 11px;
  color: #64748b;
  line-height: 1.4;
}

.model-check {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 20px;
  height: 20px;
  background: #2563eb;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
}

.action-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.cost-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.cost-label {
  font-size: 12px;
  color: #64748b;
}

.cost-value {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.btn-generate {
  padding: 12px 32px;
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-generate:hover:not(:disabled) {
  background: #1d4ed8;
}

.btn-generate:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}

.loading-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #fff;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.result-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.empty-result {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-title {
  font-size: 15px;
  font-weight: 500;
  color: #475569;
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 13px;
  color: #94a3b8;
}

.result-grid {
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.result-item {
  position: relative;
  border-radius: 6px;
  overflow: hidden;
  aspect-ratio: 1;
}

.result-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.result-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.6), transparent);
  opacity: 0;
  transition: opacity 0.2s;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 12px;
}

.result-item:hover .result-overlay {
  opacity: 1;
}

.result-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 4px;
  font-size: 12px;
  color: #1e293b;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #fff;
}
</style>
