# Contributing

- [Contributing](#contributing)
  - [Architecture Decision Records (ADR)](#architecture-decision-records-adr)
  - [Development Procedure](#development-procedure)
    - [Testing](#testing)
    - [Pull Requests](#pull-requests)
    - [Pull Request Templates](#pull-request-templates)
    - [Requesting Reviews](#requesting-reviews)
    - [Updating Documentation](#updating-documentation)
  - [Branching Model and Release](#branching-model-and-release)
    - [PR Targeting](#pr-targeting)

Thank you for considering making contributions to the Elys and related repositories!

Contributing to this repo can mean many things, such as participating in
discussion or proposing code changes. To ensure a smooth workflow for all
contributors, the general procedure for contributing has been established:

1. Start by browsing [new issues](https://github.com/elys-network/elys/issues). If you have something in your mind, you can raise an issue ticket on github or initiate chat on Discord.
2. Determine whether a GitHub issue or discussion is more appropriate for your needs:
   1. If want to propose something new that requires specification or an additional design, or you would like to change a process, start with a new discussion. With discussions, we can better handle the design process using discussion threads. A discussion usually leads to one or more issues.
   2. If the issue you want addressed is a specific proposal or a bug, then open a [new issue](https://github.com/elys-network/elys/issues/new).
   3. Review existing [issues](https://github.com/elys-network/elys/issues) to find an issue you'd like to help with.
3. Participate in thoughtful discussion on that issue.
4. If you would like to contribute:
   1. Ensure that the proposal has been accepted.
   2. Ensure that nobody else has already begun working on this issue. If they have,
      make sure to contact them to collaborate.
   3. If nobody has been assigned for the issue and you would like to work on it,
      make a comment on the issue to inform the community of your intentions
      to begin work.
5. To submit your work as a contribution to the repository follow standard GitHub best practices. See [pull request guideline](#pull-requests) below.

**Note:** For very small or blatantly obvious problems such as typos, you are
not required to an open issue to submit a PR, but be aware that for more complex
problems/features, if a PR is opened before an adequate design discussion has
taken place in a GitHub issue, that PR runs a high likelihood of being rejected.

## Architecture Decision Records (ADR)

When proposing an architecture decision for Elys network, please start by opening an [issue](https://github.com/elys-network/elys/issues/new) or a discussion with a summary of the proposal. Once the proposal has been discussed and there is rough alignment on a high-level approach to the design, the ADR creation process can begin. We are following this process to ensure all involved parties are in agreement before any party begins coding the proposed implementation.

## Development Procedure

- The latest state of development is on `main`.
- `main` must never fail `go test ./...`.
- No `--force` onto `main` (except when reverting a broken commit, which should seldom happen).
- Create a branch to start work:
  - Fork the repo (core developers must create a branch directly in the Elys repo),
    branch from the HEAD of `main`, make some commits, and submit a PR to `main`.
  - For core developers working within the `elys` repo, follow branch name conventions to ensure a clear
    ownership of branches: `{moniker}/{issue#}-branch-name`.
  - See [Branching Model](#branching-model-and-release) for more details.
- Follow the [CODING GUIDELINES](CODING_GUIDELINES.md), which defines criteria for designing and coding a software.

Code is merged into main through pull request procedure.

### Testing

Tests can be executed by running `go test ./...` at the top level of Elys repository.

### Pull Requests

Before submitting a pull request:

- merge the latest main `git merge origin/main`,
- run `go test ./...` to ensure that all checks and tests pass.

Then:

1. If you have something to show, **start with a `Draft` PR**. It's good to have early validation of your work and we highly recommend this practice. A Draft PR also indicates to the community that the work is in progress.
   Draft PRs also helps the core team provide early feedback and ensure the work is in the right direction.
2. When the code is complete, change your PR from `Draft` to `Ready for Review`.
3. Go through the actions for each checkbox present in the PR template description. The PR actions are automatically provided for each new PR.

PRs must have a category prefix that is based on the type of changes being made (for example, `fix`, `feat`,
`refactor`, `docs`, and so on). The _type_ must be included in the PR title as a prefix (for example,
`fix: <description>`). This convention ensures that all changes that are committed to the base branch follow the
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.
Additionally, each PR should only address a single issue.

Pull requests are merged by one of core developers.

NOTE: when merging, GitHub will squash commits and rebase on top of the main.

### Pull Request Templates

There are three PR templates. The [default template](./.github/PULL_REQUEST_TEMPLATE.md) is for types `fix`, `feat`, and `refactor`. We also have a [docs template](./.github/PULL_REQUEST_TEMPLATE/docs.md) for documentation changes and an [other template](./.github/PULL_REQUEST_TEMPLATE/other.md) for changes that do not affect production code. When previewing a PR before it has been opened, you can change the template by adding one of the following parameters to the url:

- `template=docs.md`
- `template=other.md`

### Requesting Reviews

In order to accommodate the review process, the author of the PR must complete the author checklist
(from the pull request template)
to the best of their abilities before marking the PR as "Ready for Review". If you would like to
receive early feedback on the PR, open the PR as a "Draft" and leave a comment in the PR indicating
that you would like early feedback and tagging whoever you would like to receive feedback from.

All PRs require at least one review approval before they can be merged. Each PR template has a reviewers checklist that must be completed before the PR can be merged. Each reviewer is responsible
for all checked items unless they have indicated otherwise by leaving their handle next to specific
items. In addition, use the following review explanations:

- `LGTM` without an explicit approval means that the changes look good, but you haven't thoroughly reviewed the reviewer checklist items.
- `Approval` means that you have completed some or all of the reviewer checklist items. If you only reviewed selected items, you must add your handle next to the items that you have reviewed. In addition, follow these guidelines:
  - You must also think through anything which ought to be included but is not
  - You must think through whether any added code could be partially combined (DRYed) with existing code
  - You must think through any potential security issues or incentive-compatibility flaws introduced by the changes
  - Naming must be consistent with conventions and the rest of the codebase
  - Code must live in a reasonable location, considering dependency structures (for example, not importing testing modules in production code, or including example code modules in production code).
  - If you approve the PR, you are responsible for any issues mentioned here and any issues that should have been addressed after thoroughly reviewing the reviewer checklist items in the pull request template.
- If you sat down with the PR submitter and did a pairing review, add this information in the `Approval` or your PR comments.
- If you are only making "surface level" reviews, submit notes as a `comment` review.

### Updating Documentation

If you open a PR on Elys, it is mandatory to update the relevant documentation in `/docs`.

- If your changes relate to a module, then be sure to update the module's spec in `x/{moduleName}/README.md`.

## Branching Model and Release

User-facing repos should adhere to the trunk based development branching model: https://trunkbaseddevelopment.com. User branches should start with a user name, example: `{moniker}/{issue#}-branch-name`.

The Elys repository is a [multi Go module](https://github.com/golang/go/wiki/Modules#is-it-possible-to-add-a-module-to-a-multi-module-repository) repository. It means that we have more than one Go module in a single repository.

Elys utilizes [semantic versioning](https://semver.org/).

### PR Targeting

Ensure that you base and target your PR on the `main` branch.

All feature additions and all bug fixes must be targeted against `main`. Exception is for bug fixes which are only related to a released version. In that case, the related bug fix PRs must target against the release branch.

If needed, we backport a commit from `main` to a release branch (excluding consensus breaking feature, API breaking and similar).
