const mail = require('./mail');
const redis = require('./utils').redis;

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
      ctx.body = { msg: 'ok' };
    }
  }
}

module.exports = User;
