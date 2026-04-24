const path = require("path");
const TerserPlugin = require("terser-webpack-plugin");

module.exports = {
    mode: "production",
    entry: "./src/main.ts",
    output: {
        path: path.resolve(__dirname, "dist"),
        filename: "gotty-bundle.js"
    },
    devtool: false,
    optimization: {
        minimizer: [
            new TerserPlugin({
                extractComments: false
            })
        ]
    },
    resolve: {
        extensions: [".ts", ".tsx", ".js"],
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                loader: "ts-loader",
                exclude: /node_modules/
            },
            {
                test: /\.js$/,
                include: /node_modules/,
                loader: 'license-loader'
            }
        ]
    }
};
