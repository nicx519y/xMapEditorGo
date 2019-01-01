const path = require('path');

module.exports = {
    mode: 'development',

    entry: './src/index.ts',
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'dist')
    },

    module: {
        rules: [
            {
                test: /\.ts$/,
                use: "ts-loader"
            },
        ],
    },
    resolve: {
        extensions: [
            '.ts', '.html', '.js',
        ]
	},
	devServer: {
		port: '9090',
		inline: true,
        hot: false,
        watchContentBase: true,
    },
};