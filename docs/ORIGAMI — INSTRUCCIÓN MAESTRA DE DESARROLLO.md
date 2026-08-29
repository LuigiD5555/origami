# ORIGAMI — INSTRUCCIÓN MAESTRA DE DESARROLLO

## PROPÓSITO

Desarrollar Origami como un lenguaje y una máquina experimental de representación de estados, relaciones, transformaciones, dinámica y percepción.

Origami NO debe reducirse a un formato de imagen, algoritmo de compresión, sistema OCR, QR gigante, VLM, LLM, base de datos, máquina virtual general ni sistema de agentes.

La representación visual es una posible proyección de Origami.

OHF — Origami HyperFold — es un perfil de carrier, memoria direccionable y laboratorio visual dentro de Origami. No es la totalidad de Origami.

Tlaloc es un sistema externo de orquestación y comportamiento. Origami debe poder funcionar independientemente de Tlaloc.

---

# FASE 1 — DEFINIR LA FUENTE DE VERDAD

Recibir una fuente.

La fuente puede ser:

texto,
PDF,
código,
datos estructurados,
grafo,
imagen,
estado sintético,
problema matemático
u otro dominio explícitamente soportado.

Generar una representación canónica.

Nunca permitir que la representación visual sustituya silenciosamente a la fuente.

Registrar hashes cuando exista una fuente exacta verificable.

Mantener separados:

EXACT_SOURCE

EXACT_CONTENT

SEMANTIC_STATE

PERCEPTUAL_PROJECTION

---

# FASE 2 — CONSTRUIR EL MODELO DE ESTADO

Transformar la fuente canónica en entidades y relaciones.

Representar como mínimo:

entidades,
atributos,
relaciones,
dependencias,
orden,
jerarquía,
contenimiento,
referencias,
grupos,
contexto.

Permitir relaciones N-arias cuando una relación binaria sea insuficiente.

La ausencia de una entidad, relación o estado puede poseer semántica y no debe tratarse automáticamente como cero o desconocido.

Mantener distintos:

ABSENT
UNKNOWN
INHIBITED
CANCELLED
UNSUPPORTED

---

# FASE 3 — DEFINIR LA MÁQUINA

Toda máquina Origami tiene:

S_t = estado actual

C_t = contexto declarado

R = reglas

F = función de transición

La evolución básica es:

S_(t+1) = F(S_t, C_t, R)

La semántica R0 debe ser determinista salvo que otro perfil declare explícitamente lo contrario.

Toda ejecución debe tener límites finitos cuando corresponda.

Detectar como mínimo:

FIXED_POINT
CYCLE
CONTRADICTION
BUDGET_EXHAUSTED

Conservar trazas reproducibles.

---

# FASE 4 — SOPORTAR ESTADOS COHERENTES

Implementar el perfil quantum-inspired sin afirmar computación cuántica física.

Estados:

DETERMINATE
SUPERPOSED
COUPLED
OBSERVED

TRANSFORM puede evolucionar ramas pero no seleccionar una.

OBSERVE constituye una frontera explícita de resolución.

INTERFERE combina las ramas declaradas y permite:

reinforcement,
partial cancellation,
total cancellation.

Una amplitud cero representa cancelación; no UNKNOWN.

Un estado COUPLED debe evaluarse como objeto conjunto hasta que una operación explícita autorice su separación.

UNFOLD expande alternativas o estructura.

FOLD restringe o resuelve según una política declarada.

---

# FASE 5 — DOMAR, NO ELIMINAR, LA EXPLOSIÓN COMBINATORIA

No reducir automáticamente todas las combinaciones posibles.

Las combinaciones pueden constituir semántica.

Explorar:

atractores,
ciclos,
oscilaciones,
metaestabilidad,
transiciones,
regímenes caóticos,
colapsos,
emergencia.

Usar estructura, simetría, equivalencias, índices, presupuestos y atención selectiva para controlar la exploración.

No destruir combinaciones únicamente porque sean numerosas.

---

# FASE 6 — SEPARAR ESTADO DE PERCEPCIÓN

Mantener como ley:

STATE != PERCEPT

Un estado válido puede existir aunque no sea perceptible mediante una observación determinada.

Mantener independientes:

COHERENT STATE

y

PERCEPTUAL AVAILABILITY.

---

# FASE 7 — IMPLEMENTAR CANALES PERCEPTUALES

Soportar como familias semánticas:

SPATIAL
INTERFERENCE
DEPTH
TEMPORAL
EMERGENT

Condiciones perceptuales:

LATENT
MASKED
MOTION_BOUND
PHASE_BOUND
STEREO_BOUND
EMERGENT
REVEALED
DECAYED

