# Use the official OpenResty base image
FROM openresty/openresty:alpine-fat

# Install required dependencies: Certbot, OpenSSL, and Certbot plugins
RUN apk add --no-cache openssl certbot certbot-nginx

# Install lua-resty-jwt using luarocks
RUN /usr/local/openresty/luajit/bin/luarocks install lua-resty-jwt

# Ensure the Certbot webroot directory exists and has proper permissions
RUN mkdir -p /var/www/certbot/.well-known/acme-challenge

# Create necessary directories
RUN mkdir -p /usr/local/openresty/nginx/conf/config

# Copy Nginx configuration and entrypoint script to their respective locations
COPY stage/nginx.conf /usr/local/openresty/nginx/conf/
COPY stage/config /usr/local/openresty/nginx/conf/config/
COPY stage/servers /usr/local/openresty/nginx/conf/servers/

COPY stage/config/entrypoint.sh /entrypoint.sh

# Make entrypoint script executable
RUN chmod +x /entrypoint.sh

# Expose HTTP and HTTPS ports
EXPOSE 80 443

# Use entrypoint script to generate certificates and start Nginx
ENTRYPOINT ["/entrypoint.sh"]
