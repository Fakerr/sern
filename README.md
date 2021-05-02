# Sern #

Sern is a Github bot that helps you manage your merges and prevents merge skew/semantic merge conflicts so your master branch stays always green.


### Motivation ###
If you have two pull requests that modify dependent code, the tests could pass on each pull request independently and Github would allow you to merge the PRs,
however, the build may break after the merge.

One option to fix that would be to block PRs that are not up to date from being merged and sync each time you would need to,
but this would mean you will move responsibility to devs to continuously watch over PRs status in order to keep things up to date which is not efficient.


### How it works ###
In order to keep master green, Sern follow those two rules:

- Make sure all branch builds are run whilst rebased against latest master before merging
- Make sure only one build can be merged at a time

For that, Sern provides a simple queue and a way to run a custom build against master.
Instead of hitting "merge", you'll get Sern to add your PR on a queue which runs one last sanity check build on a special staging branch.
The build is the equivalent of what would have run if the merge happened. If it goes green, the PR is merged automatically.


### No longer maintained ###
This project is not finished and is now archived. Many alternatives of pull requests automation branch update and merge that are more mature already exist
(ex: [Kodiak](https://github.com/chdsbd/kodiak), [mergify](https://github.com/Mergifyio/mergify-engine) , etc...)
