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

► Crée ou écrase le fichier .miniGit/config contenant le nom d'utilisateur (ex: John Doe).

---

### 3. `get-user`

Affiche le nom d'utilisateur actuellement configuré.

```bash
./minigit get-user
```

► Lit simplement le fichier .miniGit/config

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

► Ce que la commande modifie :

- Crée un fichier temporaire .miniGit/MERGE_HEAD pour indiquer qu'une fusion est en cours

- Charge les trees de :

    - la branche courante (HEAD)

    - la branche à merger

    - leur ancêtre commun (via findCommonCommitAncestorHash)

    - Ces trees sont extraits depuis .miniGit/objects/commits/<hash> (ligne tree: ...), puis chargés via .miniGit/objects/trees/<tree_hash>

- Compare chaque fichier dans les trois trees et applique un merge à 3 voies :

    - Si HEAD et la branche sont identiques → le fichier n'est pas modifié.

    - Si les contenus divergent → un fichier fusionné est généré :

        - sans conflit : le fichier est directement mis à jour

        - avec conflit : les marqueurs <<<<<<<, =======, >>>>>>> sont insérés dans le fichier

- Si aucun conflit :

    - Les fichiers fusionnés sont sauvegardés sur le disque

    - Un blob est créé pour chaque fichier modifié dans .miniGit/objects/blobs/<hash>

    - Le fichier .miniGit/index est mis à jour avec les nouveaux hash (comme après un add)

    - Un nouveau tree est généré dans .miniGit/objects/trees/<hash>

    - Un commit de merge est créé dans .miniGit/objects/commits/<hash>, avec deux parents (HEAD et la branche fusionnée)

    - Le fichier .miniGit/MERGE_HEAD est supprimé à la fin

- Si conflit :

    - Aucun commit n'est créé automatiquement

    - Le fichier .miniGit/MERGE_HEAD reste présent

    - L’utilisateur doit résoudre manuellement les conflits, faire un add des fichiers corrigés, puis un commit pour terminer le merge

---

### 12.1 `merge --abort`

Restaure tous les fichiers du disque à l’état exact du dernier commit sur la branche courante (avant le merge).

```bash
./minigit merge --abort
```

► Ce que la commande modifie :

- Lit .miniGit/MERGE_HEAD pour retrouver le commit HEAD précédent

- Restaure tous les fichiers du disque à partir du tree de ce commit (.miniGit/objects/trees<tree_hash>)

- Vide le fichier .miniGit/index (reset complet de la staging area)

- Supprime .miniGit/MERGE_HEAD

---

### 13. `revert`

Annule les changements d’un commit donné (en créant un commit inverse).

```bash
./minigit revert <hash_commit>
```

► Ce que la commande modifie :

- Lit .miniGit/objects/commits/<hash_commit_ciblé> et son commit parent

- Récupère les fichiers .miniGit/objects/trees/<tree_hash> de ces deux commits

- Compare les arbres pour calculer le diff inverse

- Applique ce diff sur les fichiers du disque

- Met à jour le fichier .miniGit/index avec les nouveaux blobs (comme un add)

- Crée un nouveau commit dans .miniGit/objects/commits/ avec :

    - un nouveau tree dans .miniGit/objects/trees/

    - les blobs modifiés dans .miniGit/objects/blobs/ (si nécessaires)

    - le message Revert "<message d'origine>"

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

► Ce que la commande modifie :

- Supprime le fichier du disque (sauf si --cached)

- Met à jour .miniGit/index en supprimant l’entrée du fichier ou en la remplaçant par un hash vide ("")

- Le commit suivant enregistrera la suppression définitive dans un nouveau tree + commit


---

## 📌 Notes importantes

- Après un `rm`, il faut **faire un commit** pour enregistrer la suppression dans l'historique.
- Les fichiers marqués avec un hash vide (`""`) dans l’index seront exclus du prochain commit.