FROM node:20-alpine

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .

CMD ["npm", "start", "--", "--host", "0.0.0.0"]
EXPOSE 4200
