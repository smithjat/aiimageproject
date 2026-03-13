App({
  globalData: {
    userInfo: null,
    token: '',
    baseUrl: 'http://localhost:8080/api/v1'
  },

  onLaunch() {
    const token = wx.getStorageSync('token')
    if (token) {
      this.globalData.token = token
      this.checkLogin()
    }
  },

  checkLogin() {
    wx.request({
      url: `${this.globalData.baseUrl}/users/profile`,
      header: {
        'Authorization': `Bearer ${this.globalData.token}`
      },
      success: (res) => {
        if (res.data.code === 0) {
          this.globalData.userInfo = res.data.data
        } else {
          this.logout()
        }
      },
      fail: () => {
        this.logout()
      }
    })
  },

  logout() {
    this.globalData.token = ''
    this.globalData.userInfo = null
    wx.removeStorageSync('token')
  },

  login() {
    return new Promise((resolve, reject) => {
      wx.login({
        success: (res) => {
          if (res.code) {
            wx.request({
              url: `${this.globalData.baseUrl}/auth/wechat-login`,
              method: 'POST',
              data: { code: res.code },
              success: (res) => {
                if (res.data.code === 0) {
                  this.globalData.token = res.data.data.token
                  this.globalData.userInfo = res.data.data.user
                  wx.setStorageSync('token', res.data.data.token)
                  resolve(res.data.data)
                } else {
                  reject(res.data.message)
                }
              },
              fail: reject
            })
          } else {
            reject(res.errMsg)
          }
        },
        fail: reject
      })
    })
  }
})
