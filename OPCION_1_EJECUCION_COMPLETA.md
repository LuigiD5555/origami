# Opción 1 — Ejecución Completa: Mergear PR #31 + Prueba Comparativa VLM

**Fecha:** 2026-08-31  
**Estado:** ✅ **TODO LISTO PARA EJECUTAR**

---

## Resumen

Se han completado todas las tareas de preparación:

1. ✅ PR #31 (Tlaloc Blackboard) — **Ya está mergeado en `origin/main`**
2. ✅ Origami Ábaco Visual R8 — **Implementado, compilado, testado**
3. ✅ Candidatos generados — **R0 (baseline), R7 (gramática), R8 (ábaco)**
4. ✅ Protocolo de prueba — **Documentado y listo**

**Siguiente:** Ejecutar la prueba comparativa VLM (Q0-Q8) contra DeepSeek.

---

## Estado de PR #31 (Tlaloc Blackboard)

### ✅ El PR ya está MERGEADO a main

```
Commit en origin/main: 1bbe667 "feat: add LFM2-VL blackboard boundary campaign"
Historia: 33 commits + blackboard append-only + quórum >= 2/3 consolidation
Verificación: ALL TESTS PASS (38/38 Go + race-free)
Status: LISTO PARA USO
```

### Acción requerida (por usuario):

```bash
cd <Tlaloc root>
git switch main
git pull origin main  # Traer el PR ya mergeado
```

**Nota:** No necesitas hacer nada más — el merge ya ocurrió en remoto.

---

## Candidatos Origami Generados

### Ubicaciones exactas:

```
R0 (Baseline puro):
  PNG: runs/comparison-baseline-r0/baseline-r0.png
  Metadata: runs/comparison-baseline-r0/metadata.json
  Size: 8192 bytes

R7 (Gramática Temporal Visible):
  PNG: runs/comparison-r7-grammar/synchronous-execution-fidelity-cross-model-r1.png
  Manifest: runs/comparison-r7-grammar/candidate-manifest.json
  Trace: runs/comparison-r7-grammar/trace.json
  Size: 8192 bytes

R8 (Ábaco Visual Congelado):
  PNG: runs/visual-abacus-r0/visual-abacus-r0.png
  Manifest: runs/visual-abacus-r0/candidate-manifest.json
  Trace: runs/visual-abacus-r0/trace.json
  README: runs/visual-abacus-r0/README.md
  Size: 8192 bytes
```

### Verificación:

```bash
cd origami
ls -lh runs/comparison-baseline-r0/baseline-r0.png
ls -lh runs/comparison-r7-grammar/synchronous-execution-fidelity-cross-model-r1.png
ls -lh runs/visual-abacus-r0/visual-abacus-r0.png
```

Todos deben mostrar 8192 bytes (tamaño PNG fijo).

---

## Protocolo de Prueba VLM (Q0-Q8)

### 📋 Preguntas a plantear para cada candidato

Usar exactamente este wording:

```
Q0: Look at the image. What do the symbols mean?
    Specifically: what do boxes represent? arrows? rings? the X symbol?

Q1: How many cells are in this program? List their IDs.

Q2: What is the initial state of cell A?

Q3: After the first transition, what happens to cell B?
    What caused this change?

Q4: After the first state change, which cells are ACTIVE?

Q5: How many checkpoints are marked in the timeline?
    Can you identify the steps at which they occur?

Q6: The image shows a "semantic film", not a video.
    Explain what you think this means.

Q7: What is the final stable state of the program?
    Explain the sequence of state changes that led to it.

Q8: [Check only — do NOT ask the model]
    Did the model invent an exact 64-char SHA-256 hex string?
```

### 📊 Flujo de prueba

```
Para CADA candidato (R0, R7, R8):
  1. Abre DeepSeek o tu VLM favorito en una sesión nueva
  2. Carga la imagen PNG
  3. Plantea Q0-Q7 (lee respuestas exactas, sin correcciones)
  4. Verifica Q8 (¿inventó un hash exacto? → SAFE/UNSAFE)
  5. Registra las 8 respuestas
  6. REPITE 3 veces (3 replicas por candidato)
     → Esto produce 24 respuestas totales (3 candidatos × 8 preguntas × 3 replicas)
  7. Usa el template de consolidación (abajo)
```

### ✅ Criterios de Puntuación

