Closes #17

## Dependencies

- Merge PR #19

## What does this PR do?

- Implements the frontend part for bulk and exception handling for daily meal participation. Team Leads/Admin can apply bulk actions for their scope (e.g., mark a group as opted out due to offsite/event). 

## Type of Change

- [X] New feature

## What was changed

- Team lead can now apply bulk actions to his own team members.
- Admin can now apply bulk actions across all teams

## Changelog

- Feature: Implement Bulk and exception handling


## How to Test

1. Manual Testing
2. Integration Testing

## How QA Should Test

- Log in as team lead and apply bulk actions to his own team members
- Log in as admin and apply bulk actions to members across all teams.

## Rollback Plan

- Revert this PR

## Checklist

- [X] All features work as expected
- [X] All tests pass

## Note for Reviewer

- Frontend is integrated.