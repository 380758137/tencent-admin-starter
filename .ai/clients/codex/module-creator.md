# Codex Prompt Template

Use this project workflow for module generation:

1. Convert requirement into `specs/modules/<module>.module.spec.json`
2. Run: `npm run module:plan -- --spec specs/modules/<module>.module.spec.json`
3. After confirmation run: `npm run module:apply -- --spec specs/modules/<module>.module.spec.json`

