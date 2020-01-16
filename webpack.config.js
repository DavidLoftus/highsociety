const webpack = require('webpack');
const path = require('path');

module.exports = {
    mode: 'development',
    entry: {
        app: './public/entry.js',
    },
    devtool: 'inline-source-map',
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'static/js'),
    },
    module: {
        loaders: [
            {test: /\.css$/, loader: "style!css"},
            {
                test: /\.jsx$/,
                loaders: ['react-hot', 'babel'],
                include: [path.join(__dirname, 'public')]
            }
        ]
    },
};