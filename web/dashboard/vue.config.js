const path = require('path')

function resolve(dir) {
    return path.join(__dirname, dir)
}

let publicPath = process.env.VUE_APP_PUBLIC_PATH
module.exports = {
    outputDir: path.resolve(__dirname, '../../cmd/server/web/dashboard'),
    productionSourceMap: true,
    devServer: {
        port: 4400,
        open: true,
        overlay: {
            warnings: false,
            errors: true
        },
        proxy: {
            '/kubepi/api': {
                target: 'http://0.0.0.0:80',
                ws: true,
                secure: false,
            }
        }
    },
    configureWebpack: {
        devtool: 'source-map',
        resolve: {
            alias: {
                '@': resolve('src')
            }
        }
    },
    publicPath: `${publicPath}`
};
