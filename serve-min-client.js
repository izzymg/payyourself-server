const http = require("http")
const fs = require("fs")

const server = http.createServer((req, res) => {
    const stream = fs.createReadStream("./min-client/index.html")
    stream.pipe(res, { end: true })
})

server.listen(3000, "0.0.0.0")