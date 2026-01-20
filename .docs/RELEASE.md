# Release docs

This doc contains information for releasing a new version.

## Create a release

Creating a release can be done by pushing a tag to the GitHub repository (beginning with `v`).

```shell
VERSION="v0.0.1-alpha.1"
TAG=$VERSION

git tag $TAG -m "tag $TAG" -a
git push origin $TAG
```
