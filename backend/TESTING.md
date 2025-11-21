# Guide de Test - SaveSync Backend

## Démarrage du Serveur

```bash
# Compiler
go build -o savesyncd ./cmd/savesyncd

# Lancer
./savesyncd

# Ou directement avec Go
go run ./cmd/savesyncd
```

Le serveur démarre sur `http://localhost:8080`

## Tests des Endpoints

### 1. Health Check

```bash
curl http://localhost:8080/health
```

**Réponse attendue:**·
```json
{"data":{"status":"ok"}}
```

### 2. Métriques Prometheus

```bash
curl http://localhost:8080/metrics
```

Affiche toutes les métriques au format Prometheus.

---

## CRUD Sources

### Lister les sources

```bash
curl http://localhost:8080/api/sources
```

### Créer une source

```bash
curl -X POST http://localhost:8080/api/sources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mes-documents",
    "path": "/home/user/documents",
    "exclusions": ["*.tmp", "*.log", "node_modules", ".git"],
    "target_id": 1
  }'
```

**Note:** Le path doit exister sur le système.

### Récupérer une source

```bash
curl http://localhost:8080/api/sources/1
```

### Mettre à jour une source

```bash
curl -X PUT http://localhost:8080/api/sources/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "mes-documents-updated",
    "path": "/home/user/documents",
    "exclusions": ["*.tmp", "*.log"],
    "target_id": 1
  }'
```

### Supprimer une source

```bash
curl -X DELETE http://localhost:8080/api/sources/1
```

---

## CRUD Targets

### Lister les targets

```bash
curl http://localhost:8080/api/targets
```

### Créer un target LOCAL

```bash
curl -X POST http://localhost:8080/api/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-local",
    "type": "local",
    "config": {
      "path": "/tmp/backups"
    }
  }'
```

### Créer un target S3

```bash
curl -X POST http://localhost:8080/api/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-s3",
    "type": "s3",
    "config": {
      "bucket": "my-backups",
      "region": "us-east-1",
      "access_key": "YOUR_ACCESS_KEY",
      "secret_key": "YOUR_SECRET_KEY",
      "endpoint": "https://s3.amazonaws.com"
    }
  }'
```

### Créer un target SFTP

```bash
curl -X POST http://localhost:8080/api/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-sftp",
    "type": "sftp",
    "config": {
      "host": "backup.example.com",
      "port": "22",
      "user": "backup-user",
      "password": "your-password",
      "path": "/backups"
    }
  }'
```

**Avec clé SSH:**
```bash
curl -X POST http://localhost:8080/api/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-sftp-key",
    "type": "sftp",
    "config": {
      "host": "backup.example.com",
      "port": "22",
      "user": "backup-user",
      "key_path": "/home/user/.ssh/id_rsa",
      "path": "/backups"
    }
  }'
```

### Récupérer un target

```bash
curl http://localhost:8080/api/targets/1
```

### Mettre à jour un target

```bash
curl -X PUT http://localhost:8080/api/targets/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "backup-local-updated",
    "type": "local",
    "config": {
      "path": "/var/backups"
    }
  }'
```

### Supprimer un target

```bash
curl -X DELETE http://localhost:8080/api/targets/1
```

---

## Backups

### Déclencher un backup

```bash
curl -X POST http://localhost:8080/api/sources/1/run
```

**Réponse:**
```json
{
  "data": {
    "job_id": 1,
    "status": "pending"
  }
}
```

---

## Jobs

### Lister tous les jobs

```bash
curl http://localhost:8080/api/jobs
```

### Récupérer un job spécifique

```bash
curl http://localhost:8080/api/jobs/1
```

**Réponse:**
```json
{
  "data": {
    "id": 1,
    "type": "backup",
    "source_id": 1,
    "status": "success",
    "started_at": "2025-01-21T17:00:00Z",
    "ended_at": "2025-01-21T17:05:30Z"
  }
}
```

---

## Scénario de Test Complet

### 1. Créer un target local

```bash
curl -X POST http://localhost:8080/api/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-backup",
    "type": "local",
    "config": {"path": "/tmp/savesync-test"}
  }'
```

**Réponse:** Notez l'ID (ex: 1)

### 2. Créer un répertoire de test

```bash
mkdir -p /tmp/test-source
echo "Hello World" > /tmp/test-source/file1.txt
echo "Test Data" > /tmp/test-source/file2.txt
```

### 3. Créer une source

```bash
curl -X POST http://localhost:8080/api/sources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-source",
    "path": "/tmp/test-source",
    "exclusions": ["*.tmp"],
    "target_id": 1
  }'
```

