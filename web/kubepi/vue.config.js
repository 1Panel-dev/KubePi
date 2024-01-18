const path = require('path')

function resolve(dir) {
    return path.join(__dirname, dir)
}

module.exports = {
    outputDir: path.resolve(__dirname, '../../cmd/server/web/kubepi'),
    productionSourceMap: true,
    lintOnSave: false,
    devServer: {
        port: 4300,
        open: true,
        proxy: {
            '/dashboard': {
                target: 'http://0.0.0.0:4400',
            },
            '/kubepi/api': {
                target: 'http://0.0.0.0:2019',
                ws: true,
                secure: false,
            },
            '/webkubectl': {
                target: 'http://0.0.0.0:2019',
                ws: true,
                secure: false,
            },
            '/terminal': {
                target: 'http://0.0.0.0:4200',
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
    publicPath: '/kubepi/',
};
