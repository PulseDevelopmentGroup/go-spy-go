const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');
const webpack = require('webpack');

const htmlPlugin = new HtmlWebpackPlugin({
  template: './src/index.html',
  filename: 'index.html',
});

const EnvPlugin = new webpack.DefinePlugin({
  'process.env.API_URL': JSON.stringify(
    process.env.API_URL || '127.0.0.1:8080'
  ),
});

module.exports = {
  devtool: 'source-map',
  entry: './src/index.js',
  output: {
    filename: '[name].bundle.js',
    path: path.resolve(__dirname, '../spyfall-server/public'),
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
        },
      },
      {
        test: /\.css$/,
        use: ['syle-loader', 'css-loader'],
      },
    ],
  },
  resolve: {
    extensions: ['*', '.js', '.jsx'],
  },
  devServer: {
    host: '127.0.0.1',
  },
  plugins: [htmlPlugin, EnvPlugin],
};
