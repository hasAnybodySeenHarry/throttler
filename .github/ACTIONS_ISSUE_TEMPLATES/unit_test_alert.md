---
title: "Unit Testing failed {{ env.COMMIT_HASH }}"
labels: bug
---

Unit tests have failed for the commit [{{ env.COMMIT_HASH }}](https://github.com/{{ env.OWNER_REPOSITORY }}/commit/{{ env.COMMIT_HASH }}).

**Branch**: `{{ env.BRANCH_NAME }}`
**Commit**: `{{ env.COMMIT_HASH }}`

Please fix the issue and push the changes to the branch `{{ env.BRANCH_NAME }}`.
