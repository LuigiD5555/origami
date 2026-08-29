# Legacy OHF R2 material

This repository originally presented **Origami HyperFold (OHF)** as the complete project and used an R2 visual carrier plus a reusable receiver prompt as its main experimental baseline.

That material is preserved for research continuity, but it predates the current architectural split:

- **Tlaloc** owns orchestration, behavior compilation, Tlaloque coordination/training and system-level verification.
- **Origami** owns representation, state semantics, dynamics and projection.

The following files are therefore historical references rather than current architecture specifications:

- `docs/SCIENTIFIC_FOUNDATIONS.md`
- `prompts/OHF_R2_MASTER_PROMPT.txt`

Concepts from those documents may later be promoted into current Origami or Tlaloc only through explicit change control and regression testing. Their presence in the repository does not make them current requirements.

The current authoritative Origami documents are:

- `README.md`
- `PROJECT_BOUNDARY.md`
- `docs/STATE_SEMANTICS_R0.md`
- `CHANGE_CONTROL_6.0.0-alpha.1.md`
