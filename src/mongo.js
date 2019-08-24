const conf = require('./config');
const mongoose = require('mongoose');
const Schema = mongoose.Schema;
const ObjectId = mongoose.Schema.ObjectId;
const model = mongoose.model;

/* init */
mongoose.connect(conf.mongo.url, { useNewUrlParser: true });

/* user */
class User {
  static schema() {
    return Schema({
      name: {
	type: String,
	unique: true
      },
      mail: {
	type: String,
	unique: true
      },
      cmomunities: []
    })
  }

  constructor() {
    this.user = model('user', User.schema());
  }
}

/* community */
class Community {
  static schema() {
    return Schema({
      name: {
	type: String,
	unique: true
      }
    });
  }

  constructor() {
    this.community = model('community', Community.schema());
  }
}

/* exports */
const _u = new User().user;
const _c = new Community().community;

module.exports = { _u, _c };
