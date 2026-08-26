# Origami HyperFold (OHF): Scientific Foundations, Architecture, and Experimental Status

**Technical research document — August 2026**  
**Project:** Origami HyperFold (OHF)  
**Canonical implementation target:** Go  
**Reference visual carrier:** PNG  

## Abstract

Origami HyperFold (OHF) is an experimental protocol for **generative, addressable, selectively unfoldable and verifiable representation of structured information**. Its central goal is not merely to compress a source into a smaller byte stream, but to construct a representation in which a receiver can locate and reconstruct only the information required by a query while leaving unrelated regions folded.

The project combines ideas from lossless compression, Minimum Description Length (MDL), algorithmic information theory, grammar-based compression, deduplication, graph representation, authenticated data structures, random-access compressed data structures, perceptual coding, bounded generative execution, and evidence-based multimodal reasoning.

OHF deliberately distinguishes four questions that conventional compression often treats together:

1. **How compactly can a source be described?**
2. **Can an arbitrary part be addressed without expanding the whole source?**
3. **Can a model perceive and navigate a visual representation reliably?**
4. **Can exact recovery be independently verified?**

This separation leads to a layered architecture containing a canonical source model, content addressing, pattern discovery, a representation tournament, bounded Generative IR, SuperIndex, attention routing, multidimensional windowing, a visual compiler, a Perception Lab, separate perception/resolution/execution swarms, residual data, and a cryptographic Verification Spine.

OHF is a research system. Some components have experimental evidence; others remain architectural targets. In particular, extreme ratios such as **1:10,000** are hypotheses for highly structured or generatively describable domains, not universal compression guarantees. Shannon information limits are not bypassed: arbitrary incompressible data cannot be losslessly represented by a substantially shorter self-contained description.

---

## 1. Research question

The central OHF hypothesis is:

> To what extent can a generative, relational, addressable and perceptually navigable protocol describe a structured corpus using a representation substantially smaller than the expanded source, while allowing a receiver to unfold only the subset required for a task and recover exact content through deterministic generation, residual information and independent verification when exactness is requested?

The intended optimization is therefore not a single compression ratio. OHF studies a multidimensional objective:

\[
\text{Utility} = f(C, G, A, W, P, V)
\]

where:

- \(C\) = physical compression,
- \(G\) = generative leverage,
- \(A\) = addressability/random access,
- \(W\) = reduction of active working context,
- \(P\) = perceptual reliability,
- \(V\) = verifiability.

A representation that minimizes bytes but requires global decompression can be inferior, for OHF, to a slightly larger representation that allows cheap local reconstruction.

---

## 2. What OHF is — and is not

OHF is best understood as a **protocol for compact executable description and selective memory**, not as a PNG compression trick.

The logical system is:

```text
SOURCE
  ↓
INGEST + CANONICAL MODEL
  ↓
STRUCTURE / SEMANTICS / REPETITION ANALYSIS
  ↓
FOLD CANDIDATES
  ↓
REPRESENTATION TOURNAMENT
  ↓
GENERATIVE IR + RESIDUAL
  ↓
SUPERINDEX + DEPENDENCY MANIFESTS
  ↓
LOGICAL OHF CARRIER
  ↓
VISUAL COMPILER (optional transport layer)
  ↓
OHF IMAGE
```

A query reverses only the necessary portion:

```text
QUERY
  ↓
SUPERINDEX
  ↓
LOGICAL ADDRESSES
  ↓
ATTENTION ROUTER
  ↓
WINDOW SLIDERS
  ↓
MINIMAL DEPENDENCY CLOSURE
  ↓
RESOLVE
  ↓
EXECUTE
  ↓
APPLY RESIDUAL
  ↓
VERIFY
  ↓
ANSWER / EXACT DATA
```

OHF is **not** based on the claim that 500 KB can universally contain 5 GB of arbitrary lossless information. Such a claim would contradict basic information theory. A very large expansion ratio is possible only when the expanded object has a short effective description relative to its literal size, when shared state is explicitly allowed, or when some fidelity has intentionally been relaxed. OHF metrics must distinguish these cases.

---

## 3. Scientific foundation I — information theory and the impossibility boundary

