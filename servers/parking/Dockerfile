FROM node
# set the current working directory in tcontainer
WORKDIR /parking
COPY package.json package-lock.json ./
RUN npm install
COPY . .
EXPOSE 80
ENTRYPOINT ["node", "index.js"]