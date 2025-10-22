# Blueprint-Maske

Die `blueprintMask` bietet eine leistungsstarke Möglichkeit, ein Blueprint für eine bestimmte Umgebung anzupassen, ohne die ursprüngliche Blueprint-Definition zu ändern. Sie fungiert als Filter, der es Ihnen ermöglicht, bestimmte Dogus, die im Hauptabschnitt `blueprint` definiert sind, selektiv zu deaktivieren.

## Anwendungsfall: Standardisierte Blueprints, angepasste Deployments

Stellen Sie sich vor, Sie pflegen ein einziges, umfassendes und bewährtes Blueprint, das eine vollständige Installation Ihrer Dogu-suite definiert. Dieses "Golden-Master"-Blueprint wird von mehreren Teams oder Kunden verwendet.

Allerdings benötigt nicht jedes Team jedes Dogu. Zum Beispiel benötigt ein Team möglicherweise das `redmine`-Dogu nicht.

Anstatt eine separate, leicht abweichende Blueprint-Datei für dieses Team zu erstellen und zu pflegen, können Sie die `blueprintMask` verwenden. Sie wenden dasselbe vollständige Blueprint auf jeden Cluster an, verwenden aber eine spezifische Maske für den Cluster dieses Teams, um die Installation von `redmine` zu verhindern.

Dieser Ansatz bietet mehrere Vorteile:
- **Single Source of Truth**: Sie verwalten ein Master-Blueprint, was die Komplexität und das Risiko von Konfigurationsabweichungen reduziert.
- **Konsistenz**: Alle Umgebungen basieren auf derselben getesteten Grundlage.
- **Flexibilität**: Sie können Dogus für jede gegebene Installation einfach und spontan aktivieren oder deaktivieren.

## Funktionsweise

Der `k8s-blueprint-operator` liest zuerst den `blueprint`-Abschnitt und wendet dann die `blueprintMask` darüber an. Wenn ein Dogu in der Maske mit `absent: true` aufgeführt ist, wird es aus der endgültigen Menge der zu installierenden oder zu verwaltenden Dogus entfernt. Das Ergebnis wird als "effektives Blueprint" bezeichnet.

## Beispiel

Betrachten Sie die folgende `Blueprint`-Ressource. Der `blueprint`-Abschnitt definiert sowohl `scm` als auch `redmine`.

```yaml
apiVersion: k8s.cloudogu.com/v2
kind: Blueprint
metadata:
  name: my-instance-blueprint
spec:
  # Dies ist das "Master"-Blueprint mit allen möglichen Dogus
  blueprint:
    dogus:
      - name: "official/scm"
        version: "3.11.0-1"
      - name: "official/redmine"
        version: "6.0.6-2"

  # Diese Maske passt das Blueprint für diese spezifische Instanz an
  blueprintMask:
    dogus:
      - name: "official/redmine"
        absent: true
```

### Ergebnis

Wenn der Operator diese Ressource verarbeitet:
1. Er sieht `scm` und `redmine` im `blueprint`.
2. Er wendet dann die `blueprintMask` an, die besagt, dass `redmine` abwesend sein soll.
3. Das resultierende **effektive Blueprint** enthält nur das `scm`-Dogu.

Infolgedessen installiert der Operator nur `official/scm:3.11.0-1` und ignoriert das `redmine`-Dogu vollständig. Wenn `redmine` bereits installiert war, würde es zur Deinstallation markiert.

## Syntax

Die Struktur der `blueprintMask` spiegelt das `blueprint` selbst wider. Um ein Dogu auszuschließen, müssen Sie nur seinen `name` und das Flag `absent: true` angeben.

```yaml
blueprintMask:
  dogus:
    - name: "<namespace>/<dogu-name>"
      absent: true
    # Fügen Sie hier weitere auszuschließende Dogus hinzu
```