Claude Shannon's information theory establishes the fundamental relationship between source statistics and achievable coding rates. Lossless compression exploits statistical regularity; it cannot assign uniquely decodable short representations to every possible long input.

This gives OHF its first invariant:

> **Generative leverage is not a violation of entropy.**

If a 5 GB object is generated exactly from a 500 KB self-contained carrier plus a fixed decoder, then that particular object has a compact description under that decoding system. This says something about the structure of the object; it does not imply that every 5 GB object admits such a description.

OHF therefore measures **description opportunity** rather than assuming compressibility. High-entropy residual regions remain literal or conventionally compressed.

---

## 4. Scientific foundation II — Kolmogorov complexity and Minimum Description Length

Kolmogorov complexity provides the conceptual ideal behind generative folding: the complexity of an object is related to the length of the shortest program that generates it. The shortest program is generally uncomputable, so OHF does not attempt to calculate Kolmogorov complexity directly.

Instead, OHF follows a practical idea closely related to **Minimum Description Length (MDL)**. Rissanen's MDL principle selects a model by considering the combined cost of describing the model and the data that remain unexplained by it.

For an OHF candidate representation \(M\) of region \(D\):

\[
L_{OHF}(D,M) = L(M) + L(D\mid M) + L_I + L_A + L_V + L_R
\]

where:

- \(L(M)\): cost of the rule/grammar/graph/transform,
- \(L(D\mid M)\): residual information required to reconstruct the region,
- \(L_I\): index cost,
- \(L_A\): expected attention/dependency cost,
- \(L_V\): verification metadata cost,
- \(L_R\): perceptual or operational risk penalty.

This is intentionally broader than conventional MDL. OHF cares not only about storage size but also about selective access.

A reusable rule is promoted only if its total system cost is lower than the literal alternative.

---

## 5. Scientific foundation III — compression by generation

Grammar-based compression replaces repeated substrings or structures with productions that generate them. Research on grammar-compressed strings demonstrates that compressed representations can support random access and substring extraction without decompressing the entire string.

OHF generalizes this principle beyond strings. Candidate descriptions may include:

- literals,
- references,
- dictionaries,
- grammars,
- templates with slots,
- transforms,
- repeated graph motifs,
- hierarchies,
- defaults plus overrides,
- prediction plus exact patch,
- residual blocks.

The key concept is **Generative Leverage**:

\[
GL = \frac{\text{bytes produced exactly}}{\text{bytes of generative description + required residual}}
\]

A repeated structure can have \(GL \gg 1\). Random or already-compressed content will generally have \(GL \approx 1\) or worse after metadata, in which case OHF should refuse the generative representation and retain a literal/residual representation.

---

## 6. Canonical source model

Before folding, OHF separates the source from its transport format.

Each source unit retains:

```text
SourceUnit
 ├─ identity
 ├─ path / logical location
 ├─ media type
 ├─ original bytes
 ├─ canonical representation
 ├─ original cryptographic digest
 └─ reverse patch when canonicalization is reversible
```

Canonicalization exists because semantically irrelevant variation can hide reusable structure. Examples include whitespace, line endings, formatting, serialization order, or representation-specific metadata.

However, OHF distinguishes two fidelity classes:

- **EXACT_CONTENT**: exact canonical content required by the application.
- **EXACT_SOURCE**: original source bytes must be reproduced byte-for-byte.

Canonicalization must never be silently presented as exact-source recovery. If exact-source recovery is required, all removed information must be preserved through a reversible patch or literal residual.

---

## 7. Content addressing and exact deduplication

Reusable content receives a content identifier derived from a cryptographic digest, conceptually:

\[
CID(x) = SHA256(x)
\]

Paths identify locations; CIDs identify content. Two different paths may therefore reference the same content.

Exact deduplication is the simplest folding operation:

```text
fragment
  ↓
CID
  ↓
already stored?
 ├─ yes → REF(CID)
 └─ no  → STORE(fragment)
```

The scientific basis is conventional content-addressed storage and hash-based deduplication. OHF extends it by allowing deduplicated nodes to participate in later grammars, graphs, transforms and dependency manifests.

---

