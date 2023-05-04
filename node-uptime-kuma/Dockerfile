FROM node:14-alpine

WORKDIR /app

RUN apk update && \
    apk upgrade && \
    apk add --no-cache \
      git \
      make \
      g++ \
      python \
      libc6-compat \
      curl

RUN git clone https://github.com/louislam/uptime-kuma.git .

RUN yarn install --production=true

EXPOSE 3001

CMD ["yarn", "start"]
