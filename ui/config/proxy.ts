export default {
  dev: {
    "/api/": {
      // target: "http://127.0.0.1:9001",
      target: "http://172.16.21.219:9001",
      changeOrigin: true,
      pathRewrite: { "^": "" },
    },
  },
  test: {},
  pre: {},
};