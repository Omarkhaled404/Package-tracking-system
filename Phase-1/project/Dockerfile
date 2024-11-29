# # Build stage
# FROM node:alpine AS build-stage

# WORKDIR /app
# COPY . .

# # Install dependencies
# RUN npm install

# # Install Angular CLI globally
# RUN npm install -g @angular/cli

# # Build the Angular app in production mode
# RUN ng build --configuration=production

# # Production stage
# FROM nginx:alpine AS production-stage
# COPY --from=build-stage /app/dist/project /usr/share/nginx/html
# EXPOSE 80
# CMD ["nginx", "-g", "daemon off;"]

# Stage 1: Build the Angular application
FROM node:18 as build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the Angular project files
COPY . .

# Build the Angular project
RUN npm run build --prod

# Stage 2: Serve the application using Nginx
FROM nginx:1.25

# Copy the built Angular files (client-side browser files) from the previous stage
COPY --from=build /app/dist/project/browser /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf


# Expose the desired custom port
EXPOSE 8080

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]





