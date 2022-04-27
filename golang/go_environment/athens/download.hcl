downloadURL = "https://goproxy.cn"

mode = "async_redirect"

download "gitlab.test.com/*" {
    mode = "sync"
}

