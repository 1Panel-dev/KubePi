const path = require('path')

function resolve(dir) {
    return path.join(__dirname, dir)
}

module.exports = {
    outputDir: path.resolve(__dirname, '../../cmd/server/web'),
    assetsDir: 'dashboard',
    indexPath: 'index.html',
    productionSourceMap: true,
    devServer: {
        port: 4400,
        open: true,
        overlay: {
            warnings: false,
            errors: true
        },
        proxy: {
            '/api': {
                target: 'http://0.0.0.0:2019',
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
    publicPath: '/dashboard/',
};
