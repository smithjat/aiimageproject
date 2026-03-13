const app = getApp<IAppOption>()

Page({
  data: {
    username: '',
    password: '',
    loading: false
  },

  onInputChange(e: WechatMiniprogram.Input) {
    const { field } = e.currentTarget.dataset
    this.setData({
      [field]: e.detail.value
    })
  },

  async handleLogin() {
    const { username, password } = this.data
    
    if (!username || !password) {
      wx.showToast({ title: '请输入用户名和密码', icon: 'none' })
      return
    }

    this.setData({ loading: true })

    try {
      const res = await wx.request({
        url: `${app.globalData.baseUrl}/auth/login`,
        method: 'POST',
        data: { username, password }
      }) as WechatMiniprogram.RequestSuccessCallbackResult<{ code: number; data: { token: string; user: any } }>

      if (res.data.code === 0) {
        app.globalData.token = res.data.data.token
        app.globalData.userInfo = res.data.data.user
        wx.setStorageSync('token', res.data.data.token)
        wx.showToast({ title: '登录成功', icon: 'success' })
        setTimeout(() => {
          wx.switchTab({ url: '/pages/index/index' })
        }, 1500)
      } else {
        wx.showToast({ title: '登录失败', icon: 'none' })
      }
    } catch (error) {
      wx.showToast({ title: '网络错误', icon: 'none' })
    } finally {
      this.setData({ loading: false })
    }
  },

  handleWechatLogin() {
    app.login().then(() => {
      wx.showToast({ title: '登录成功', icon: 'success' })
      setTimeout(() => {
        wx.switchTab({ url: '/pages/index/index' })
      }, 1500)
    }).catch((err) => {
      wx.showToast({ title: err || '登录失败', icon: 'none' })
    })
  }
})
