<!-- Link to the issue this PR addresses -->
Closes #7

## Dependencies


- none needed

## What does this PR do?

Implements the backend part for role-based team visibility for daily meal participation. 

## Type of Change

<!-- Check one -->
- [X] New feature

## What was changed

- Employees can see the team they are assigned to and can view their own participation status.
- Team leads can view the pariticipation for members of their own team.
- Admins can view the participation across all teams
## Changelog

- Feature: Users can now see the participation status of themselves and the employees they lead based on their respective roles.


## How to Test

<!-- Describe how you tested your changes -->

- Manual testing
- Integration testing

## How QA Should Test

- Check if employee can see their participation status.
- Check if team leads can see their team members participation status.
- Check if admins can see the participation status of all the members. 

## Rollback Plan

- Revert this PR

## Checklist

- [X] My code follows the project style guidelines
- [X] All tests pass
- [X] I have updated documentation (if applicable)



## Note for Reviewer

- N/A