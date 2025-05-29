#!/bin/sh

echo "Starting SSL generation process..."

# Stop OpenResty (Nginx)
echo "Stopping OpenResty..."
/usr/local/openresty/nginx/sbin/nginx -s stop

# Backup existing Nginx configuration
echo "Backing up existing Nginx configuration..."
mv /usr/local/openresty/nginx/conf/nginx.conf /usr/local/openresty/nginx/conf/nginx.conf.bak

# Define the Certbot webroot directory
CERTBOT_WEBROOT="/var/www/certbot"

# Ensure the Certbot webroot directory exists
if [ ! -d "$CERTBOT_WEBROOT/.well-known/acme-challenge" ]; then
    echo "Creating Certbot webroot directory: $CERTBOT_WEBROOT"
    mkdir -p "$CERTBOT_WEBROOT/.well-known/acme-challenge"
fi

# Create a temporary Nginx configuration for Certbot challenge
echo "Creating temporary Nginx configuration for Certbot validation..."
cat <<EOL > /usr/local/openresty/nginx/conf/nginx.conf
worker_processes 1;
events { worker_connections 1024; }
http {
    server {
        listen 80;
        server_name squidweb.ir www.squidweb.ir;

        location /.well-known/acme-challenge/ {
            root $CERTBOT_WEBROOT;
        }

        location / {
            return 404;
        }
    }
}
EOL

# Test and start Nginx with temporary configuration
echo "Testing temporary Nginx configuration..."
/usr/local/openresty/nginx/sbin/nginx -t

echo "Starting Nginx for Certbot validation..."
/usr/local/openresty/nginx/sbin/nginx

# Generate certificates using Certbot
echo "Generating certificates with Certbot..."
certbot certonly --webroot -w "$CERTBOT_WEBROOT" \
    -d squidweb.ir -d www.squidweb.ir \
    --agree-tos --non-interactive --email amirex128@gmail.com

if [ $? -ne 0 ]; then
    echo "Certbot failed! Generating self-signed SSL certificate..."
    SSL_DIR="/etc/letsencrypt/live/squidweb.ir"
    mkdir -p $SSL_DIR
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout $SSL_DIR/privkey.pem -out $SSL_DIR/fullchain.pem \
        -subj "/C=US/ST=State/L=City/O=Company/OU=Org/CN=squidweb.ir"

    echo "Updating Nginx configuration for self-signed SSL..."
    cat <<EOL > /usr/local/openresty/nginx/conf/nginx.conf
worker_processes 1;
events { worker_connections 1024; }
http {
    server {
        listen 443 ssl;
        server_name squidweb.ir www.squidweb.ir;
        ssl_certificate /etc/letsencrypt/live/squidweb.ir/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/squidweb.ir/privkey.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;
        location / {
            root /var/www/html;
            index index.html;
        }
    }
}
EOL
fi

# Stop temporary Nginx
echo "Stopping temporary Nginx..."
/usr/local/openresty/nginx/sbin/nginx -s stop

# Restore original Nginx configuration if Certbot succeeded
if [ -f "/etc/letsencrypt/live/squidweb.ir/fullchain.pem" ]; then
    echo "Restoring original Nginx configuration..."
    mv /usr/local/openresty/nginx/conf/nginx.conf.bak /usr/local/openresty/nginx/conf/nginx.conf
fi

# Restart OpenResty
echo "Restarting OpenResty with SSL configuration..."
/usr/local/openresty/nginx/sbin/nginx -g 'daemon off;'
