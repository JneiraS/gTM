# 1. Utiliser une image officielle Go comme base
FROM golang:1.23.3-alpine

# Installer les dépendances nécessaires
RUN apk add --no-cache gcc musl-dev sqlite-dev


# 2. Définir le répertoire de travail dans le conteneur
WORKDIR /app

# 3. Copier les fichiers Go dans le conteneur
COPY . .

# 4. Télécharger les dépendances
RUN go mod tidy

# 5. Compiler l'application en un binaire statique
# RUN go run src/main.go

# 6. Exposer le port utilisé par l'application
EXPOSE 7263

# 7. Définir la commande pour exécuter l'application
CMD ["go", "run", "./src/main.go"]
