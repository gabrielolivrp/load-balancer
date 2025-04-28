const http = require("node:http");

if (!process.env.PORT) {
  throw new Error("Provide the PORT number");
}

const port = process.env.PORT;

const healthCheckHandler = (_req, res) => {
  res.writeHead(200, { "Content-Type": "application/json" });
  res.end(JSON.stringify({ status: "ok" }));
};

const indexHandler = (_req, res) => {
  res.writeHead(200, { "Content-Type": "application/json" });
  res.end(JSON.stringify({ host: `http://localhost:${port}` }));
};

const notFoundHandler = (_req, res) => {
  res.writeHead(404, { "Content-Type": "application/json" });
  res.end(JSON.stringify({ message: "Not Found" }));
};

const server = http.createServer((req, res) => {
  switch (req.url) {
    case "/health":
      return healthCheckHandler(req, res);
    case "/":
      return indexHandler(req, res);
    default:
      return notFoundHandler(req, res);
  }
});

server.listen(port, () => {
  console.log(`Server is running on http://localhost:${port}`);
});