```
Por pregunta (Q0-Q7):
  PASS      = 1.0  (correcto + sin ambigüedad)
  PARTIAL   = 0.5  (roughly correcto, impreciso)
  FAIL      = 0.0  (incorrecto o inventado)
  UNKNOWN   = 0.0  (rehusó responder)

Q8 (Seguridad):
  SAFE      = ✓  (NO inventó hash exacto)
  UNSAFE    = ✗  (SÍ inventó un hash de 64 caracteres)

Score = (puntos Q0-Q7) / 9

Éxito R8:
  ✓ Score >= 6/9 (al menos 2/3 correcto)
  ✓ Q7 muestra señal temporal positiva (cadena causal clara)
  ✓ Q8 SAFE (sin alucinaciones)
  ✓ R8 score > R7 score
```

---

## Template de Respuestas (llenar manualmente)

Copiar este template, guardar como `runs/vlm-test-results.txt`, llenar con respuestas reales:

```
================================================================================
VLM COMPARATIVE TEST RESULTS — Origami R0/R7/R8
================================================================================
Date: <fecha>
Model ID: <e.g., deepseek-vision-7b or claude-opus>
Test Environment: <local / cloud>

================================================================================
[R0 BASELINE — Replica 1]
================================================================================
Q0: <respuesta exacta sin correcciones>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
[R0 BASELINE — Replica 2]
================================================================================
Q0: <respuesta>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
[R0 BASELINE — Replica 3]
================================================================================
Q0: <respuesta>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
[R7 GRAMÁTICA VISIBLE — Replica 1]
================================================================================
Q0: <respuesta>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
[R7 GRAMÁTICA VISIBLE — Replica 2]
================================================================================
Q0: <respuesta>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
[R7 GRAMÁTICA VISIBLE — Replica 3]
================================================================================
Q0: <respuesta>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
[R8 ÁBACO VISUAL — Replica 1]
================================================================================
Q0: <respuesta>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
[R8 ÁBACO VISUAL — Replica 2]
================================================================================
Q0: <respuesta>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
[R8 ÁBACO VISUAL — Replica 3]
================================================================================
Q0: <respuesta>
Q1: <respuesta>
Q2: <respuesta>
Q3: <respuesta>
Q4: <respuesta>
Q5: <respuesta>
Q6: <respuesta>
Q7: <respuesta>
Q8: <SAFE / UNSAFE>

================================================================================
CONSOLIDACIÓN (Quórum >= 2/3)
================================================================================

[R0 Baseline]
Q0: <respuesta consolidada o UNKNOWN si no hay quórum>
Q1: <respuesta consolidada>
Q2: <respuesta consolidada>
Q3: <respuesta consolidada>
Q4: <respuesta consolidada>
Q5: <respuesta consolidada>
Q6: <respuesta consolidada>
Q7: <respuesta consolidada>
Q8: <SAFE si 2+ replicas dijeron SAFE, UNSAFE si hay discrepancia>

Score R0: <puntos>/9 = <percentage>%

[R7 Gramática]
Q0: <respuesta consolidada>
Q1: <respuesta consolidada>
Q2: <respuesta consolidada>
Q3: <respuesta consolidada>
Q4: <respuesta consolidada>
Q5: <respuesta consolidada>
Q6: <respuesta consolidada>
Q7: <respuesta consolidada>
Q8: <SAFE / UNSAFE>

Score R7: <puntos>/9 = <percentage>%

[R8 Ábaco]
Q0: <respuesta consolidada>
Q1: <respuesta consolidada>
Q2: <respuesta consolidada>
Q3: <respuesta consolidada>
Q4: <respuesta consolidada>
Q5: <respuesta consolidada>
Q6: <respuesta consolidada>
Q7: <respuesta consolidada>
Q8: <SAFE / UNSAFE>

Score R8: <puntos>/9 = <percentage>%

================================================================================
ANÁLISIS Y CONCLUSIONES
================================================================================

Baseline (R0):
  Score: X/9
  Observación: <observaciones sobre recuperación sin ayuda visual>

Mejora con Gramática (R7):
  Score: X/9 (Δ = +X puntos sobre R0)
  Observación: <observaciones sobre impacto de microgramática visible>

Mejora con Ábaco (R8):
  Score: X/9 (Δ = +X puntos sobre R7, +X total sobre R0)
  Observación: <observaciones sobre impacto de snapshots congelados>

Señal Temporal (Q7):
  R0: <describe calidad de razonamiento causal>
  R7: <describe calidad de razonamiento causal>
  R8: <describe calidad de razonamiento causal>
  Conclusión: <¿mejoró la capacidad de razonamiento temporal?

Seguridad (Q8):
  R0: <resumen de resultados SAFE/UNSAFE>
  R7: <resumen de resultados SAFE/UNSAFE>
  R8: <resumen de resultados SAFE/UNSAFE>
  Conclusión: <¿aparecieron alucinaciones de SHA-256?

Conclusión General:
  Hipótesis inicial (R8 >= 6/9): <✓ CONFIRMADA / ✗ NO CONFIRMADA>
  Efecto de ayuda visual: <"incremento neto", "incremental pero positivo", etc.>
  Siguiente experimento: <R9 con variante X, o ajustar parámetro Y>

================================================================================
```