Implementar incrementalmente:

MOIRE
PHASE_SHIFT
STEREO_BIND
PARALLAX_RESOLVE
KINETIC_REVEAL
TEMPORAL_INTEGRATE
TEMPORAL_DECAY

Una operación declarada pero todavía no implementada debe devolver UNSUPPORTED.

Nunca aproximarla silenciosamente.

---

# FASE 8 — SOPORTAR TEMPORAL LATENT IMAGES

Permitir representaciones cuyo percepto no exista completamente en ningún frame estático.

Una Temporal Latent Image puede requerir:

carrier,
capas,
máscara,
fase,
trayectoria,
velocidad,
ventana temporal,
integración,
persistencia.

Una imagen ausente de un frame estático no constituye automáticamente evidencia de ausencia.

Pero toda afirmación temporal debe poder ser falsada mediante un Observation Contract finito.

---

# FASE 9 — EXIGIR OBSERVATION CONTRACTS

Toda afirmación que no sea directamente observable debe declarar antes de ejecutar:

claim_id,
subject,
preconditions,
observer,
trajectory,
budget,
expected predicate,
failure semantics.

El presupuesto debe ser finito.

Los únicos resultados válidos del contrato son:

PASS

FAIL

INVALID_CONTRACT

UNSUPPORTED

FAIL sólo ocurre cuando las precondiciones son válidas, se ejecuta el procedimiento autorizado y se agota el presupuesto sin satisfacer el predicado esperado.

No redefinir el criterio de éxito después del experimento.

---

# FASE 10 — CREAR REPRESENTACIONES GENERATIVAS

Cuando una fuente tenga estructura, competir entre distintas formas de representarla.

Candidatos posibles:

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

Seleccionar representaciones considerando:

tamaño,
complejidad,
profundidad de unfolding,
dependencias,
riesgo perceptual,
coste de ejecución,
coste de verificación.

No seleccionar únicamente por tamaño.

---

# FASE 11 — CONSTRUIR ÍNDICES DIRECCIONABLES

Toda representación grande debe poder consultarse sin desplegarla completamente.

Construir un SuperIndex capaz de direccionar, cuando sean aplicables:

páginas,
secciones,
bloques,
conceptos,
palabras,
entidades,
relaciones,
grafos,
jerarquías,
dependencias,
orden,
verificación.

La consulta debe producir direcciones lógicas.

Calcular después la mínima clausura de dependencias necesaria.

---

# FASE 12 — IMPLEMENTAR SELECTIVE UNFOLDING

Nunca utilizar por defecto:

carrier → unfold completo → búsqueda.

Utilizar:

query
→ SuperIndex
→ addresses
→ dependency closure
→ selective unfold
→ evidence.

Permitir ventanas independientes sobre:

espacio,
semántica,
radio de grafo,
profundidad jerárquica,
dependencias,
orden,
verificación.

---

# FASE 13 — PROYECTAR VISUALMENTE CUANDO SEA ÚTIL

Una proyección visual puede utilizar:

estado,
posición,
forma,
color,
orientación,
textura,
relaciones,
topología,
contenimiento,
orden,
grafo,
jerarquía,
composición,
interferencia,
profundidad,
tiempo.

No asumir que el producto cartesiano matemático de esas dimensiones constituye capacidad usable.

---

# FASE 14 — MEDIR CAPACIDAD PERCEPTUAL

Representar dimensiones mediante DVR.

Separar:

NOMINAL_PSS

de

SAFE_PSS.

Aplicar Perceptual Orthogonality:

una mutación controlada de una dimensión debe preservar las demás dimensiones que supuestamente permanecen constantes.

Mantener:

NO EVIDENCE = UNKNOWN

PAIRWISE PASS != HIGHER-ORDER PASS

UNKNOWN > FALSE EXACTNESS

---

# FASE 15 — MANTENER LA PERCEPTION WALL

La percepción sólo produce evidencia candidata.

Estados permitidos:

KNOWN
AMBIGUOUS
UNKNOWN
INVALID

La percepción no puede promover por sí sola un dato a verdad exacta.

Después de percepción ejecutar:

RESOLUTION

DETERMINISTIC EXECUTION

VERIFICATION.

---

# FASE 16 — VERIFICAR ANTES DE PROMOVER

Mantener la ley:

FALSE_EXACT = 0

Utilizar según corresponda:

length,
CRC,
block hash,
SHA-256,
Merkle proof,
root hash,
source hash,
invariants,
reference traces.

Sólo el verificador puede convertir un resultado candidato en exacto/verificado.

---

