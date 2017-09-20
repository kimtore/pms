# How to contribute to PMS

Thank you for considering contributing to Practical Music Search.

You're welcome to report any bugs or request features using the Github issue
tracker. Code contributions are warmly received through pull requests on
Github. Make sure there is an open issue for the feature or bug you want to
tackle, and that it is coherent with the project strategy.

For general discussion about the project, or to contact the project devs, you
can use the IRC channel `#pms` on Freenode.

This project adheres to the
[Contributor Covenant Code of Conduct](code_of_conduct.md).
By participating, you are expected to uphold this code.

## Setting up the development environment

If you want to work on PMS, fork this repository and clone the fork to your
`$GOPATH` (`$GOPATH/src/github.com/[YOUR_GITHUB_USERNAME]/pms`). Due to
hardcoded dependencies, cloning your fork into `$GOPATH` is not sufficient.
Thus, you should also create a symlink to trick Go into looking for these
dependencies in your fork. For example:

```
cd $GOPATH/src/github.com
mkdir -p ambientsound
ln -s [YOUR_GITHUB_USERNAME]/pms ambientsound/pms
```

After this, `cd` into your fork and run `make`. Now the `pms` command is built
from your fork.
