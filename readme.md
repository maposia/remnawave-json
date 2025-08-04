### 📖 V2Ray Json Subscription Proxy with User-Agent Routing

This application serves as a **proxy** for generating and serving **V2Ray subscription configurations** based on `User-Agent`. It dynamically provides the correct configuration format depending on the client version, ensuring seamless compatibility across different applications.

Work with <https://remna.st>

For use your own web page template, use {{.MetaTitle}} {{.MetaDescription}} {{.PanelData}} to get panel data and other.

## 🇷🇺 [Happ Routing](https://github.com/hydraponique/roscomvpn-happ-routing/tree/main)

## ✨ Features

- **🔀 User-Agent-Based Routing**
  - Automatically detects and serves the correct subscription format for supported clients:
    - **Streisand** (JSON)
    - **Happ** (JSON)
- **Web page template**
  - Supported web page template.
- **🌍 Direct Proxy Fallback**
  - If `User-Agent` is unsupported, the server provides a **default proxy response**.

---

## ⚙️ Configuration

Modify `.env.sample` to adjust the application settings:

```
REMNAWAVE_URL=sub_domain
APP_PORT=4000
# V2RAY_TEMPLATE_PATH=/app/templates/v2ray/default.json
# V2RAY_MUX_ENABLED=true
# V2RAY_MUX_TEMPLATE_PATH=/app/templates/v2ray/mux_default.json
# WEB_PAGE_TEMPLATE_PATH=/app/templates/subscription/index.html
META_TITLE=Zalupa
META_DESCRIPTION=Pupa
REMNAWAVE_TOKEN=
MODE=
```

After modifying execute this

```bash
mv .env.sample .env
```

---

# How to Run

1. Clone the repo

```bash
git clone https://github.com/Jolymmiles/remnawave-json
```

2. Go to the cloned repo

```bash
cd remnawave-json
```

3. Configure .env

4. Run Docker Compose

```bash
docker compose up -d
```

---

# How to update

1. Go to directory with docker-compose.yaml

```bash
remnawave-json
```

2. Update image

```bash
docker compose pull
```

3. Restart container

```bash
docker compose down --remove-orphans && docker compose up -d
```

---

## 🌿 **Environment Variables**

1. **REMNAWAVE_URL**  
   _Description:_ The base URL for the subdomain. [Installation Environment Variables](https://remna.st/installation/env#subscription-public-domain)  
   _Example:_ `REMNAWAVE_URL=https://panel.com`

2. **APP_PORT**  
   _Description:_ The port on which the application will run.  
   _Example:_ `APP_PORT=4000`

3. **V2RAY_TEMPLATE_PATH**  
   _Description:_ The file path to the default V2Ray configuration template.  
   _Example:_ `V2RAY_TEMPLATE_PATH=/app/templates/v2ray/default.json`

4. **V2RAY_MUX_ENABLED**  
   _Description:_ A flag to enable or disable the V2Ray Mux feature. Set to `true` to enable Mux.  
   _Example:_ `V2RAY_MUX_ENABLED=true`

5. **V2RAY_MUX_TEMPLATE_PATH**  
   _Description:_ The file path to the V2Ray Mux configuration template.  
   _Example:_ `V2RAY_MUX_TEMPLATE_PATH=/app/templates/v2ray/mux_default.json`

6. **WEB_PAGE_TEMPLATE_PATH**  
   _Description:_ The file path to the subscription template.  
   _Example:_ `WEB_PAGE_TEMPLATE_PATH=/app/templates/subscription/index.html`

7. **HAPP_JSON_ENABLED**  
   _Description:_ A flag to enable or disable JSON output for Happ.  
   _Example:_ `HAPP_JSON_ENABLED=false`

8. **HAPP_ROUTING**  
   _Description:_ The routing path for Happ connections.  
   _Example:_ `HAPP_ROUTING=happ://routing/...`

9. **HAPP_ANNOUNCEMENTS**  
   _Description:_ Announcement text. It can be provided in plain text.
   _Example:_ `HAPP_ANNOUNCEMENTS=zalupa`

10. **RU_OUTBOUND_NAME**  
    _Description:_ RU outbound name.
    _Example:_ `RU_OUTBOUND_NAME=RU`

11. **RU_USER_HOST**  
    _Description:_ RU user host.
    _Example:_ `RU_USER_HOST=Россия`

12. **REMNAWAVE_TOKEN**  
    _Description:_ REMNAWAVE token.
    _Example:_ `REMNAWAVE_TOKEN=zalupa`

13. **MetaDescription**  
    _Description:_ MetaDescription.
    _Example:_ `MetaDescription=Zalupa`

14. **MetaTitle**  
    _Description:_ MetaTitle.
    _Example:_ `MetaTitle=Zalupa`
15. **MODE****
    _Description:_ Set if use remnawave:3000.
    _Example:_ `MODE=local`

---

## Nginx example

```nginx configuration
server
{
        listen 443 ssl;
        listen [::]:443 ssl;
        http2 on;

        #from .env remnawave SUB_PUBLIC_DOMAIN
        server_name sub_domain ;


        # you certs
        ssl_certificate /root/.acme.sh/example.com_ecc/example.com.cer;
        ssl_certificate_key /root/.acme.sh/example.com_ecc/example.com.key;
        ssl_dhparam /etc/nginx/ssl/dhparams.pem;

        ssl_protocols TLSv1.3;
        ssl_ciphers TLS13-CHACHA20-POLY1305-SHA256:TLS13-AES-256-GCM-SHA384:TLS13-AES-128-GCM-SHA256:EECDH+CHACHA20:EECDH+AESGCM:EECDH+AES;
        ssl_prefer_server_ciphers on;

        # HSTS
        add_header Strict-Transport-Security "max-age=63072000" always;

        # OCSP stapling
        ssl_stapling on;
        ssl_stapling_verify on;
        ssl_trusted_certificate /root/.acme.sh/example.com_ecc/ca.cer;

        location / {
        proxy_http_version 1.1;
        # APP_PORT
        proxy_pass http://127.0.0.1:4000;
        proxy_set_header Host $host;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Port $server_port;

        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
 }
```

## 📜 License

This project is open-source and available under the **AGPL v2.0**.

🚀 _Enhance your V2Ray subscription management with automated User-Agent routing!_