# FASE 17 — MANTENER OHF COMO PERFIL EXPERIMENTAL

OHF debe continuar investigando:

carrier compacto,
representación generativa,
SuperIndex,
selective unfolding,
DVR/PSS,
Glyph Calculus,
SAFE_MICRO_ISA,
Context SIMD,
Macro-Gestalt,
Native,
Computational,
Hybrid.

El objetivo <=500 KB pertenece actualmente a OHF.

No convertirlo automáticamente en una restricción de todo Origami.

No convertir el carrier en una VM de propósito general.

---

# FASE 18 — PROBAR ORIGAMI SIN IMÁGENES

Antes de atribuir una capacidad al lenguaje visual, comprobar que su semántica subyacente funciona independientemente.

Usar el Reference Engine para:

enumerar estados,
aplicar reglas,
producir trazas,
detectar contradicciones,
detectar ciclos,
detectar puntos fijos,
controlar presupuestos.

EXP-001 constituye la primera base de esta línea.

Después extender con experimentos independientes para:

interferencia,
coupling,
superposición,
observación,
Temporal Latent Images,
emergencia,
múltiples escalas temporales.

---

# FASE 19 — PROBAR ORIGAMI CON DOCUMENTOS

Para un PDF no intentar primero demostrar compresión extrema.

Construir:

PDF
→ canonical document
→ page structure
→ headings
→ tables
→ figures
→ references
→ semantic graph
→ document indexes
→ Origami state
→ optional OHF carrier.

Después aceptar preguntas arbitrarias sobre el documento.

Cada consulta debe realizar:

question
→ query interpretation
→ SuperIndex
→ relevant addresses
→ minimal dependency closure
→ selective unfold
→ evidence reconstruction
→ answer.

Cuando sea posible, conservar referencias de página y fuente.

Medir:

retrieval recall,
exact-content recovery,
semantic answer accuracy,
amount unfolded,
carrier cost,
verification rate,
false-exact rate.

El objetivo principal de esta prueba es demostrar:

“Origami puede convertir una fuente compleja en conocimiento direccionable y consultable sin desplegar toda la fuente”.

No:

“Origami mágicamente comprime cualquier archivo”.

---

# FASE 20 — CONTROLAR LOS EXPERIMENTOS

Todo experimento debe registrar antes de ejecutarse:

hipótesis,
fixture,
reglas,
seed,
presupuesto,
criterios PASS,
criterios FAIL.

Después registrar:

artefactos,
hashes,
trazas,
resultado,
regresiones,
impacto,
decisión.

Una hipótesis nunca debe modificarse para coincidir con un resultado.

---

# FASE 21 — PROMOVER SÓLO EVIDENCIA

Una idea puede encontrarse en uno de estos estados:

HYPOTHESIS
SPECIFIED
IMPLEMENTED
TESTED
EVIDENCE_SUPPORTED
PROMOTED

No confundirlos.

Que una operación aparezca en una especificación no significa que exista su runtime.

Que funcione matemáticamente no significa que sobreviva perceptualmente.

Que una pareja de dimensiones funcione no demuestra combinaciones de orden superior.

---

# FASE 22 — INTEGRACIÓN CON TLALOC

Mantener la frontera:

ORIGAMI
define estados, relaciones, representación, dinámica, observación y proyección.

TLALOC
coordina trabajo, agentes Tlaloque, comportamiento, entrenamiento, ejecución contra modelos, evaluación y promoción.

Tlaloc puede consumir un contrato versionado de Origami.

Tlaloc no puede redefinir su semántica.

Origami debe continuar siendo usable sin Tlaloc.

---

# FASE 23 — TONAL

Tonal debe funcionar como raíz de composición y política compartida cuando ambos proyectos se distribuyan conjuntamente.

No fusionar conceptualmente:

Tonal,
Tlaloc,
Origami.

Deben poder evolucionar y versionarse independientemente.

---

# FASE 24 — PRINCIPIO FINAL

Para cada nueva idea preguntar, en este orden:

1. ¿Qué estado representa?
2. ¿Qué relaciones existen?
3. ¿Qué reglas lo transforman?
4. ¿Cómo evoluciona en el tiempo?
5. ¿Qué propiedades son directas y cuáles emergentes?
6. ¿Cómo puede observarse?
7. ¿Cómo puede falsarse?
8. ¿Cómo se direcciona?
9. ¿Cómo se despliega únicamente lo necesario?
10. ¿Cómo se verifica?
11. ¿Qué evidencia permitiría promoverlo?

Sólo después preguntar:

“¿Cómo lo dibujamos?”

Porque la imagen es una proyección de Origami.

No Origami mismo.