**Réponse:** Notez l'ID (ex: 1)

### 4. Lancer le backup

```bash
curl -X POST http://localhost:8080/api/sources/1/run
```

### 5. Vérifier le job

```bash
curl http://localhost:8080/api/jobs/1
```

### 6. Vérifier les fichiers créés

```bash
ls -la /tmp/savesync-test/
```

Vous devriez voir:
- `chunks/` - Répertoire contenant les chunks
- `manifests/` - Répertoire contenant les manifests JSON

---

## Tests avec jq (formatage JSON)

Si vous avez `jq` installé:

```bash
# Lister les sources avec formatage
curl -s http://localhost:8080/api/sources | jq

# Récupérer un job avec formatage
curl -s http://localhost:8080/api/jobs/1 | jq

# Extraire uniquement le status
curl -s http://localhost:8080/api/jobs/1 | jq '.data.status'
```

---

## Variables d'Environnement

```bash
# Changer le port
export SERVER_PORT=9000
./savesyncd

# Changer le niveau de log
export LOG_LEVEL=debug
./savesyncd

# Changer le chemin de la base de données
export DATABASE_PATH=/var/lib/savesync/db.sqlite
./savesyncd
```

---

## Tests Unitaires

```bash
# Tous les tests
go test ./...

# Avec verbosité
go test -v ./...

# Avec couverture
go test -cover ./...

# Test d'un package spécifique
go test -v ./internal/app/sourceservice/...
```

---

## Monitoring

### Vérifier les métriques Prometheus

```bash
# Toutes les métriques
curl -s http://localhost:8080/metrics

# Filtrer les métriques savesync
curl -s http://localhost:8080/metrics | grep savesync

# Métriques de backup
curl -s http://localhost:8080/metrics | grep backup
```

**Métriques principales:**
- `savesync_backup_last_run_timestamp` - Timestamp du dernier backup
- `savesync_backup_status` - Statut (1=succès, 0=échec)
- `savesync_backup_duration_seconds` - Durée du backup
- `savesync_bytes_transferred_total` - Octets transférés
- `savesync_error_count_total` - Nombre d'erreurs

---

## Logs

Les logs sont au format JSON structuré:

```bash
# Lancer avec logs debug
LOG_LEVEL=debug ./savesyncd

# Filtrer les logs avec jq
./savesyncd 2>&1 | jq 'select(.level == "error")'
```

---

## Troubleshooting

### Erreur: "invalid path: directory does not exist"

Le répertoire source n'existe pas. Créez-le:
```bash
mkdir -p /path/to/source
```

### Erreur: "backend initialization failed"

Vérifiez la configuration du target:
- Pour local: le path doit être accessible
- Pour S3: vérifiez les credentials et le bucket
- Pour SFTP: vérifiez host, user, password/key

### Port déjà utilisé

```bash
# Changer le port
export SERVER_PORT=9000
./savesyncd
```

### Base de données verrouillée

SQLite n'autorise qu'une connexion en écriture. Fermez les autres instances.

---

## Docker

### Build

```bash
docker build -t savesync-backend .
```

### Run

```bash
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v /path/to/backup:/backups \
  savesync-backend
```

### Test avec Docker

```bash
# Health check
curl http://localhost:8080/health

# Créer un target
docker exec -it <container-id> sh
# Puis utiliser curl depuis le container
```

---

## Exemples Avancés

### Backup avec exclusions multiples

```bash
curl -X POST http://localhost:8080/api/sources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "project-backup",
    "path": "/home/user/project",
    "exclusions": [
      "*.tmp",
      "*.log",
      "node_modules",
      ".git",
      "dist",
      "build",
      "__pycache__",
      "*.pyc"
    ],
    "target_id": 1
  }'
```

### MinIO (S3-compatible)

```bash
curl -X POST http://localhost:8080/api/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "minio-backup",
    "type": "s3",
    "config": {
      "bucket": "backups",
      "region": "us-east-1",
      "access_key": "minioadmin",
      "secret_key": "minioadmin",
      "endpoint": "http://localhost:9000"
    }
  }'
```

### Backblaze B2

```bash
curl -X POST http://localhost:8080/api/targets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "b2-backup",
    "type": "s3",
    "config": {
      "bucket": "my-bucket",
      "region": "us-west-002",
      "access_key": "YOUR_KEY_ID",
      "secret_key": "YOUR_APP_KEY",
      "endpoint": "https://s3.us-west-002.backblazeb2.com"
    }
  }'
```
