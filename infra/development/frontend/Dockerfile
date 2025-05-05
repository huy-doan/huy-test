FROM node:20.17.0-alpine

WORKDIR /app

ENV NODE_ENV=development

# Copy package.json and yarn.lock first for caching
COPY ["package.json", "yarn.lock", "./"]

# Install dependencies (including devDependencies)
RUN yarn install --frozen-lockfile && yarn cache clean
ARG ENV_TYPE=develop
ENV ENV_TYPE=${ENV_TYPE}

# Copy source files
COPY . ./

# Expose the port
EXPOSE 3000

CMD ["yarn", "dev"]