## 8. Near deduplication and templates

Exact hashes cannot detect structures that are nearly identical. OHF therefore considers **near duplicates** as candidate templates.

For example:

```text
foo(x) = x + 10
bar(x) = x + 15
```

may be represented conceptually as:

```text
T(name, constant) := name(x) = x + constant
T(foo, 10)
T(bar, 15)
```

The template is accepted only when:

\[
Cost(T) + Cost(slots) + Cost(residual) < Cost(literals)
\]

Near similarity never licenses exact deduplication. The differing information remains explicit.

---

## 9. Source graph

OHF progressively converts source structure into graph views over shared identities.

Possible nodes include:

- files and directories,
- pages and sections,
- paragraphs,
- functions and AST nodes,
- records and schemas,
- entities,
- dictionary entries,
- rules,
- motifs.

Possible edges include:

- contains,
- references,
- imports,
- calls,
- defines,
- inherits,
- depends,
- precedes,
- similar-to,
- generated-by.

The graph is not itself a claim that graph encoding always compresses better. It provides a common relational substrate from which motifs, dependency closures, indexes and navigation can be derived.

---

## 10. Grammar, transform and motif discovery

### 10.1 Grammar discovery

Repeated symbolic structures are proposed as grammar rules. A rule candidate records definition cost, occurrence count, reference cost and literal alternative cost.

A simple savings estimate is:

\[
S = L_{literal} - (L_{definition} + nL_{reference} + L_{exceptions})
\]

Only positive and operationally useful candidates should survive.

### 10.2 Transform discovery

OHF can describe one object as a transformation of another when the transformation is cheaper than a literal copy. Candidate operators include COPY, REORDER, SUBSTITUTE, DELTA, PREFIX, SUFFIX, INSERT and DELETE; spatial domains may additionally admit geometric transforms.

### 10.3 Graph motif discovery

Repeated relational structures can be represented as parameterized motifs such as chains, stars, trees, diamonds, cycles, fan-in/fan-out structures or repeated subgraphs.

The representation must include all parameters and exceptions needed for deterministic expansion.

---

## 11. The Representation Tournament

A central OHF idea is that no single compression mechanism is assumed to be optimal.

For a region \(R\), several candidates compete:

```text
R
├─ Literal
├─ Reference
├─ Dictionary
├─ Grammar
├─ Template
├─ Transform
├─ Graph
├─ Motif
├─ Schema
└─ Predict + Patch
```

Each candidate produces a multidimensional cost vector:

```text
CarrierBits
ResidualBits
DefinitionBits
DependencyCost
ExpectedWindowCost
PerceptualCost
VerificationCost
DecodeCost
Risk
```

This is one of the main differences between OHF and a compressor optimized only for byte count.

A 20-byte representation requiring 50,000 dependencies may lose to a 40-byte representation requiring two local dependencies.

Local winners are subsequently checked by a global optimizer to detect dependency bombs, duplicated definitions, cycles, excessive fan-out and window inflation.

---

## 12. Generative IR

Winning fold representations compile to a finite, bounded intermediate representation rather than a general-purpose virtual machine.

Representative operations are:

```text
LITERAL
REF
CONCAT
REPEAT
SLICE
PATCH
TRANSFORM
RULE
GRAPH_EXPAND
MOTIF_EXPAND
DEFAULT
OVERRIDE
RESIDUAL
VERIFY
```

The IR is deliberately constrained:

- no arbitrary network access,
- no arbitrary system calls,
- no unrestricted filesystem mutation,
- no unbounded loops,
- explicit resource budgets.

Execution budgets include maximum instruction count, output bytes, recursion depth and expanded nodes.

This makes the carrier a bounded generative program rather than an unconstrained executable.

---

## 13. Residual information

Residual data are not a failure mode. They are mathematically necessary whenever a generative model does not fully determine the source.

OHF follows a two-part description:

\[
Source = Generate(Model) + Residual
\]

A practical representation may therefore be:

```text
prototype
+ transformation
+ exceptions
+ exact residual
```

The **Residual Ratio** is a critical metric. A domain with a high residual ratio offers little generative leverage even if sophisticated patterns can be discovered.

---

