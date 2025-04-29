<p align="center">
  <img width="320" src="assets/logo-black.png#gh-light-mode-only" alt="HomeShare Logo">
  <img width="320" src="assets/logo-white.png#gh-dark-mode-only" alt="HomeShare Logo">
</p>
<h2 align="center">HomeShare - A File Server in your Home PC ⚡</h2>

<p align="center">
    <a href="/LICENSE"><img alt="GPL-V3.0 License" src="https://img.shields.io/badge/License-GPLv3-orange.svg"></a>
    <a href="https://github.com/jugeekuz/HomeShare/graphs/contributors"><img alt="Contributors" src="https://img.shields.io/github/contributors/jugeekuz/HomeShare?color=green"></a>
</p><br>

Turn your Home PC, into a globally accessible performant file server. With **HomeShare** you can push or pull files from anywhere to your computer, with rock‑solid authentication or create sharing folders with time‑limited share links to share with your friends.

<p align="center">
    <img src="./assets/demo-upload.gif" alt="HomeShare Short Demo" width="320">
</p>


--- 
## 📝 Description
HomeShare is a self-hosted file-sharing solution that transforms your home PC into a high-performance, globally accessible file server. With HomeShare you can skip the unnecessary step of the cloud and keep your data protected :
 - **Push & Pull Anywhere**: Upload or download files to your home machine from any location, using your own custom domain.
 - **Rock-Solid Authentication**: Protect access with user accounts and strong passwords.
 - **Time-Limited Share Links**: Create expiring URLs to share password-protected folders and individual files with friends, colleagues or clients.
 - **Automatic DNS & TLS**: Built-in DDNS updater let's your HomeShare app be globally accessible withouth the need for a Static IP.
 - **Easy Docker-Compose Deploy**: Get up and running in under 10 minutes with just Docker and a single `docker compose up` command.


## ⚡ Quick Start Guide
Follow the steps below to get HomeShare up and running in under 10 minutes:
#### Prerequisites
 1. **Git**: Make sure you have git installed in your system:
    ```bash
    git --version
    ```
 2. **Docker & Docker Compose**: Install Docker and Docker Compose. Verify installation:
    ```bash
    docker --version
    docker compose version
    ```
 3. **Cloudflare Account**: You need a Cloudflare account with a registered domain name. You’ll use Cloudflare to manage DNS records and issue SSL certificates.


---

#### 1. Configure Cloudflare DNS
 1. Log into your Cloudflare dashboard and select your domain.
 2. Under **Overview**, note down:
    - **Zone ID**
    - **Account ID**
 3. Go to **DNS** and create the following records:
    - **A** record for `@` (root) pointing to any placeholder IP (e.g., `192.0.2.1`).
    - **CNAME** record for `api` pointing to `@` (your root domain).
    - Make sure the **Cloudflare Cloud Icon** is checked, so that traffic is proxied through Cloudflare.

---

#### 2. Generate Cloudflare API Token
 1. In Cloudflare, click your avatar (top right) → **Profile** → **API Tokens**.
 2. Click **Create Token** and use the **Edit DNS** template or a custom template with the following permissions:
    - **Zone.Zone**: Read
    - **Zone.DNS**: Edit
 3. Under **Zone Resources**, select **Include** → **Specific zone** → choose your domain.
 4. Finalize and save the generated token. Keep it handy for later.

---

#### 3. (Optional but Recommended) Reserve a Static IP
Reserving a static local IP address ensures port forwarding remains valid:
 1. Log into your home router configuration page (usually at `http://192.168.0.1` or `http://192.168.1.1`).
 2. Under **DHCP Settings** or *LAN Settings*, create a **DHCP Reservation** for your computer’s MAC address.
 3. Choose a fixed IP address (e.g., `192.168.1.100`).


---

#### 4. Configure Port Forwarding
 1. In your router settings, go to **Port Forwarding** (sometimes under **Advanced**)
 2. Forward **TCP port 443** from the WAN (public) side to port **443** on your computer’s reserved IP address.

---

#### 5. Clone the HomeShare Repository
On your local machine, run:
```bash
git clone https://github.com/jugeekuz/HomeShare.git
cd HomeShare
```

---

#### 6. Create and Populate the .env File
 - In the project root (next to `docker-compose.yml`), create a file named `.env` and add the following variables:
   ```env
   # Cloudflare settings
   CLOUDFLARE_ZONE_ID=<YOUR_ZONE_ID>
   CLOUDFLARE_ACCOUNT_ID=<YOUR_ACCOUNT_ID>
   CF_DNS_API_TOKEN=<YOUR_API_TOKEN>
   CF_API_EMAIL=<YOUR_CLOUDFLARE_EMAIL>
   CLOUDFLARE_RECORD_NAME=<YOUR_ROOT_DOMAIN>   # e.g. example.com

   # Application domain and URLs
   DOMAIN=<YOUR_ROOT_DOMAIN>                   # e.g. example.com
   DOMAIN_ORIGIN=https://${DOMAIN}             # used for CORS and redirects

   # Database credentials (change to secure values)
   POSTGRES_USER=myuser
   POSTGRES_PASSWORD=mypassword
   POSTGRES_DB=userdb

   # File storage folder - files will be uploaded here (absolute or relative path)
   # You can point this to any directory on your host system.
   FILES_UPLOAD_FOLDER="./files"       # e.g. /home/username/hs-files or ./files

   # User Credentials (change to secure values)
   # You will use those credentials to log in to the HomeShare UI
   ADMIN_USERNAME=user@email.com
   ADMIN_EMAIL=user@email.com
   ADMIN_PASSWORD=mypassword
   ```

   > **Notes:**
   > - Replace placeholders (`<...>`) with your actual Cloudflare IDs, token, email, and domain.
   > - You may use an **absolute path** (e.g., `/home/user/hs-files`) or a **relative path** (e.g., `./files`) for `FILES_UPLOAD_FOLDER`.

   - In `traefik/traefik.yml`, change the `changeme@changeme.org` address to your own address:
      ```yml
      certificatesResolvers:
         letencrypt:
            acme:
               caServer: https://acme-v02.api.letsencrypt.org/directory
               email: changeme@changeme.org   # change this value
      ```
---

#### 7. Build and Launch the Stack

From the project root, execute:

```bash
docker compose up -d --build
```

- The `-d` flag runs containers in detached mode.
- `--build` forces a rebuild of the images if you made changes.

---

#### 8. Access the UI:
   - In your browser, go to `https://<YOUR_ROOT_DOMAIN>/` for the frontend.
   - Log in using the credentials we configured before in Step #6
   - Upload files to your computer, or create sharing folders.


### 🎉 Congratulations!
Your HomeShare instance is live and secured by Cloudflare. Enjoy seamless file sharing on your own domain!

---

## 📜 License
Licensed under the GPL V3.0 License.
<a href="https://github.com/jugeekuz/AirRemote-Embedded/blob/master/LICENSE">🔗 View License Details </a>


---

## 🤝 Contributing
Feel free to fork the repository and contribute! Pull requests and feedback are welcome.