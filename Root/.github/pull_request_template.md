<!-- Link to the issue this PR addresses -->
Closes #13

## Dependencies

- Merge PR #22

## What does this PR do?

- Allows admin to generate a copy/paste ready annoucement message for a selected date, summarizing meal participation and special-day notice.

## Type of Change

- [x] New feature

## What was changed

- **Admin**
  - Can select a date.
  - Can generate a formatted announcement message.
  - Can copy the message easily for sharing (e.g., Slack, Email).

- **Message Content Includes**
  - Meal-wise totals 
  - Overall headcount
  - Office vs WFH split (if applicable)
  - Special day status:
    - Office Closed
    - Government Holiday
    - Special Celebration Day (including note)

## Changelog


- Feature: Daily announcement draft



## How to Test

<!-- Describe how you tested your changes -->

1. Manual Testing
2. Integration Testing

## How QA Should Test

- Log in as admin
- Select a date
- Copy the message that includes employee informations

## Rollback Plan

- Revert this PR


## Checklist

- [ ] All requirements have been met
- [ ] All tests pass


## Note for Reviewer

- Only backend has been implemented