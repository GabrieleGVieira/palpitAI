# PalpitAI

O bolao inteligente: app mobile de Copa do Mundo com grupos privados, rankings, palpites em tempo real e insights de IA.

## Stack

| Camada | Tecnologias |
| ------ | ----------- |
| Mobile | React Native, Expo, TypeScript |
| Qualidade | ESLint, Prettier, Husky, Commitlint |

## Estrutura

```text
palpitAI/
├── frontend/          # App Expo React Native
│   ├── App.tsx
│   ├── index.ts
│   ├── app.json
│   ├── eslint.config.js
│   ├── commitlint.config.js
│   └── package.json
├── codealike.json
└── README.md
```

## Requisitos

- Node.js
- npm
- Expo Go no dispositivo fisico, Android Emulator ou iOS Simulator

## Como rodar

Instale as dependencias:

```bash
cd frontend
npm install
```

Inicie o Metro Bundler:

```bash
npm run start
```

Ou rode diretamente em uma plataforma:

```bash
npm run android
npm run ios
npm run web
```

## Qualidade de codigo

```bash
npm run lint          # Executa ESLint
npm run lint:fix      # Corrige problemas automaticos do ESLint
npm run format        # Formata com Prettier
npm run format:check  # Verifica formatacao
npm run typecheck     # Verifica TypeScript
```

## Commits

O projeto usa Husky e Commitlint para validar commits no padrao Conventional Commits.

Exemplos:

```bash
feat: add login screen
fix: correct prediction score
chore: update dependencies
```

Hooks configurados:

- `pre-commit`: roda `npm run lint` e `npm run typecheck`
- `commit-msg`: valida a mensagem com Commitlint

## Scripts do frontend

```bash
npm run start
npm run android
npm run ios
npm run web
npm run lint
npm run lint:fix
npm run format
npm run format:check
npm run typecheck
```
