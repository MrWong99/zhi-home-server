# Home Server

A [zhi](https://github.com/MrWong99/zhi) workspace that generates and manages a Docker Compose deployment for a self-hosted home server.

## Services

| Component | Service | Description |
|-----------|---------|-------------|
| `core` | *(shared settings)* | Timezone, domain, data root path (mandatory) |
| `pihole` | [PiHole](https://pi-hole.net/) | Network-wide DNS ad-blocking |
| `plex` | [Plex](https://www.plex.tv/) | Media server for movies, TV, music |
| `nextcloud` | [Nextcloud](https://nextcloud.com/) | File sync, sharing, collaboration |
| `mariadb` | [MariaDB](https://mariadb.org/) | Database backend for Nextcloud |
| `redis` | [Redis](https://redis.io/) | Cache and file locking for Nextcloud |
| `nginx-proxy-manager` | [NPM](https://nginxproxymanager.com/) | Reverse proxy with Let's Encrypt SSL |

`nextcloud` depends on `mariadb` and `redis` — enabling Nextcloud automatically enables both.

## Prerequisites

- [zhi CLI](https://github.com/MrWong99/zhi) (v1.1.3+)
- Docker (27.0+) and Docker Compose (v2.20+)
- HashiCorp Vault (running and accessible)
- zhi Vault plugins installed:
  ```sh
  zhi plugin install oci://ghcr.io/mrwong99/zhi/zhi-store-vault:latest
  zhi plugin install oci://ghcr.io/mrwong99/zhi/zhi-store-vault-manager:latest
  ```

### Vault Bootstrap

If you don't have Vault running, use the zhi vault workspace to set one up:

```sh
zhi_version="1.1.3"
zhi workspace install ghcr.io/mrwong99/zhi/zhi-workspace-vault:v${zhi_version} ./vault
cd vault && zhi edit && cd ..
# SAVE THE VAULT UNSEAL KEY(S) IN A SAFE PLACE
```

## Quick Start

```sh
# Clone the repository
git clone <repo-url> home-server
cd home-server

# Verify workspace loads
zhi list components
zhi list paths

# Enable the services you want
zhi component enable pihole
zhi component enable nextcloud   # auto-enables mariadb + redis
zhi component enable plex
zhi component enable nginx-proxy-manager

# Configure interactively
zhi edit

# Or deploy directly (export + docker compose up)
zhi apply
```

## Common Operations

```sh
# Open the interactive configuration editor
zhi edit

# Export the Docker Compose file without deploying
zhi export

# Deploy (runs docker compose up -d --remove-orphans)
zhi apply

# Stop all containers (preserves data volumes)
zhi apply stop

# Restart all containers
zhi apply restart

# Destroy everything including data volumes
zhi apply destroy

# Check configuration for errors
zhi validate

# Enable or disable a service
zhi component enable pihole
zhi component disable pihole
```

## Configuration

Each service has a dedicated config file under `config/`. Configuration values are edited through `zhi edit` (interactive TUI) or by modifying the YAML files directly.

### Required Values

The following values must be set before deploying (enforced by blocking validation):

| Path | Description |
|------|-------------|
| `core/domain` | Base domain (e.g., `home.example.com`) |
| `pihole/admin-password` | PiHole web admin password |
| `nextcloud/admin-password` | Nextcloud admin password |
| `mariadb/root-password` | MariaDB root password |
| `mariadb/nextcloud-password` | MariaDB password for Nextcloud user |

### Network Topology

- **frontend**: Nginx Proxy Manager, PiHole, Nextcloud
- **backend**: MariaDB, Redis, Nextcloud (bridges both)
- **host**: Plex uses `network_mode: host` for DLNA/UPnP discovery

### Volume Strategy

- **Bind mounts** under `${core/data-root}/<service>/` for user-accessible data (Plex config, Nextcloud files)
- **Named Docker volumes** for internal state (MariaDB data, Redis data, NPM data, PiHole config)

## Notes

### systemd-resolved Conflict (PiHole)

On systems with systemd-resolved (most Linux desktops/servers), port 53 is already in use. `zhi validate` will warn about this. To fix:

```sh
sudo sed -i 's/#DNSStubListener=yes/DNSStubListener=no/' /etc/systemd/resolved.conf
sudo systemctl restart systemd-resolved
```

Alternatively, change `pihole/dns-port` to a different port (e.g., 5353).

### First-Run vs Reconfiguration

MariaDB and Nextcloud passwords are only applied during **initial container creation**. Changing passwords in zhi config after the first run will not update the running databases. To change passwords after initial setup, you must also update them manually inside the containers.

### Plex Claim Token

The Plex claim token (`plex/claim-token`) is only needed during first setup to link the server to your Plex account. Get one at https://plex.tv/claim — it expires after 4 minutes.
