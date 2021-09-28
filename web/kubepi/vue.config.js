const path = require('path')

function resolve(dir) {
    return path.join(__dirname, dir)
}

module.exports = {
    outputDir: path.resolve(__dirname, '../../cmd/server/web/kubepi'),
    indexPath: 'index.html',
    productionSourceMap: true,
    devServer: {
        port: 4300,
        open: true,
        overlay: {
            warnings: false,
            errors: true
        },
        proxy: {
            '/api': {
                target: 'http://0.0.0.0:80',
                ws: true,
                secure: false,
            },
            '/dashboard': {
                target: 'http://0.0.0.0:4400',
            },
            '/terminal': {
                target: 'http://0.0.0.0:4200',
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
    publicPath: '/kubepi/kubepi/',
};
