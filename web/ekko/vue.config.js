const path = require('path')

function resolve(dir) {
    return path.join(__dirname, dir)
}

module.exports = {
    outputDir: path.resolve(__dirname, '../cmd/server/web/ekko'),
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
            },
            '/dashboard': {
                target: 'http://0.0.0.0:4400',
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
    publicPath: '/ekko/',
};
