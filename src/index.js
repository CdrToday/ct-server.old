const Koa = require('koa');
const cors = require('@koa/cors');
const logger = require('koa-logger');
const Router = require('koa-router');
const bodyparser = require('koa-body');
const conf = require('./config');
const user = require('./user');
const articles = require('./articles');

// main
class Index {
  static router() {
    let r = new Router();

    r.all('/', ctx => { ctx.body = 'hello, world';})
      .post('/api_v0/:mail/code', user.sendCode)
      .post('/api_v0/:mail/verify', user.verifyCode)
      .post('/api_v0/:mail/publish', user.publish)
      .get('/api_v0/:mail/articles', articles.articles)

    return r;
  }

  static server(r) {
    const server = new Koa();
    console.log(`Server listen to ${conf.server.port}...`);

    server
      .use(cors())
      .use(logger())
      .use(bodyparser())
      .use(r.routes())
      .use(r.allowedMethods())
      .listen(conf.server.port);
  }

  static main() {
    Index.server(Index.router());
  }
}

Index.main();