## 14. SuperIndex and addressability

Traditional decompression often begins from the start of a stream. OHF instead treats addressability as a first-class requirement.

The SuperIndex is a family of indexes over logical addresses, potentially including:

- PATH,
- PAGE / DOCUMENT IDENTITY,
- SEMANTIC,
- STRUCTURAL,
- ENTITY,
- DEPENDENCY,
- GRAPH,
- HIERARCHY,
- CYCLE,
- MOTIF,
- RULE,
- HASH,
- SCENE,
- TIME.

Indexes return addresses, not expanded content.

This is related to research on random access over compressed representations: compressed structure and selective access can coexist, but the representation must be designed for both.

---

## 15. Identity resolution: why offsets are not semantics

An important experimental lesson came from book-page tests. A constant offset between a container page index and a printed page number happened to work for one document but is not a valid general rule.

OHF therefore requires an explicit identity resolver capable of distinguishing:

- **ContainerIndex** — position in the physical container,
- **PDFLabel** — page label supplied by a PDF or similar format,
- **PrintedLabel** — number printed on the visible page,
- **LogicalID** — semantic identity intended by the query.

The resolver may be segmented because front matter, Roman numerals, appendices and inserted pages can change mappings within a document.

This converts accidental offsets into explicit metadata and prevents false exact retrieval.

---

## 16. Query planning

A query is classified before data are unfolded.

Conceptually:

```text
QueryPlan
 ├─ fidelity
 ├─ paths
 ├─ concepts
 ├─ entities
 ├─ relations
 └─ constraints
```

Fidelity is explicitly typed:

- SEMANTIC,
- STRUCTURAL,
- EXACT_CONTENT,
- EXACT_SOURCE.

The distinction matters. A semantic answer may tolerate abstraction; an exact-source request requires deterministic byte recovery and cryptographic verification.

---

## 17. Attention Router

The Attention Router does not answer the query. It selects which indexes and regions are likely to contain the answer.

For example:

```text
"return exactly internal/parser/parser.go"
        ↓
Fidelity = EXACT_SOURCE
Indexes  = PATH + DEPENDENCY + HASH
```

This is the first half of OHF's active-context principle: the represented corpus may be large, while the working set remains small.

---

## 18. Multidimensional Window Sliders

A window is not merely a rectangular crop of an image. It is a logical neighborhood with independent dimensions:

- spatial radius,
- semantic radius,
- graph radius,
- dependency depth,
- hierarchy up/down,
- order before/after,
- verification scope.

The core policy is:

> Failure to answer does not authorize global expansion.

Instead:

```text
query unresolved
  ↓
identify missing evidence/dependency class
  ↓
expand only corresponding slider
  ↓
retry
```

For example, a missing parent rule expands hierarchy upward; a missing referenced definition expands dependency depth. This is intended to prevent the common failure mode in which a selective system degenerates into loading the entire corpus.

---

## 19. Dependency manifests and dependency bombs

Every generative unit declares what it needs.

```text
Target
├─ Required dependencies
└─ Optional dependencies
```

Dependency traversal follows the direction of need. A global rule referenced by 50,000 nodes does not imply that all 50,000 dependents should be loaded when reconstructing one target.

OHF may intentionally duplicate small definitions when the storage penalty is lower than the expected attention cost of repeatedly crossing large dependency graphs.

This explicitly trades compression against locality.

---

## 20. Visual carrier as a transport layer

The visual layer is optional at the logical architecture level. PNG is the current reference visual carrier, but Folding IR must not depend on PNG.

A visual compiler receives:

```text
Logical OHF IR
+ Perceptual ABI
+ Attention Profile
        ↓
Visual Scene
        ↓
Renderer
        ↓
Carrier Image
```

The renderer does not invent semantics. Channel allocation occurs before rendering.

A future OHF implementation could use another transport while preserving the logical protocol.

---

## 21. BOOT, Rosetta and SuperIndex in a native carrier

A self-contained Native carrier must allow a fresh receiver to bootstrap from only:

1. the universal Master Prompt, and
2. one OHF image.

No external ABI, manifest, registry or helper file is a legitimate hidden dependency for Native mode.

