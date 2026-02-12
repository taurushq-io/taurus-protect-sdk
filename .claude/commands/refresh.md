# Review and Update CLAUDE Memory Files

**Instructions:**  
**Step 1: Get Overview**
List all CLAUDE.md and CLAUDE.local.md files in the project hierarchy.
**Step 2: Iterative Review**
Process each file systematically, starting with the root `CLAUDE.md` file:

- Load the current content
- Compare documented patterns against actual implementation
- Identify outdated, incorrect, or missing information
  **Step 3: Update and Refactor**
  For each memory file:
- Verify all technical claims against the current codebase
- Remove obsolete information
- Consolidate duplicate entries
- Ensure information is in the most appropriate file
  When information belongs to a specific subcomponent, ensure it's placed correctly:

* UI-specific patterns → `apps/myproject-ui/CLAUDE.md`
* API conventions → `apps/myproject-api/CLAUDE.md`
* Infrastructure details → `cdk/CLAUDE.md` or `infrastructure/CLAUDE.md`
  Focus on clarity, accuracy, and relevance. Remove any information that no longer serves the project.
