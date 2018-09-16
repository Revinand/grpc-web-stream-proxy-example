# gRPC-web example with channels streaming

### Get started (with HTTP 1.1)

* `npm install`
  * The postinstall will run `npm run get_go_deps` to install Golang dependencies. [dep](https://github.com/golang/dep) is used for go dependency management.
* `npm start` to start the Golang server and Webpack dev server
* Go to `http://localhost:8081`

### grpcwebproxy binary

`npm start` calls `grpcwebproxy` binary installed globally, e.g symlink at least is required - `ln -s /usr/local/bin $GOPATH/bin/grpcwebproxy`