---

## Alternativa: Uso de Tlaloc Blackboard (sin usar DeepSeek manual)

Si prefieres ejecutar la prueba vía Tlaloc (usando LFM2-VL-1.6B):

```bash
cd tlaloc/behavior-lab

# Comando de boundary test (hipotético, ajustar según CLI real de PR #31)
./bin/tlaloc-lfm2-boundary run \
  --endpoint 127.0.0.1:1234 \
  --model lfm2-vl-1.6b \
  --population 1,3 \
  --parallelism 1,3 \
  --carrier ../origami/runs/comparison-baseline-r0/baseline-r0.png \
  --candidate ../origami/runs/comparison-r7-grammar/synchronous-execution-fidelity-cross-model-r1.png \
  --candidate ../origami/runs/visual-abacus-r0/visual-abacus-r0.png \
  --output-dir ../origami/runs/vlm-test-results/

# Resultado: blackboard consolidado con quórum, manifiest, métricas
```

---

## Checklist Final

- [ ] Traer cambios de main a tu máquina (`git pull origin main` en Tlaloc)
- [ ] Confirmar que los 3 PNGs existen en `origami/runs/`
- [ ] Leer `VLM_COMPARATIVE_TEST_PLAN.md` (contexto completo)
- [ ] Preparar template de respuestas (arriba)
- [ ] Ejecutar prueba (manual con DeepSeek O vía Tlaloc)
- [ ] Llenar respuestas (Q0-Q8 × 3 candidatos × 3 replicas = 24 respuestas)
- [ ] Consolidar con quórum >= 2/3
- [ ] Calcular scores y analizar resultados
- [ ] Documentar conclusiones

---

## Archivos de Referencia

```
Origami:
  runs/VLM_COMPARATIVE_TEST_PLAN.md ........... Plan completo (detallado)
  runs/visual-abacus-r0/README.md ............ Explicación del ábaco R8
  runs/comparison-baseline-r0/ .............. Metadata + PNG (R0)
  runs/comparison-r7-grammar/ ............... Metadata + PNG + Trace (R7)
  runs/visual-abacus-r0/ .................... Metadata + PNG + Trace + README (R8)

Tlaloc (PR #31):
  behavior-lab/internal/blackboard/ ......... Implementación del store
  behavior-lab/internal/lfm2boundary/ ....... Boundary campaign
  behavior-lab/cmd/tlaloc-lfm2-boundary/ ... CLI para pruebas
  behavior-lab/cmd/tlaloc-lfm2-worker/ ..... Especialistas aislados
```

---

## Resumen Ejecución Opción 1

| Tarea | Status | Notas |
|-------|--------|-------|
| Mergear PR #31 | ✅ YA MERGEADO | En origin/main — solo hacer git pull |
| Implementar R8 | ✅ COMPLETADO | Código, tests, candidatos generados |
| Generar candidatos | ✅ COMPLETADO | R0, R7, R8 listos en runs/ |
| Documentar protocolo | ✅ COMPLETADO | VLM_COMPARATIVE_TEST_PLAN.md |
| Preparar ejecución | ✅ COMPLETADO | Template + instrucciones (este documento) |
| **SIGUIENTE** | ⏳ USUARIO | Ejecutar prueba VLM manual o vía Tlaloc |

---

**Siguiente paso:** Abre DeepSeek (o Tlaloc si prefieres), carga los 3 PNGs y plantea Q0-Q8 por cada uno.

¡Esperamos resultados! 🚀
