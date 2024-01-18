const path = require('path')

function resolve(dir) {
    return path.join(__dirname, dir)
}

let publicPath = process.env.VUE_APP_PUBLIC_PATH
module.exports = {
    outputDir: path.resolve(__dirname, '../../cmd/server/web/dashboard'),
    productionSourceMap: true,
    lintOnSave: false,
    devServer: {
        port: 4400,
        open: true,
        proxy: {
            '/kubepi/api': {
                target: 'http://0.0.0.0:80',
                ws: true,
                secure: false,
            }
        }
    },
    configureWebpack: {
        performance: { 
            hints: false,
        },
        optimization: {
            minimize: true,
            splitChunks: {
                chunks: 'all',
            }
        },
        devtool: 'source-map',
        resolve: {
            alias: {
                '@': resolve('src')
            },
            fallback: {
                "path": false,
            }
        }
    },
    publicPath: `${publicPath}`
};