The intended bootstrap sequence is:

```text
FIND ORIENTATION
→ FIND ANCHORS
→ READ PROTOCOL VERSION
→ FIND ROSETTA
→ LEARN CARRIER DIALECT
→ LOCATE SUPERINDEX
→ LOAD ROOT ADDRESSES
→ READY
```

**BOOT** contains the minimum universal orientation/version information.

**Rosetta** explains carrier-specific visual semantics that cannot safely be assumed universal.

**SuperIndex** exposes logical navigation roots.

The prompt defines the interpreter contract; the image defines the particular carrier's dialect and data.

---

## 22. Perception Lab

The visual language must be derived from measured model capabilities rather than human intuition.

The Perception Lab is the experimental subsystem that estimates those capabilities.

Its staged program is:

```text
P0    primitive capacity
P1    pair interactions
P1.5  correlated failures
P2    compounds
P2.5  promotion
P3    graphs
P3.5  hypergraphs / meta-relations
P4    degradation
P5    attention and distraction
P6    swarm behavior
P7    complete protocol
```

### P0 — Primitive capacity

Tests whether a model reliably distinguishes individual visual states: shape, topology, orientation, position, fill, markers and other candidate subdimensions.

### P1 — Interactions

A channel that works alone may fail when combined with another. P1 measures pairwise interference.

### P1.5 — Correlated failures

Redundancy is useful only when error sources are sufficiently independent. OHF therefore measures quantities conceptually equivalent to:

\[
P(A_{wrong}),\; P(B_{wrong}),\; P(A_{wrong}\cap B_{wrong})
\]

and especially the probability that two channels produce the **same wrong interpretation**.

### P2 — Compounds

Tests multi-feature symbols rather than assuming primitive accuracies compose multiplicatively.

### P3/P3.5 — Graphs and higher relations

Tests whether nodes, edges, direction, grouping, hierarchy, cycles, motifs and meta-relations survive actual VLM perception.

### P4 — Degradation

Tests resize, resampling, compression, noise and other transformations likely to occur in transport or model preprocessing.

### P5 — Attention

Measures how target recognition changes with distractors and whether local-window inspection restores accuracy.

### P6 — Swarm

Measures whether multiple narrow observations improve reliability and whether their errors are actually independent enough to justify consensus.

### P7 — Complete protocol

The valid Native P7 test is a fresh model session receiving exactly the Master Prompt plus one OHF image. If helper metadata or a decoder file is required, Native self-containment fails.

---

## 23. Perceptual ABI

Perception Lab results compile into a **Perceptual ABI** used by the encoder/compiler, not supplied as a third input to the receiving model.

Conceptually it contains:

- capability by visual subdimension,
- interaction matrix,
- evidence-dependence matrix,
- compound capabilities,
- graph capabilities,
- robustness envelope.

The visual compiler can then avoid channels known to be unreliable for a target model profile.

This architecture permits different carriers to use different visual dialects while keeping the Master Prompt relatively stable.

---

## 24. Why three swarms are separated

OHF treats probabilistic perception and deterministic truth as different computational stages.

```text
PERCEPTION SWARM
       ↓
Evidence
════════════════ PERCEPTION WALL
       ↓
RESOLUTION SWARM
       ↓
ResolvedValue
════════════════ EXECUTION GATE
       ↓
EXECUTION SWARM
       ↓
GeneratedValue
       ↓
Residual
       ↓
VERIFICATION SPINE
════════════════ COMMIT GATE
       ↓
VerifiedValue
```

### Perception Swarm

Observes. It may return candidates with confidence and provenance. It never produces cryptographic truth.

### Resolution Swarm

Combines candidate sets, constraints, grammar, graph structure, hierarchy and residual evidence. It does not inspect pixels directly.

### Execution Swarm

Executes bounded deterministic Generative IR after sufficient symbols have been resolved.

### Verification Spine

Checks independently supplied invariants such as lengths, CRCs, hashes, block hashes, Merkle proofs or source digests.

The separation prevents a dangerous epistemic shortcut:

> **Agreement is not verification.**

Ten models agreeing on the same value do not make it byte-exact.

---

## 25. Typed epistemic states

