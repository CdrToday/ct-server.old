const mail = require('./mail');
const redis = require('./utils').redis;
const { _a, _u } = require('./mongo');

class User {
  static async sendCode(ctx) {
    const _mail = ctx.params.mail;
    if (await mail(_mail)) {
      ctx.body = { msg: 'ok' };
    }
  }

  static async verifyCode(ctx) {
    const mail = ctx.params.mail;
    const code = ctx.request.body.code;
    const _code = await redis.get(mail);
    
    if(_code === code) {
      if (!_u.exists({mail: mail})) {
	await _u.create({mail: mail});
      }
      ctx.body = { msg: 'ok' };
    }
  }

  static async publish(ctx) {
    const mail = ctx.params.mail;
    const { title, content } = ctx.request.body;
    
    let res = await _a.create({
      title: title,
      content: content,
      timestamp: new Date().getTime()
    });
    let user = await _u.findOne({mail: mail});
    let articles = user.articles;
    articles.push(res._id);

    await _u.findOneAndUpdate({mail: mail}, {articles: articles});
    
    ctx.body = { msg: 'ok' };
  }
}

module.exports = User;
