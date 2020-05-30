FROM node:12-alpine

LABEL maintainer="Yaroslav Muravskyi <mail@myarik.com>"

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
COPY package*.json ./
RUN npm install --only=production

# Bundle app source
COPY . .

CMD [ "node", "first_response.js" ]