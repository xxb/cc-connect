# Contributor Fork Workflow

[中文](#中文) | [English](#english)

This document defines the preferred Git setup for contributors who occasionally work on **cc-connect** through a personal fork. It is written for both humans and AI coding assistants so a freshly cloned checkout can be normalized consistently across machines.

## English

### Summary

This is the standard workflow for contributing to someone else's open-source repo:

1. fork the upstream repo
2. clone **your fork**
3. add the original repo as `upstream`
4. keep local `main` aligned with `upstream/main`
5. do all work on feature branches
6. push branches to your fork
7. open PRs from your fork branch to `upstream/main`

### Canonical Remote Layout

- `origin` = your fork (`xxb/cc-connect`), writable
- `upstream` = canonical repo (`chenhg5/cc-connect`), fetch-only
- `upstream` push URL should be `DISABLED`

### Preferred Fresh-Machine Setup

Use GitHub CLI if available:

```bash
gh repo clone xxb/cc-connect
cd cc-connect
git remote add upstream git@github.com:chenhg5/cc-connect.git
git remote set-url --push upstream DISABLED
git branch --set-upstream-to=upstream/main main
git config branch.main.pushRemote origin
git config remote.pushDefault origin
git config fetch.prune true
git config pull.ff only
git config rebase.autoStash true
git config push.autoSetupRemote true
git fetch --all --prune
git remote -v
git branch -vv
```

If that machine uses HTTPS instead of SSH, replace the remote URLs accordingly.

### Normalize An Existing Checkout

#### Case A: legacy layout where `origin` = upstream and `fork` = personal fork

```bash
git switch main
git remote rename origin upstream
git remote rename fork origin
git branch --set-upstream-to=upstream/main main
git config branch.main.pushRemote origin
git config remote.pushDefault origin
git remote set-url --push upstream DISABLED
git config fetch.prune true
git config pull.ff only
git config rebase.autoStash true
git config push.autoSetupRemote true
git fetch --all --prune
git reset --hard upstream/main
```

#### Case B: cloned directly from upstream and no personal fork remote exists yet

```bash
git switch main
git remote rename origin upstream
git remote add origin git@github.com:xxb/cc-connect.git
git branch --set-upstream-to=upstream/main main
git config branch.main.pushRemote origin
git config remote.pushDefault origin
git remote set-url --push upstream DISABLED
git config fetch.prune true
git config pull.ff only
git config rebase.autoStash true
git config push.autoSetupRemote true
git fetch --all --prune
git reset --hard upstream/main
```

### Branch Rules

- `main` is a sync branch, not a development branch
- local `main` should track `upstream/main`
- `origin/main` should usually mirror `upstream/main`
- every PR should come from a dedicated feature branch
- push feature branches to `origin/<branch>`
- open PRs from `origin/<branch>` to `upstream/main`

### Daily Commands

#### Sync `main`

```bash
git switch main
git fetch upstream
git reset --hard upstream/main
```

#### Start a new branch

```bash
git switch -c feat/my-change
git push -u origin feat/my-change
```

#### Continue working on an existing PR branch from another machine

```bash
git fetch origin
git switch -c spike/discord-progress-card --track origin/spike/discord-progress-card
```

### Rebase Guidance

If a PR is already open and only you use the branch, rebasing is technically safe, but do **not** rebase by default when:

- the PR is already under review
- there is no merge conflict
- CI is green
- the maintainer did not ask for a rebase

Rebase only when there is a clear benefit, such as conflict resolution, maintainer request, or a required refresh against latest `main`.

If you rebase a branch that was already pushed, update the remote branch with:

```bash
git push --force-with-lease
```

### Guidance For AI Assistants

Unless the user explicitly asks for another workflow:

- prefer `gh repo clone xxb/cc-connect` for a fresh checkout
- ensure `origin` is the user's fork and `upstream` is `chenhg5/cc-connect`
- configure local `main` to track `upstream/main`
- set `upstream` push URL to `DISABLED`
- avoid feature work on `main`
- ask before destructive operations such as `reset --hard`, branch deletion, or force-push

---

## 中文

### 总结

这就是参与别人维护的开源仓库时最常见、最标准的流程：

1. fork 上游仓库
2. clone **你自己的 fork**
3. 把原仓库加为 `upstream`
4. 让本地 `main` 始终对齐 `upstream/main`
5. 所有开发都在功能分支上完成
6. 把分支推到你自己的 fork
7. 从你的 fork 分支向 `upstream/main` 提 PR

### 统一远端布局

- `origin` = 你的 fork（`xxb/cc-connect`），可写
- `upstream` = 上游仓库（`chenhg5/cc-connect`），只拉不推
- `upstream` 的 push URL 应设置成 `DISABLED`

### 新机器上的推荐初始化步骤

如果机器上装了 GitHub CLI，优先这样做：

```bash
gh repo clone xxb/cc-connect
cd cc-connect
git remote add upstream git@github.com:chenhg5/cc-connect.git
git remote set-url --push upstream DISABLED
git branch --set-upstream-to=upstream/main main
git config branch.main.pushRemote origin
git config remote.pushDefault origin
git config fetch.prune true
git config pull.ff only
git config rebase.autoStash true
git config push.autoSetupRemote true
git fetch --all --prune
git remote -v
git branch -vv
```

如果这台机器使用 HTTPS 而不是 SSH，把远端地址替换成 HTTPS 即可。

### 规范化已有仓库

#### 情况 A：历史布局是 `origin` = 上游、`fork` = 个人 fork

```bash
git switch main
git remote rename origin upstream
git remote rename fork origin
git branch --set-upstream-to=upstream/main main
git config branch.main.pushRemote origin
git config remote.pushDefault origin
git remote set-url --push upstream DISABLED
git config fetch.prune true
git config pull.ff only
git config rebase.autoStash true
git config push.autoSetupRemote true
git fetch --all --prune
git reset --hard upstream/main
```

#### 情况 B：一开始直接 clone 了上游，还没有个人 fork 远端

```bash
git switch main
git remote rename origin upstream
git remote add origin git@github.com:xxb/cc-connect.git
git branch --set-upstream-to=upstream/main main
git config branch.main.pushRemote origin
git config remote.pushDefault origin
git remote set-url --push upstream DISABLED
git config fetch.prune true
git config pull.ff only
git config rebase.autoStash true
git config push.autoSetupRemote true
git fetch --all --prune
git reset --hard upstream/main
```

### 分支规则

- `main` 只是同步分支，不直接开发
- 本地 `main` 应跟踪 `upstream/main`
- `origin/main` 通常也应尽量与 `upstream/main` 保持一致
- 每个 PR 都应来自单独的功能分支
- 功能分支推到 `origin/<branch>`
- PR 从 `origin/<branch>` 提到 `upstream/main`

### 日常命令

#### 同步 `main`

```bash
git switch main
git fetch upstream
git reset --hard upstream/main
```

#### 新建开发分支

```bash
git switch -c feat/my-change
git push -u origin feat/my-change
```

#### 在另一台机器继续已有 PR 分支

```bash
git fetch origin
git switch -c spike/discord-progress-card --track origin/spike/discord-progress-card
```

### 关于 rebase

如果某个 PR 已经打开，而且这条分支只有你自己在使用，那么 rebase 在技术上通常是安全的；但在以下情况中，**默认不要为了“整理历史”而随便 rebase**：

- PR 已经在 review 中
- 当前没有冲突
- CI 正常
- maintainer 没有要求你 rebase

只有在确实有收益时再做，比如：

- 需要解决冲突
- maintainer 明确要求更新分支
- 必须同步到最新 `main`

如果该分支之前已经 push 过，rebase 后要用：

```bash
git push --force-with-lease
```

### 给 AI 助手的规则

除非用户明确要求别的 Git 工作流，否则应默认：

- 新机器优先使用 `gh repo clone xxb/cc-connect`
- 确保 `origin` 是用户自己的 fork，`upstream` 是 `chenhg5/cc-connect`
- 确保本地 `main` 跟踪 `upstream/main`
- 把 `upstream` 的 push URL 设为 `DISABLED`
- 不要直接在 `main` 上做功能开发
- 对 `reset --hard`、删分支、force-push 这类破坏性操作，先说明影响再征求确认
