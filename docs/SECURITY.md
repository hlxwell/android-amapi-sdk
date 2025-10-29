# 安全指南

## 敏感信息保护

本项目使用 Google Cloud 服务账号进行身份验证。请务必保护好您的凭证信息。

### ⚠️ 不要提交的文件

以下文件包含敏感信息，**切勿提交到 Git 仓库**：

- `sa-key.json` - 服务账号密钥文件
- `*-sa-key.json` - 任何服务账号密钥文件
- `service-account*.json` - 服务账号相关文件
- `amapi.yaml` / `amapi.yml` / `amapi.json` - 配置文件（包含项目 ID）
- `.env` - 环境变量文件

### ✅ 安全配置步骤

#### 1. 创建服务账号密钥

```bash
# 在 Google Cloud Console 创建服务账号
# 下载密钥文件到项目目录
cp ~/Downloads/your-key.json ./sa-key.json
```

#### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填写实际值
vi .env
```

#### 3. 配置 YAML 文件

```bash
# 复制配置文件模板
cp config.yaml.example config.yaml

# 编辑配置文件
vi config.yaml
```

### 🔒 已保护的文件

`.gitignore` 文件已配置忽略以下敏感文件：

```gitignore
# 密钥文件
sa-key.json
*-sa-key.json
service-account*.json

# 配置文件
amapi.yaml
amapi.yml
amapi.json
.env
.env.local
```

### 📝 使用示例文件

项目提供了以下示例文件，可以安全地提交到版本控制：

- `sa-key.json.example` - 服务账号密钥模板
- `.env.example` - 环境变量模板
- `config.yaml.example` - YAML 配置文件模板

### 🚨 如果意外提交了敏感信息

如果您不小心提交了包含敏感信息的文件：

1. **立即撤销服务账号密钥**
   ```bash
   gcloud iam service-accounts keys delete KEY_ID \
     --iam-account=SERVICE_ACCOUNT_EMAIL
   ```

2. **创建新的密钥**
   ```bash
   gcloud iam service-accounts keys create ./sa-key.json \
     --iam-account=SERVICE_ACCOUNT_EMAIL
   ```

3. **从 Git 历史中删除敏感文件**
   ```bash
   # 使用 git filter-branch 或 BFG Repo-Cleaner
   git filter-branch --force --index-filter \
     "git rm --cached --ignore-unmatch sa-key.json" \
     --prune-empty --tag-name-filter cat -- --all

   # 强制推送（谨慎使用）
   git push origin --force --all
   ```

4. **通知团队成员**更新他们的本地仓库

### 🔐 最佳实践

1. **使用环境变量**：生产环境优先使用环境变量而非配置文件
2. **最小权限原则**：服务账号仅授予必要的权限
3. **定期轮换密钥**：建议每 90 天轮换一次服务账号密钥
4. **使用 Secret Manager**：生产环境考虑使用 Google Secret Manager
5. **审计日志**：定期检查 Cloud Audit Logs 以发现异常访问

### 📚 相关资源

- [项目主文档](../README.md)
- [快速开始](QUICKSTART.md)
- [构建指南](BUILD_GUIDE.md)
- [Google Cloud 服务账号最佳实践](https://cloud.google.com/iam/docs/best-practices-service-accounts)
- [管理服务账号密钥](https://cloud.google.com/iam/docs/creating-managing-service-account-keys)
- [Android Management API 安全性](https://developers.google.com/android/management/security)

### 🆘 报告安全问题

如果发现安全漏洞，请通过以下方式报告：

- 创建私有 Security Advisory（推荐）
- 发送邮件到项目维护者

**请勿在公开 Issue 中讨论安全漏洞。**

