Closes #33

## Dependencies

- Merge PR #34

## What does this PR do?

- Allows Admin to declare a date range as a company-wide “WFH for everyone” period. During this period, employees are treated as WFH by default for reporting and headcount calculations.

## Type of Change

- [x] New feature

## What was changed

- **Admin**
  - Can define a start date and end date for a WFH period.
  - Can update or remove an existing WFH period.

- **System Behavior**
  - During the declared period:
    - Employees are treated as WFH by default.
    - Headcount reporting reflects all users as WFH unless explicitly overridden (if allowed).
    - Meals are opted out
  - Outside the declared period:
    - Normal work-location rules apply.


## Changelog


Feature: Company-wide WFH period



## How to Test

1. Manual Testing
2. Integration Testing

## How QA Should Test

- Create a WFH period
- Check employee meal and work status

## Rollback Plan

- Revert this PR


## Checklist

- [x] All requirements have been met
- [x] All tests pass

## Note for Reviewer

- Frontend has been implemented.