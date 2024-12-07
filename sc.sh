REPO=hasAnybodySeenHarry/throttler
CURRENT_BRANCH=main
CURRENT_COMMIT=24056b0c4f4999e7e405fe88921d0e79e6512cfc

ISSUES=$(gh issue list --repo "$REPO" --label "technical debt" --state open --json number,title,body 2>/dev/null || echo '[]')

ISSUES=$(echo "$ISSUES" | tr -d '\000-\031')

echo $ISSUES

echo "$ISSUES" | jq -c '.[]' | while read -r ISSUE; do
    ISSUE_NUMBER=$(echo "$ISSUE" | jq -r '.number')
    ISSUE_TITLE=$(echo "$ISSUE" | jq -r '.title')
    ISSUE_BODY=$(echo "$ISSUE" | jq -r '.body')

    echo "Processing issue: Number=$ISSUE_NUMBER, Title=$ISSUE_TITLE"

    if [[ "$ISSUE_TITLE" == "Code Linting failed"* ]]; then
        FAILED_BRANCH=$(echo "$ISSUE_BODY" | grep -oE '\*\*Branch\*\*: `[^\`]+' | sed 's/.*`//')
        FAILED_COMMIT=$(echo "$ISSUE_BODY" | grep -oE '\*\*Commit\*\*: `[^\`]+' | sed 's/.*`//')

        echo "Extracted Branch=$FAILED_BRANCH, Commit=$FAILED_COMMIT"

        if [[ "$FAILED_BRANCH" == "$CURRENT_BRANCH" ]]; then
        if git merge-base --is-ancestor "$FAILED_COMMIT" "$CURRENT_COMMIT"; then
            echo "Matching issue found: #$ISSUE_NUMBER"
            gh issue close "$ISSUE_NUMBER" --repo "$REPO" --comment "The issue has been resolved in commit $CURRENT_COMMIT."
        else
            echo "Commit is not a descendant of the current commit."
        fi
        else
        echo "Branch does not match."
        fi
    else
        echo "Skipping issue as it doesn't match the expected title."
    fi
done