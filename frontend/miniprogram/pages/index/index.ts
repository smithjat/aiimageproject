const app = getApp<IAppOption>()

Page({
  data: {
    hasLogin: false
  },

  onLoad() {
    this.checkLogin()
  },

  onShow() {
    this.checkLogin()
  },

  checkLogin() {
    this.setData({
      hasLogin: !!app.globalData.token
    })
  },

  goLogin() {
    wx.navigateTo({
      url: '/pages/login/login'
    })
  }
})
