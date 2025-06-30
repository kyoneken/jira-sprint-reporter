# Jira Sprint Reporter

Jiraのアクティブスプリントからストーリーやタスクを取得し、報告用のリスト形式で出力するCLIツールです。Confluence等への貼り付けに最適化されたタブ区切り形式で出力します。

## セットアップ

1. 環境変数ファイルを作成:
```bash
cp .env.example .env
```

2. `.env`ファイルを編集して以下の情報を設定:
```
JIRA_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-api-token
```

### Jira API トークンの取得方法

1. Atlassianアカウントにログイン
2. アカウント設定 → セキュリティ → API トークンを作成
3. 生成されたトークンを`JIRA_API_TOKEN`に設定

#### 重要：スコープなしのAPIトークンが必要

**このアプリケーションは、スコープ制限のないAPIトークンが必要です。**

APIトークン作成時：
- **「スコープを選択しない」** - 制限なしのトークンを作成してください
- 特定のスコープ（`read:jira-work`、`read:jira-user`など）を選択すると、Agile API (`/rest/agile/1.0/`) へのアクセスができません

**なぜスコープなしが必要か:**
- Jira の Agile API（スプリント、ボード情報）にアクセスするため
- スコープ付きトークンでは401 Unauthorizedエラーが発生します
- 読み取り専用での利用なので、セキュリティ上の問題はありません

**トークン作成手順:**
1. Jira Cloud → アカウント設定 → セキュリティ → API tokens
2. 「Create API token」をクリック
3. 名前を設定（例：「Jira Sprint Reporter」）
4. **スコープ選択画面では何も選択せず、デフォルトの制限なしトークンを作成**
5. 生成されたトークンをコピーして`.env`ファイルに設定

## 使い方

```bash
# ビルド
go build -o jira-sprint-reporter

# 実行
./jira-sprint-reporter
```

## 機能

- アクティブスプリントの取得
- インタラクティブなスプリント選択
- 報告用リスト形式での出力（タブ区切り）
- Confluenceテーブルへの貼り付け対応
- チケット詳細情報（キー、リンク、タイプ、概要、エピック、ストーリーポイント、ステータス、担当者）

## 必要な権限

- Jiraプロジェクトへの読み取りアクセス
- アジャイルボードへのアクセス