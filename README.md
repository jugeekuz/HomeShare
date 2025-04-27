<p align="center">
  <img width="320" src="assets/logo-black.png#gh-light-mode-only" alt="HomeShare Logo">
  <!-- <img width="300" src="assets/logo-white.png#gh-dark-mode-only" alt="HomeShare Logo"> -->
</p>
<!-- <h2 align="center">HomeShare Repository ⚡</h2> -->
<p align="center">
    <a href="/LICENSE"><img alt="GPL-V3.0 License" src="https://img.shields.io/badge/License-GPLv3-orange.svg"></a>
    <a href="https://github.com/jugeekuz/HomeShare/graphs/contributors"><img alt="Contributors" src="https://img.shields.io/github/contributors/jugeekuz/HomeShare?color=green"></a>
</p><br>

<strong>HomeShare</strong> — your own private, behind‑the‑cloud file server. Push or pull files from anywhere, with rock‑solid authentication and time‑limited share links.

<p align="center">
    <img src="./assets/homeshare-demo-short.gif" alt="HomeShare Short Demo" width="320">
</p>


--- 
## 🔍 What is HomeShare?
 Homeshare is a way to 

## 🔧 Container Architecture

| Component | Role | Stack |
|----------|------|------|
| DDNS Updater | Keep your domain pointing home | Go + Cloudflare API |
| Traefik | TLS & reverse‑proxy | Traefik |
| Database | User & share‑link metadata | PostgreSQL |
| File Server | Core upload/download logic | Go |
| Web UI | Frontend interface | React + Nginx |
---


## ⚡ Quick Start
 1. Set a Static IP via DHCP Reservation (Optional - Recommended)
    - Navigate into your home router's settings and reserve the same Private IP Address for your computer.
 2. Set up port forwarding in your home router to your computer.
    - Navigate into your home router's settings and set up port forwarding from port 443 into your computer's Private IP address.
 3. 
---

## 📜 License
Licensed under the GPL V3.0 License.
<a href="https://github.com/jugeekuz/AirRemote-Embedded/blob/master/LICENSE">🔗 View License Details </a>


---

## 🤝 Contributing
Feel free to fork the repository and contribute! Pull requests and feedback are welcome.