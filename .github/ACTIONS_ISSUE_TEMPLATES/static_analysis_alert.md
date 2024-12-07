---
title: "Code Linting failed {{ env.COMMIT_HASH }}"
labels: technical debt
---

Static code analysis has failed for the commit [{{ env.COMMIT_HASH }}](https://github.com/{{ env.OWNER_REPOSITORY }}/commit/{{ env.COMMIT_HASH }}).

**Branch**: `{{ env.BRANCH_NAME }}`
**Commit**: `{{ env.COMMIT_HASH }}`

Please resolve the issues and push the changes to the branch `{{ env.BRANCH_NAME }}`.
