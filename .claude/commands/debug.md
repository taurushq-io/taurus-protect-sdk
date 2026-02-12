# Debugger Mode

You are in DEBUGGER mode. Investigate the bug and identify the root cause.

## Your Task

1. **Understand the problem**: Expected vs actual behavior
2. **Form hypotheses**: List possible causes ranked by likelihood
3. **Investigate**: Trace through the code path
4. **Identify root cause**: Not just WHERE but WHY
5. **Recommend fix**: Specific, actionable guidance

## Common Bug Patterns

**Null/Undefined Reference**
- Missing null check
- Async timing issue (data not loaded yet)
- Failed API call not handled
- Incorrect assumption about data shape

**Off-by-One**
- `<` vs `<=` confusion
- Index starting at 0 vs 1
- Inclusive vs exclusive ranges

**Race Conditions**
- Missing synchronization
- Assumed order of async operations
- Shared mutable state

**State Issues**
- Mutating state directly
- Stale closures
- Multiple sources of truth

## Output Format

```
## Bug Analysis

**Issue**: [One-line description]
**Severity**: [Critical / High / Medium / Low]

## Problem Statement
**Expected**: [What should happen]
**Actual**: [What happens instead]

## Hypotheses
| # | Hypothesis | Likelihood |
|---|------------|------------|
| 1 | [Most likely cause] | High |
| 2 | [Another possibility] | Medium |

## Investigation
[Walkthrough of the code path and findings]

## Root Cause
**Proximate cause**: [Immediate trigger]
**Root cause**: [Underlying issue]

## Recommended Fix
[Specific code changes with examples]

## Verification
[How to verify the fix works]

## Prevention
[How to prevent similar bugs]
```

## Bug to Investigate

$ARGUMENTS
