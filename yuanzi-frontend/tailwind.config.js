export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      // Stitch 设计系统令牌 (完全匹配 stitch_yuanzi_baby_app)
      colors: {
        // 主色：珊瑚粉
        primary: '#ff998a',
        'primary-light': '#ffb4a2',
        'primary-dark': '#e8887a',
        'primary-soft': '#ffeeeb',

        // 背景色（暖白色）
        'background-light': '#fffbf7',
        'background-dark': '#23110f',
        'background-card': '#ffffff',
        'neutral-soft': '#f4e8e6',

        // 辅助色：薄荷绿
        mint: '#a2d5c6',
        'mint-light': '#e7f3eb',
        'mint-text': '#4e9767',

        // 功能色（图标背景）
        'blue-100': '#dbeafe',
        'orange-100': '#fed7aa',
        'indigo-100': '#e0e7ff',

        // 文字色（使用 slate）
        'text-primary': '#0f172a',
        'text-secondary': '#475569',
        'text-tertiary': '#94a3b8',

        // 边框色
        'border-light': '#f0e6de',

        // 语义色
        success: '#6bcb77',
        warning: '#ffd54f',
        error: '#ff6b6b',
        info: '#4fc3f7',
      },

      // 字体（Plus Jakarta Sans + 中文）
      fontFamily: {
        display: [
          'Plus Jakarta Sans',
          'Noto Sans SC',
          'system-ui',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'sans-serif',
        ],
      },

      // 圆角（匹配 Stitch：0.5rem, 1rem, 1.5rem）
      borderRadius: {
        'none': '0',
        'sm': '0.375rem',   // 6px
        'DEFAULT': '0.5rem', // 8px
        'md': '0.75rem',    // 12px
        'lg': '1rem',       // 16px
        'xl': '1.5rem',     // 24px
        '2xl': '2rem',      // 32px
        'full': '9999px',
      },

      // 阴影（匹配 Stitch）
      boxShadow: {
        'sm': '0 1px 2px rgba(0,0,0,0.04)',
        'card': '0 2px 8px rgba(0,0,0,0.06)',
        'md': '0 4px 12px rgba(0,0,0,0.08)',
        'lg': '0 8px 24px rgba(0,0,0,0.1)',
        'xl': '0 12px 32px rgba(0,0,0,0.12)',
        '2xl': '0 16px 48px rgba(0,0,0,0.15)',
        'primary': '0 4px 14px rgba(255, 153, 138, 0.3)',
      },

      // 间距
      spacing: {
        'xs': '4px',
        'sm': '8px',
        'md': '12px',
        'lg': '16px',
        'xl': '20px',
        '2xl': '24px',
        '3xl': '32px',
      },

      // 字体大小
      fontSize: {
        'xs': ['12px', { lineHeight: '1.5' }],
        'sm': ['14px', { lineHeight: '1.6' }],
        'base': ['16px', { lineHeight: '1.5' }],
        'lg': ['18px', { lineHeight: '1.4' }],
        'xl': ['20px', { lineHeight: '1.3' }],
        '2xl': ['24px', { lineHeight: '1.3' }],
        '3xl': ['28px', { lineHeight: '1.2' }],
      },
    },
  },
  plugins: [],
}