OHF proposes encoding epistemic distinctions in the implementation type system rather than relying only on prompt discipline.

Conceptually:

```text
CandidateValue
      ↓ resolution
ResolvedValue
      ↓ execution
GeneratedValue
      ↓ verification
VerifiedValue
```

A function requiring `VerifiedValue` should not accept `ResolvedValue`.

This turns the project's rule **FALSE EXACT = 0** into a software architecture property.

If exact verification cannot be completed, the correct output is UNKNOWN rather than a plausible reconstruction mislabeled as exact.

---

## 26. Verification Spine

Verification is independent of semantic confidence.

Potential levels are:

```text
length
CRC
local hash
block hash
Merkle proof
root hash
source hash
```

Cryptographic hash trees are especially relevant because they allow a small retrieved region to be authenticated against a known root without reading every unrelated region. This is conceptually aligned with OHF selective unfolding.

Expected verification values originate at encoding time. A decoder must not generate both the candidate and its supposed proof from the same untrusted reconstruction and then call that independent verification.

---

## 27. Three receiver profiles

OHF distinguishes three operational profiles.

### Native

```text
Master Prompt + OHF image + VLM
```

Used to test self-contained visual navigation. Exact capabilities are bounded by actual model perception.

### Computational

```text
OHF carrier + deterministic runtime
```

Used for exact addressing, bounded execution and cryptographic verification.

### Hybrid

```text
VLM + deterministic runtime
```

The VLM performs semantic routing and interpretation; deterministic software performs exact addressing, expansion and verification.

The Hybrid profile is expected to be operationally strongest, but Native remains an independent research target.

---

## 28. R2 exact carrier baseline

The current R2 exact-carrier experiment provides an important baseline but must not be confused with the final HyperFold architecture.

R2 can encode an exact compressed payload into deterministic binary visual cells with a visible bootstrap describing geometry and envelope structure. A capable environment can sample cells, reconstruct bytes and verify CRC/SHA-256.

This demonstrates a **visual transport / exact round-trip path** when pixel-level deterministic processing is available.

It does **not** by itself demonstrate:

- 1:10,000 generative folding,
- robust Native VLM decoding,
- full Representation Tournament optimization,
- cross-model Perceptual ABI,
- exact recovery from visual guessing,
- universal sub-500 KB carriers.

R2 should therefore be treated as an exactness and transport baseline from which more generative versions can be compared.

---

## 29. Experimental result: document random access

Book experiments demonstrated the usefulness of page/block indexes and selective block extraction, but they also exposed an important distinction between content recovery and source recovery.

A carrier built from canonicalized page text can support exact textual recovery relative to that canonical representation while still being unable to reconstruct the original PDF byte-for-byte, including images, fonts, layout streams and metadata.

Accordingly, all future reports should separately state:

```text
NativeSemanticCoverage
NativeStructuralCoverage
NativeExactContentCoverage
ComputationalExactContentCoverage
ComputationalExactSourceCoverage
```

This avoids inflating a successful text retrieval experiment into a stronger claim than the evidence supports.

---

## 30. Delta carriers and shared dictionaries

OHF's protocol model allows future carriers to depend explicitly on previously shared state.

```text
C0 = base carrier
C1 = changes relative to C0
C2 = changes relative to C1
```

Similarly, participants may share a dictionary identified by CID.

However, reported compression must distinguish:

- **SelfContainedRatio** — everything required is counted,
- **SharedStateRatio** — pre-existing state is excluded from transmitted bytes.

Failing to distinguish these would make compression claims scientifically misleading.

---

## 31. Metrics

A serious OHF benchmark should report at least:

```text
SourceBytes
CarrierBytes
SelfContainedRatio
SharedStateRatio
ResidualRatio
GrammarRatio
GraphRatio
TransformRatio
DedupRatio
GenerativeLeverage
NativeSemanticCoverage
NativeStructuralCoverage
NativeExactCoverage
ComputationalExactCoverage
MedianWindowNodes
P95WindowNodes
AttentionReductionRatio
EncodeLatency
DecodeLatency
RandomAccessLatency
FalseExactRate
```

The primary safety metric is:

\[
FalseExactRate = 0
\]

