# SSL Certificates for Rojgarsetu

For development, generate self-signed certificates:

On Linux/Mac:
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout rojgarsetu.key -out rojgarsetu.crt \
  -subj "/C=IN/ST=Delhi/L=Delhi/O=Rojgarsetu/OU=IT/CN=localhost"
```

On Windows (with OpenSSL installed):
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout rojgarsetu.key -out rojgarsetu.crt -subj "/C=IN/ST=Delhi/L=Delhi/O=Rojgarsetu/OU=IT/CN=localhost"
```

For production, use Let's Encrypt or a commercial CA certificate.