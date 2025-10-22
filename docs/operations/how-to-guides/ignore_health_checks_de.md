## Health-Checks ignorieren

Vorgeschaltete Health-Checks können deaktiviert werden:
- für Dogus, wenn `spec.ignoreDoguHealth` auf `true` gesetzt ist.

Dies ermöglicht es, Fehler an Dogus via Blueprint zu beheben.
Für ein Dogu-Upgrade muss ein Dogu jedoch healthy sein, um Pre-Upgrade-Skripte ausführen zu können.
Das Ignorieren des Dogu-Health-Status kann daher zu Folgefehlern während der Ausführung des Blueprints führen.