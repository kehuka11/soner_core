# 認知負荷分析ロジック

このドキュメントでは、SonerCore におけるコードの認知負荷を測定するために使用されるロジックについて説明します。

## 概要
認知負荷を推定するための主要な指標として **Cognitive Complexity（認知的複雑度）** を使用します。この指標は、循環的複雑度（Cyclomatic Complexity）の限界に対処するために SonarSource によって開発されました。

## 参考文献
- [Cognitive Complexity: A new way of measuring understandability](https://www.sonarsource.com/docs/CognitiveComplexity.pdf) (SonarSource)

## アルゴリズムの説明

### 基本原則
Cognitive Complexity は、コードブロックの制御フローを理解するのがどれほど難しいかを評価します。以下の要素に基づいてスコアを加算します：
1.  **ネスト（入れ子）**: ネストが深くなるほど、コストが高くなります。
2.  **構造的切断**: 線形フローを中断するコード構造（例：`if`, `for`, `while`, `catch`）。

### スコアリングルール（Go言語向け簡易版）

スコアは 0 から始まります。

1.  **加算**:
    *   `if`, `else if`, `else`
    *   `switch`
    *   `for`, `range`
    *   `select`
    *   `&&`, `||` のシーケンス
    *   `catch`（他の言語の場合。Go の `recover` ロジックも適用される可能性がありますが、標準的なエラー処理 `if err != nil` は通常免除されるか軽く扱われます。ただし、厳密な Cognitive Complexity では `if` としてカウントされます）
    *   再帰

2.  **ネストペナルティ**:
    *   構造がネストレベルを上げるたびに、加算値は `1 + ネストレベル` となります。

3.  **免除**:
    *   複雑さを加えない短縮構造（例：一部の言語における単純な null チェック。ただし Go は明示的です）。

### 実装の詳細
Go AST（抽象構文木）を使用してコードを走査します。

#### AST 走査
- **Visitor パターン**: AST をウォークスルーします。
- **状態**: `nesting`（ネスト）カウンターを維持します。
- **スコアリング**:
    - `*ast.IfStmt` 訪問時: +1 + ネスト
    - `*ast.ForStmt`, `*ast.RangeStmt` 訪問時: +1 + ネスト
    - `*ast.SwitchStmt`, `*ast.TypeSwitchStmt` 訪問時: +1 + ネスト
    - `*ast.SelectStmt` 訪問時: +1 + ネスト
    - `*ast.BinaryExpr`（`LAND` または `LOR`）訪問時: +1（通常ネストペナルティはありませんが、シーケンスはカウントされます）

### 今後の改善点
- "Halstead Complexity Measures"（ハルステッド複雑度）の追加
- "Maintainability Index"（保守性指数）の追加
