const path = require('path')

function resolve(dir) {
    return path.join(__dirname, dir)
}

module.exports = {
    outputDir: path.resolve(__dirname, '../cmd/server/web/dashboard'),
    productionSourceMap: true,
    devServer: {
        port: 4300,
        open: true,
        overlay: {
            warnings: false,
            errors: true
        },
        proxy: {
            '/proxy': {
                target: 'http://0.0.0.0:2019',
                ws: true,
                secure: false,
            },
            '/api': {
                target: 'http://0.0.0.0:2019',
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
    publicPath: '/ui/',
};
