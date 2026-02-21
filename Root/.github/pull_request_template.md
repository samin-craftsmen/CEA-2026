<!-- Link to the issue this PR addresses -->
Closes #9

## Dependencies

- Merge PR #8

## What does this PR do?

- Implements the backend part for bulk and exception handling for daily meal participation. Team Leads/Admin can apply bulk actions for their scope (e.g., mark a group as opted out due to offsite/event). It also refactors the code to make it more maintainable

## Type of Change

<!-- Check one -->
- [X] New feature
- [X] Refactor

## What was changed

- Team lead can now apply bulk actions to his own team members.
- Admin can now apply bulk actions across all teams

## Changelog

- Feature: Implement Bulk and exception handling


## How to Test

<!-- Describe how you tested your changes -->

1. Manual Testing
2. Integration Testing

## How QA Should Test

- Log in as team lead and apply bulk actions to his own team members
- Log in as admin and apply bulk actions to members across all teams.

## Rollback Plan

- Revert this PR

## Checklist

- [X] My code follows the project style guidelines
- [X] All tests pass
- [X] The features work as expected

## Note for Reviewer

- This only implements the backend part.