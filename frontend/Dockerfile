FROM node:20-alpine AS build
WORKDIR /app

COPY package*.json ./
RUN npm install --registry=https://registry.npmjs.org/

COPY . .
ARG VITE_API_BASE=http://localhost:8080
ENV VITE_API_BASE=${VITE_API_BASE}
RUN npm run build

FROM nginx:alpine
WORKDIR /usr/share/nginx/html
COPY --from=build /app/dist ./
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

