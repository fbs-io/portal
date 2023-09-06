/*
 * @Author: reel
 * @Date: 2023-07-30 22:36:56
 * @LastEditors: reel
 * @LastEditTime: 2023-09-05 23:13:32
 * @Description: 请填写简介
 */
import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {

  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    base: "/website/",
    plugins: [
      vue(),
      AutoImport({
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
        'vue-i18n': 'vue-i18n/dist/vue-i18n.cjs.js',
      },
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue']
    },
    server: {
      host: '0.0.0.0',
      port: env.VITE_APP_PORT,
      proxy: {
          "api": {
            target: env.VITE_APP_API_BASEURL,
            changeOrigin: true,
            ws: true,
            rewrite: (path) => path.replace(/^\/api/, ""),
          }
      }
    },
    build: {
      rollupOptions: { // 配置rollup的一些构建策略
          // output: { // 控制输出
          //     // 在rollup里面, hash代表将你的文件名和文件内容进行组合计算得来的结果
          //     // assetFileNames: "[hash].[name].[ext]"
          // }
          // publicPath:"static",
         
      },
      assetsInlineLimit: 4096000, // 4000kb  超过会以base64字符串显示
      outDir: "website", // 输出名称
      // assetsDir: "static" // 静态资源目录
    },
  }
})
