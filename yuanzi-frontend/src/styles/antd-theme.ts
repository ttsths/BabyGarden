export const antdTheme = {
  token: {
    // 品牌色 - 珊瑚橙
    colorPrimary: '#FF9A8B',
    colorPrimaryLight: '#FFB4A2',
    colorPrimaryDark: '#E8887A',
    
    // 背景色
    colorBackground: '#FFFBF7',
    colorBackgroundLight: '#FFFFFF',
    colorBackgroundSoft: '#FFF5ED',
    
    // 文字色
    colorText: '#2D2D2D',
    colorTextSecondary: '#6B7280',
    colorTextPlaceholder: '#9CA3AF',
    
    // 边框
    colorBorder: '#F0E6DE',
    
    // 圆角
    borderRadius: 12,
    borderRadiusLg: 16,
    
    // 字体
    fontFamily: '"PingFang SC", "Source Han Sans CN", -apple-system, sans-serif',
    fontSizeBase: 15,
    fontSizeLg: 18,
  },
  
  // 组件级定制
  components: {
    Button: {
      algorithm: true,
      token: {
        height: 40,
        heightLg: 48,
        heightSm: 32,
        paddingContentHorizontal: 16,
      },
    },
    Input: {
      token: {
        height: 40,
        fontSize: 15,
      },
    },
    Card: {
      token: {
        padding: 16,
        borderRadius: 12,
      },
    },
  },
};
