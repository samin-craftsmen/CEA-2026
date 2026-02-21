
Closes #10

## Dependencies

- Merge PR #20

## What does this PR do?

- Allows Admin and Logistics roles to configure special day types that affect meal availability and participation behavior. Implements Only the backend.

## Type of Change

- [X] New feature

## What was changed

- **Admin**
  - Can mark a date as:
    - Office Closed
    - Government Holiday
    - Special Celebration Day (with optional note)
  - Can update or remove a special day setting.

## Changelog

Feature: Implemented Special Day Controls



## How to Test

1. Manual Testing
2. Integration Testing
3.

## How QA Should Test

- Choose a date as office closed and check meal participation.
- Choose government holiday and check only day status change.
- Choose special day celebration and see special day message

## Rollback Plan

- Revert this PR


## Checklist

- [x] Features work as expected
- [x] All tests pass
## Note for Reviewer

- Only backend implemented