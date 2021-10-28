# ssh-login-notify


## deploy

install the binary to `/usr/local/bin/ssh-login-notify`

then edit `/etc/pam.d/sshd`

append this to last:

```
session optional pam_exec.so debug log=/tmp/sshd.log seteuid /usr/bin/env SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxx MAIL_FROM_NAME=ssh-guard MAIL_FROM=ssh-login-guard@example.com MAIL_TO=admin@example.com /usr/local/bin/ssh-login-notify
```