The central scaling experiment is the **Generative Leverage Curve**: measure description size, residual size, indexing overhead and query cost as source size grows within a domain.

---

## 32. Scientific falsifiability

OHF should be considered successful only to the degree supported by reproducible experiments.

Important falsification conditions include:

1. Generative representations do not beat literals after residual/index costs.
2. Random-access metadata eliminates the compression benefit.
3. Dependency closures routinely approach the full corpus.
4. Visual channels fail cross-model robustness tests.
5. Perceptual redundancy is strongly correlated and provides little effective error reduction.
6. Native decoding requires hidden helper artifacts.
7. Exact reconstruction fails cryptographic verification.
8. Extreme compression ratios disappear when shared state and omitted fidelity are counted honestly.

A negative result in the visual layer would not invalidate the logical OHF engine. It would invalidate or constrain the visual-native hypothesis.

---

## 33. Current status: demonstrated, partially demonstrated, proposed

### Demonstrated / baseline evidence

- deterministic exact carrier transport can be constructed and round-tripped when pixel processing is available;
- block/index based selective retrieval has been exercised on document experiments;
- cryptographic hashes can distinguish verified reconstruction from plausible reconstruction;
- Perception Lab infrastructure covers primitive, interaction, compound, graph, degradation, attention and swarm-style experiments;
- failures in Native model interpretation have already exposed real design weaknesses, including overreliance on text-like visual regions, loops in perception attempts and model-specific capability differences.

### Partially demonstrated / under calibration

- robust page identity resolution across heterogeneous documents;
- cross-model Native bootstrapping;
- stable visual compounds and graph navigation;
- correlation-aware redundancy;
- bounded local-window navigation under heavy distraction;
- selective exact retrieval performance across multiple document types.

### Proposed / not yet established as a general result

- full 42-family visual language validated across models;
- mature Representation Tournament across heterogeneous domains;
- production-quality global optimizer;
- large-scale grammar/graph/motif folding;
- 1:10,000 self-contained exact representation for real large corpora;
- multi-GB Generative Leverage Curve;
- mature delta-carrier ecosystem;
- OHF 1.0 cross-model Native protocol.

This separation is essential. Architecture is not evidence, and a successful experiment in one domain is not a universal result.

---

## 34. Implementation strategy

The canonical implementation target is Go.

The development order intentionally isolates deterministic correctness from visual uncertainty:

```text
1. Source model
2. Canonicalization
3. CID + exact dedup
4. Source graph
5. Fold candidates
6. Cost model
7. Generative IR
8. Deterministic exact round-trip
9. Grammar / transforms / motifs
10. Representation Tournament
11. Global optimization
12. Residual
13. SuperIndex
14. Query engine
15. Window Sliders
16. Attention benchmark
17. Perceptual registry
18. Perception Lab calibration
19. Perceptual ABI
20. Visual Compiler
21. Rosetta / Native bootstrap
22. Swarms
23. Verification Spine integration
24. Cross-model protocol tests
25. Delta carriers
```

The critical engineering principle is that visual perception must not be allowed to hide bugs in the logical encoder/decoder.

---

## 35. Threats to validity

### Dataset bias

Highly repetitive books, repositories or synthetic graphs may exaggerate generative leverage.

### Decoder prior knowledge

If a model already knows the source, apparent reconstruction may be memorization rather than carrier decoding.

### Hidden shared state

External prompts, dictionaries, runtimes or metadata can silently move information outside the measured carrier.

### Fidelity ambiguity

Semantic equivalence, canonical textual equality and byte-identical source reconstruction are different targets.

### Perceptual preprocessing

Image resizing and model-specific vision pipelines may alter the effective visual channel before the model sees it.

### Correlated model errors

Multiple agents using the same model family or visual representation may fail together, making apparent redundancy weaker than the number of votes suggests.

### Benchmark leakage

Repeatedly testing the same carriers while modifying prompts can overfit the protocol to a small benchmark set.

---

## 36. Relationship to established fields

OHF does not claim to invent the individual mathematical mechanisms on which it relies. Its research contribution, if validated, lies in their **composition around selective generative communication**:

