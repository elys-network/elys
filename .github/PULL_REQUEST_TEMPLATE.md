<!--
The default pull request template is for types feat, fix, or refactor.
For other templates, add one of the following parameters to the url:
- template=docs.md
- template=other.md
-->

## Description

Closes:
* https://github.com/elys-network/issues/issues/XXX
* https://github.com/elys-network/issues/issues/XXX

<!-- Add a description of the changes that this PR introduces and the files that
are the most critical to review. -->

### What has Changed?

What specific problem were you aiming to address, and how did you successfully resolve it? If tests were not uploaded for this pull request or if coverage decreased, please provide an explanation for the change.

---

### Author Checklist

_All items are required. Please add a note to the item if the item is not applicable and
please add links to any relevant follow up issues._

I have...

- [ ] included the correct [type prefix](https://github.com/commitizen/conventional-commit-types/blob/v3.0.0/index.json) in the PR title
- [ ] added `!` to the type prefix if API or client breaking change
- [ ] targeted the correct branch (see [PR Targeting](https://github.com/elys-network/elys/blob/main/CONTRIBUTING.md#pr-targeting))
- [ ] provided a link to the relevant issue or specification
- [ ] followed the guidelines for [building modules](https://github.com/elys-network/elys/blob/main/docs/docs/building-modules)
- [ ] included the necessary unit and integration [tests](https://github.com/elys-network/elys/blob/main/CONTRIBUTING.md#testing)
- [ ] included comments for [documenting Go code](https://blog.golang.org/godoc)
- [ ] updated the relevant documentation or specification
- [ ] reviewed "Files changed" and left comments if necessary
- [ ] confirmed all CI checks have passed

### Reviewers Checklist

_All items are required. Please add a note if the item is not applicable and please add
your handle next to the items reviewed if you only reviewed selected items._

I have...

- [ ] confirmed the correct [type prefix](https://github.com/commitizen/conventional-commit-types/blob/v3.0.0/index.json) in the PR title
- [ ] confirmed `!` in the type prefix if API or client breaking change
- [ ] confirmed all author checklist items have been addressed
- [ ] reviewed state machine logic
- [ ] reviewed API design and naming
- [ ] reviewed documentation is accurate
- [ ] reviewed tests and test coverage
- [ ] manually tested (if applicable)

### Deployment Notes

Are there any specific considerations to take into account when deploying these changes? This may include new dependencies, scripts that need to be executed, or any aspects that can only be evaluated in a deployed environment.

### Screenshots and Videos

Please provide any relevant before and after screenshots by uploading them here. Additionally, demo videos can be highly beneficial in demonstrating the process.
