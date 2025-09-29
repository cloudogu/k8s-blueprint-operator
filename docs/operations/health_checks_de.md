# Health-Checks

Vor und nach dem Anwenden des Blueprints wird gewartet, dass das Ecosystem healthy ist.
Dabei wird folgendes geprüft:
- Health aller Dogus anhand der Dogu-CRs
- Überprüfung, ob alle Dogus bereits die neueste Version und Konfiguration verwenden

## Health ignorieren

Die Health-Checks vor der Ausführung des Blueprints können deaktiviert werden:
- für Dogus, wenn `spec.ignoreDoguHealth` auf `true` gesetzt wird,

So ist es möglich, per Blueprint Fehler an Dogus und Komponenten zu beheben.
Für ein Dogu-Upgrade muss ein Dogu allerdings healthy sein, um Pre-Upgrade-Skripte ausführen zu können.
Das Ignorieren der Dogu-Health kann also zu Folgefehlern während der Ausführung des Blueprints führen.