| OHF mechanism | Established foundation |
|---|---|
| literal compression | source coding / conventional lossless compression |
| generative rules | algorithmic information theory / MDL |
| grammar folding | grammar-based compression |
| exact dedup | content-addressed storage / fingerprinting |
| motifs | graph compression / structural reuse |
| SuperIndex | indexing / random-access compressed structures |
| Window Sliders | bounded local retrieval / working-set control |
| residual | model + error/residual coding |
| hashes | cryptographic integrity |
| Merkle-style verification | authenticated data structures |
| Perception Lab | psychophysics-inspired controlled capability measurement / empirical VLM evaluation |
| evidence separation | probabilistic inference separated from deterministic verification |
| bounded IR | safe declarative execution / resource-bounded interpreters |

The novelty question must therefore be evaluated at the architectural and experimental level, not by claiming novelty for hashes, grammars, graphs or MDL individually.

---

## 37. Research roadmap

A defensible path to OHF 1.0 is:

### Phase A — exact logical engine

Prove byte-identical round-trip and selective file reconstruction without any visual layer.

### Phase B — generative folding

Measure grammar, transform, graph and motif savings on heterogeneous datasets.

### Phase C — addressability

Prove that random access and dependency windows remain bounded as corpus size grows.

### Phase D — perceptual calibration

Run P0–P6 across target VLMs and generate empirical Perceptual ABIs.

### Phase E — Native protocol

Test fresh-session `Master Prompt + one image` operation with no helper artifacts.

### Phase F — scaling

Construct Generative Leverage Curves from KB through MB and GB scales.

### Phase G — delta communication

Measure incremental carriers and shared-state scenarios separately from self-contained compression.

---

## 38. Conclusion

Origami HyperFold is an attempt to change the unit of optimization from **compressed file size alone** to **compact description plus selective reconstructability**.

Its strongest idea is not that an image magically stores an arbitrarily larger file. Its strongest idea is that structured information can sometimes be represented as a network of reusable generative descriptions, residuals, indexes and proofs, and that a receiver should unfold only the region demanded by the current task.

The resulting research problem sits at the intersection of compression, programmatic generation, graph representation, random access, multimodal perception and cryptographic verification.

The scientific standard for the project is therefore deliberately strict:

> A representation is only as compact as all information required to reconstruct it; a result is only exact when independently verified; and a visual feature is only part of the protocol when experiments show that target receivers can perceive it reliably.

Under those constraints, OHF's extreme compression objective becomes a measurable hypothesis rather than a claim of impossible universal compression.

---

## References and theoretical anchors

1. C. E. Shannon. **A Mathematical Theory of Communication.** *Bell System Technical Journal*, 1948.
2. J. Rissanen. **Modeling by Shortest Data Description.** *Automatica*, 14(5), 465–471, 1978. DOI: 10.1016/0005-1098(78)90005-5.
3. M. Li and P. Vitányi. Work on Kolmogorov complexity and Minimum Description Length; see also their analysis of MDL induction and algorithmic complexity.
4. P. Bille, G. M. Landau, R. Raman, K. Sadakane, S. R. Satti, O. Weimann. **Random Access to Grammar-Compressed Strings.** arXiv:1001.1565; later journal work on grammar-compressed random access.
5. C. Policroniades and I. Pratt. **Alternatives for Detecting Redundancy in Storage Systems Data.** USENIX Annual Technical Conference, 2004.
6. R. C. Merkle. Foundational work on hash trees and authenticated structures; subsequent authenticated-data-structure literature formalizes selective proofs over hashed trees.
7. A. Miller, M. Hicks, J. Katz, E. Shi. **Authenticated Data Structures, Generically.** POPL 2014.
8. Recent work on compression by code generation, including **The KoLMogorov Test: Compression by Code Generation**, is conceptually relevant to empirical measurement of how short executable descriptions can reproduce structured sequences.

### Project-specific note

The references above provide theoretical foundations and neighboring techniques. They should not be read as claiming that those works describe OHF itself. OHF's architecture, Representation Tournament, multidimensional Window Sliders, Perceptual ABI, swarm separation and visual-carrier protocol are project-specific design hypotheses that require their own empirical validation.
