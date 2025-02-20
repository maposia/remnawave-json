### 📖 V2Ray Subscription Proxy with User-Agent Routing

This application serves as a **proxy** for generating and serving **V2Ray subscription configurations** based on `User-Agent`. It dynamically provides the correct configuration format depending on the client version, ensuring seamless compatibility across different applications.

Work with https://remna.st

---

## ✨ Features
- **🔀 User-Agent-Based Routing**
   - Automatically detects and serves the correct subscription format for supported clients:
      - **v2rayN** (`>=6.40` JSON, older versions Base64)
      - **v2rayNG** (`>=1.8.29` JSON, older versions Base64)
      - **Streisand** (JSON)
      - **Happ** (`>=1.63.1` JSON, older versions Base64)
- **🛠 Mux support**
   - Supported `mux` template.
- **🌍 Direct Proxy Fallback**
   - If `User-Agent` is unsupported or the request doesn’t match `/v2ray-json`, the server provides a **default proxy response**.

---

## ⚙️ Configuration
Modify `.env.sample` to adjust the application settings:
```
REMNAWAWE_URL=sub_domain
APP_PORT=4000
# V2RAY_TEMPLATE_PATH=/app/templates/v2ray/default.json
# V2RAY_MUX_ENABLED=true
# V2RAY_MUX_TEMPLATE_PATH=/app/templates/v2ray/mux_default.json
# WEB_PAGE_TEMPLATE_PATH=/app/templates/subscription/index.html
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

## Environment Variables

1. **REMNAWAWE_URL**  
   _Description:_ The base URL for sub domain. [Installation Environment Variables](https://remna.st/installation/env#subscription-public-domain)  
   _Example:_ `REMNAWAWE_URL=domain`

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
   _Example:_ `V2RAY_MUX_TEMPLATE_PATH=/app/templates/subscription/index.html`


## 📜 License
This project is open-source and available under the **MIT License**.

🚀 *Enhance your V2Ray subscription management with automated User-Agent routing!*

