# MiniGit

MiniGit est une version simplifiée de Git, développée en Go, permettant de suivre les modifications de fichiers, effectuer des commits, gérer des branches, etc.

---

## 🔧 Compilation

```bash
go build -o minigit
```

Cela génère un exécutable nommé `minigit`.

Place le binaire à la racine de ton projet, ou adapte les chemins lors des appels si tu le déplaces.

---

## Commandes

### 1. `init`

Initialise un dépôt MiniGit (création du dossier `.miniGit`).

```bash
./minigit init
```

---

### 2. `set-user`

Modifie le nom d'utilisateur pour les futurs commits.

```bash
./minigit set-user <username>
```

---

### 3. `get-user`

Affiche le nom d'utilisateur actuellement configuré.

```bash
./minigit get-user
```

---

### 4. `add`

Ajoute un ou plusieurs fichiers à l'index (staging area).

```bash
./minigit add <file1> <file2>
./minigit add .   # ajoute tous les fichiers du dossier courant
```

---

### 5. `commit`

Crée un commit avec les fichiers présents dans l’index.

```bash
./minigit commit -m "Message de commit"
```

---

### 6. `status`

Affiche l’état du dépôt :
- fichiers modifiés et prêts à être commités,
- fichiers modifiés mais non indexés,
- fichiers non suivis.

```bash
./minigit status
```

---

### 7. `diff`

Affiche les différences entre les fichiers présents dans le disque et ceux dans l’index.

```bash
./minigit diff
```

---

### 8. `log`

Affiche l'historique des commits de la branche actuelle.

```bash
./minigit log
```

---

### 9. `reset`

Réinitialise l’index (retire tous les fichiers du staging area, sans modifier les fichiers du disque).

```bash
./minigit reset
```

---

### 10. `branch`

Crée une nouvelle branche ou affiche les branches existantes.

```bash
./minigit branch           # liste les branches
./minigit branch <name>    # crée une nouvelle branche
```

---

### 11. `checkout`

Change de branche.

```bash
./minigit checkout <branch_name>
```

---

### 12.0 `merge`

Fusionne une branche dans la branche actuelle.

```bash
./minigit merge <branch_name>
```

---

### 12.1 `merge --abort`

Restaure tous les fichiers du disque à l’état exact du dernier commit sur la branche courante (avant le merge).

```bash
./minigit merge --abort
```

---

### 13. `revert`

Annule les changements d’un commit donné (en créant un commit inverse).

```bash
./minigit revert <hash_commit>
```

---

### 14. `rm`

Supprime un fichier du disque et/ou de l’index.

- Supprime du disque **et** du prochain commit (sauf s’il est modifié en staging) :
```bash
./minigit rm <fichier>
```

- Supprime uniquement de l’index (s'il est modifié en staging) :
```bash
./minigit rm --cached <fichier>
```

- Supprime de force, même s’il est modifié :
```bash
./minigit rm -f <fichier>
```

---

## 📌 Notes importantes

- Après un `rm`, il faut **faire un commit** pour enregistrer la suppression dans l'historique.
- Les fichiers marqués avec un hash vide (`""`) dans l’index seront exclus du prochain commit.