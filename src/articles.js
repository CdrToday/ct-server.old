const { _a, _u } = require('./mongo');

class Articles {
  static async articles(ctx) {
    const mail = ctx.params.mail;
    
    let user = await _u.findOne({mail: mail});
    let list = await _a.find({ '_id': {$in: user.articles }});
    let res = list.map(e => {
      return {
	_id: e._id,
	title: e.title,
	content: e.content
      };
    })
    
    ctx.body = {
      articles: res
    };
  }
}

module.exports = Articles;
