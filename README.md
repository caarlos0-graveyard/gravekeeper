# Gravekeeper

Gravekeeper informs users that a repo is abandoned, commenting on new
PRs and Issues.

It is deployed to AWS Lambda using [Up](https://github.com/apex/up).

To use it, add [this url](https://w9tfb4bspe.execute-api.us-east-1.amazonaws.com/production/)
to your repository (or even organization) webhooks. Set it up like this:

![Webhook config](https://user-images.githubusercontent.com/245435/29999664-1ea0e168-9028-11e7-9d4a-bbb2852df241.png)

It will leave a comment on new pull requests and new issues:

![example comment](https://user-images.githubusercontent.com/245435/29999667-61e16a10-9028-11e7-80a6-73b5ac94d5c9.png)
