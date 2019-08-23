const nodemailer = require('nodemailer');
const uuid = require('uuid/v4');
const conf = require('./config');
const redis = require('./utils').redis;



async function mail(to) {
  let _uuid = uuid();
  
  let transporter = nodemailer.createTransport({
    host: 'smtp.qq.com',
    secure: false,
    auth: {
      user: conf.mail.user,
      pass: conf.mail.pass
    }
  });

  await transporter.sendMail({
    from: conf.mail.from,
    to: to,
    subject: conf.mail.subject,
    text: _uuid,
    html: `<b>${_uuid}</b>`
  });

  await redis.set(to, _uuid);
  return await transporter.verify();
}

module.exports = mail;
